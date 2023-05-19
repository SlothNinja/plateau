package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	gameKey   = "Game"
	jsonKey   = "JSON"
	statusKey = "Status"
	hParam    = "hid"
	msgEnter  = "Entering"
	msgExit   = "Exiting"
)

func (cl Client) requireLogin(ctx *gin.Context) (sn.User, error) {
	cu, err := cl.User.Current(ctx)
	if err != nil {
		return sn.User{}, fmt.Errorf("must login to access resource: %w", err)
	}
	return cu, nil
}

// func (cl Client) gamesIndex(c *gin.Context) {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	_, err := cl.requireLogin(c)
// 	if err != nil {
// 		sn.JErr(c, err)
// 		return
// 	}
//
// 	status := sn.ToStatus(c.Param("status"))
// 	q := datastore.
// 		NewQuery(headerKind).
// 		FilterField("Status", "=", status.String()).
// 		Order("-UpdatedAt")
//
// 	var es []*GHeader
// 	_, err = cl.DS.GetAll(c, q, &es)
// 	if err != nil {
// 		sn.JErr(c, err)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{"gheaders": es})
// }

// func (cl Client) showHandler(c *gin.Context) {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	cu, err := cl.User.Current(c)
// 	if err != nil {
// 		cl.Log.Warningf(err.Error())
// 	}
//
// 	g, err := cl.getGame(c, cu, noUndo)
// 	if err != nil {
// 		sn.JErr(c, err)
// 		return
// 	}
//
// 	if cu.IsZero() {
// 		c.JSON(http.StatusOK, gin.H{"game": g, "unread": 0})
// 		return
// 	}
//
// 	// subs, err := cl.getSubsFor(c, g.ID, cu.ID())
// 	// if err != nil {
// 	// 	sn.JErr(c, err)
// 	// 	return
// 	// }
//
// 	// unread, err := cl.MLog.Unread(c, g.ID, cu)
// 	// if err != nil {
// 	// 	sn.JErr(c, err)
// 	// 	return
// 	// }
//
// 	// c.JSON(http.StatusOK, gin.H{"game": g, "unread": unread, "subs": subs})
// 	c.JSON(http.StatusOK, gin.H{})
// }

func (cl Client) cuHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	if cu.IsZero() {
		ctx.JSON(http.StatusOK, gin.H{"CU": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"CU": cu})
}

func (cl Client) resetHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	g, err := cl.getCommitted(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.clearCached(ctx, g, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type stackFunc func(*sn.Stack) bool

var noUndo stackFunc = func(s *sn.Stack) bool { return false }
var undo stackFunc = (*sn.Stack).Undo
var redo stackFunc = (*sn.Stack).Redo

func (cl Client) undoHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	cu, err := cl.requireLogin(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	ref := stackDocRef(cl.FS, id, cu.ID())
	snap, err := ref.Get(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}
	var stack sn.Stack
	err = snap.DataTo(&stack)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	// do nothing if stack does not change
	if !stack.Undo() {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	_, err = ref.Set(ctx, &stack)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (cl Client) redoHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	cu, err := cl.requireLogin(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	ref := stackDocRef(cl.FS, id, cu.ID())
	snap, err := ref.Get(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	var stack sn.Stack
	err = snap.DataTo(&stack)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	// do nothing if stack does not change
	if !stack.Redo() {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	_, err = ref.Set(ctx, &stack)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (cl Client) rollbackHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	err = validateAdmin(cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	obj := struct {
		Rev int
	}{}

	err = ctx.ShouldBind(&obj)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	rev := obj.Rev - 1
	g, err := cl.getRev(ctx, rev)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.save(ctx, g, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (cl Client) rollforwardHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	err = validateAdmin(cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	obj := struct {
		Rev int
	}{}

	err = ctx.ShouldBind(&obj)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	rev := obj.Rev + 1
	g, err := cl.getRev(ctx, rev)
	if status.Code(err) == codes.NotFound {
		sn.JErr(ctx, fmt.Errorf("cannot roll forward any further: %w", sn.ErrValidation))
		return
	}
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.save(ctx, g, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (cl Client) mlogHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// id, err := getID(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// ml, err := cl.MLog.Get(c, id)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// cl.Log.Debugf("ml: %#v", ml.Messages)

	// cu, err := cl.User.Current(c)
	// if err == nil {
	// 	ml, err = cl.MLog.UpdateRead(c, ml, cu)
	// 	if err != nil {
	// 		sn.JErr(c, err)
	// 		return
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"messages": ml.Messages,
	// 	"unread":   0,
	// })
}

func (cl Client) mlogAddHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// cu, err := cl.User.Current(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// id, err := getID(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// obj := struct {
	// 	Message string  `json:"message"`
	// 	Creator sn.User `json:"creator"`
	// }{}

	// err = c.ShouldBind(&obj)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// if obj.Creator.ID() != cu.ID() {
	// 	sn.JErr(c, fmt.Errorf("invalid creator: %w", sn.ErrValidation))
	// 	return
	// }

	// ml, err := cl.MLog.Get(c, id)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// m := ml.AddMessage(cu, obj.Message)
	// _, err = cl.MLog.Put(c, id, ml)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": m,
	// 	"unread":  0,
	// })
}
