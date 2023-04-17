package main

import (
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"firebase.google.com/go/messaging"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

const subKind = "Sub"
const subParentKey = "SubPK"

// Sub provides support for players to subscribe to real-time updates
type Sub struct {
	Key       *datastore.Key `datastore:"__key__" json:"-"`
	Token     string         `json:"token"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// Load of PropertyLoadSaver Interface
func (s *Sub) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(s, ps)
}

// Save of PropertyLoadSaver Interface
func (s *Sub) Save() ([]datastore.Property, error) {
	t := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = t
	}
	s.UpdatedAt = t
	return datastore.SaveStruct(s)
}

// LoadKey of KeyLoader Interface
func (s *Sub) LoadKey(k *datastore.Key) error {
	s.Key = k
	return nil
}

func newSubParentKey(gid int64, uid sn.UID) *datastore.Key {
	return datastore.IDKey(subParentKey, int64(uid), rootKey(gid))
}

func newSubKey(id, gid int64, uid sn.UID) *datastore.Key {
	return datastore.IDKey(subKind, id, newSubParentKey(gid, uid))
}

// uid return the user ID associated with the subscription
func (s *Sub) uid() sn.UID {
	if s.Key != nil && s.Key.Parent != nil {
		return sn.UID(s.Key.Parent.ID)
	}
	return 0
}

func (cl *Client) subHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	gid, err := getID(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	obj := struct {
		Token string `json:"token"`
		CUID  sn.UID `json:"cuid"`
	}{}

	err = c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	subs, err := cl.getSubsFor(c, gid, obj.CUID)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	const notFound = -1
	if pie.FindFirstUsing(subs, func(s *Sub) bool { return s.Token == obj.Token }) == notFound {
		sub := &Sub{Token: obj.Token}
		sub.Key, err = cl.DS.Put(c, newSubKey(0, gid, obj.CUID), sub)
		if err != nil {
			sn.JErr(c, err)
			return
		}
		subs = append(subs, sub)
	}

	// Found or not, return subs
	c.JSON(http.StatusOK, gin.H{
		"subs": subs,
	})
}

func (cl *Client) getSubsFor(c *gin.Context, gid int64, cuid sn.UID) ([]*Sub, error) {
	var subs []*Sub
	_, err := cl.DS.GetAll(c, datastore.NewQuery(subKind).Ancestor(newSubParentKey(gid, cuid)), &subs)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}
	return subs, err
}

func (cl *Client) sendNotifications(c *gin.Context, g game) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	var subs []*Sub
	_, err := cl.DS.GetAll(c, datastore.NewQuery(subKind).Ancestor(rootKey(g.id())), &subs)
	if err != nil {
		return err
	}

	// Send refresh notifications to other users
	otherSubs := g.otherSubs(subs)
	if len(otherSubs) > 0 {
		resp, err := cl.Messaging.SendMulticast(c, &messaging.MulticastMessage{
			Tokens: pie.Map(otherSubs, func(sub *Sub) string { return sub.Token }),
			Data:   map[string]string{"action": "refresh"},
		})
		if resp != nil {
			cl.Log.Debugf("batch response: %+v", resp)
			for _, r := range resp.Responses {
				cl.Log.Debugf("response: %+v", r)
			}
		}
		if err != nil {
			return err
		}
	}

	// Send turn notifications to users associated with current player(s)
	cpSubs := g.cpSubs(subs)
	for _, sub := range cpSubs {
		_, err := cl.Messaging.Send(c, &messaging.Message{
			Token: sub.Token,
			Data: map[string]string{
				"title": "SlothNinja Games Turn Notification",
				"body":  "It is your turn in one or more games.",
			},
			Webpush: &messaging.WebpushConfig{
				Notification: &messaging.WebpushNotification{
					Title: "SlothNinja Games Turn Notification",
					Body:  "It is your turn in one or more games.",
				},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *game) otherSubs(subs []*Sub) []*Sub {
	if len(g.CPIDS) == 0 {
		return subs
	}

	uids := g.uidsForPIDS(g.CPIDS)
	return pie.FilterNot(subs, func(s *Sub) bool { return pie.Contains(uids, s.uid()) })
}

func (g *game) cpSubs(subs []*Sub) []*Sub {
	if len(g.CPIDS) == 0 {
		return subs
	}

	uids := g.uidsForPIDS(g.CPIDS)
	sn.Debugf("cpSubs uids: %v", uids)
	return pie.Filter(subs, func(s *Sub) bool {
		sn.Debugf("s.UID: %v", s.uid())
		return pie.Contains(uids, s.uid())
	})
}
