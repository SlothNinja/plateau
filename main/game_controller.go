package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

const (
	gameKey   = "Game"
	jsonKey   = "JSON"
	statusKey = "Status"
	hParam    = "hid"
	msgEnter  = "Entering"
	msgExit   = "Exiting"
)

func (cl *Client) requireLogin(c *gin.Context) (*sn.User, error) {
	cu, err := cl.User.Current(c)
	if err != nil {
		return nil, fmt.Errorf("must login to access resource: %w", err)
	}
	return cu, nil
}

func (cl *Client) gamesIndex(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	_, err := cl.requireLogin(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	status := sn.ToStatus(c.Param("status"))
	q := datastore.
		NewQuery(headerKind).
		FilterField("Status", "=", status.String()).
		Order("-UpdatedAt")

	var es []*GHeader
	_, err = cl.DS.GetAll(c, q, &es)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"gheaders": es})
}

func (cl *Client) showHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(c, cu, noUndo)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	subs, err := cl.getSubsFor(c, g.id(), cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	if cu == nil {
		c.JSON(http.StatusOK, gin.H{
			"game":   g,
			"unread": 0,
			"cu":     cu,
		})
		return
	}

	unread, err := cl.MLog.Unread(c, g.id(), cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game":   g,
		"unread": unread,
		"subs":   subs,
		"cu":     cu,
	})
}

func (cl *Client) cuHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"cu": cu})
}

func (cl *Client) resetHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	g, err := cl.getCommitted(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = cl.clearCached(c, g, cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": g,
		"cu":   cu,
	})
}

type stackFunc func(*sn.Stack) bool

var noUndo stackFunc = func(s *sn.Stack) bool { return false }
var undo stackFunc = (*sn.Stack).Undo
var redo stackFunc = (*sn.Stack).Redo

func (cl *Client) undoHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(c, cu, undo)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	stack := g.Undo
	_, err = cl.DS.Put(c, stackKey(g.id(), cu.ID()), &stack)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": g,
		"cu":   cu,
	})
}

func (cl *Client) redoHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(c, cu, redo)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	undo := g.Undo
	_, err = cl.DS.Put(c, stackKey(g.id(), cu.ID()), &undo)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": g,
		"cu":   cu,
	})
}

func (cl *Client) rollbackHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	err = validateAdmin(cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	obj := struct {
		Rev int64 `json:"rev"`
	}{}

	err = c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}
	cl.Log.Debugf("obj.Rev: %v", obj.Rev)

	var g *game
	for rev := obj.Rev - 1; rev >= 0; rev-- {
		cl.Log.Debugf("rev: %v", rev)
		g, err = cl.getRev(c, rev)
		if err == datastore.ErrNoSuchEntity {
			continue
		}
		if err == nil {
			break
		}
		sn.JErr(c, err)
		return
	}

	err = cl.save(c, g, cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": g,
		"cu":   cu,
	})
}

func (cl *Client) rollforwardHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	err = validateAdmin(cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	obj := struct {
		Rev int64 `json:"rev"`
	}{}

	err = c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	rev := obj.Rev + 1
	g, err := cl.getRev(c, rev)
	if err == datastore.ErrNoSuchEntity {
		sn.JErr(c, fmt.Errorf("cannot roll forward any further: %w", sn.ErrValidation))
		return
	}
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = cl.save(c, g, cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": g,
		"cu":   cu,
	})
}

func (cl *Client) mlogHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id, err := getID(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	ml, err := cl.MLog.Get(c, id)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cl.Log.Debugf("ml: %#v", ml.Messages)

	cu, err := cl.User.Current(c)
	if err == nil {
		ml, err = cl.MLog.UpdateRead(c, ml, cu)
		if err != nil {
			sn.JErr(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": ml.Messages,
		"unread":   0,
	})
}

func (cl *Client) mlogAddHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	id, err := getID(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	obj := struct {
		Message string   `json:"message"`
		Creator *sn.User `json:"creator"`
	}{}

	err = c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	if obj.Creator.ID() != cu.ID() {
		sn.JErr(c, fmt.Errorf("invalid creator: %w", sn.ErrValidation))
		return
	}

	ml, err := cl.MLog.Get(c, id)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	m := ml.AddMessage(cu, obj.Message)
	_, err = cl.MLog.Put(c, id, ml)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": m,
		"unread":  0,
	})
}
