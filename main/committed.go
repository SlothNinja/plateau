package main

import (
	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

const (
	committedKind  = "Committed"
	cachedRootKind = "CachedRoot"
)

func (g *game) committedKey() *datastore.Key {
	return committedKey(g.id())
}

func committedKey(id int64) *datastore.Key {
	return datastore.IDKey(committedKind, id, rootKey(id))
}

func cachedRootKey(id int64, uid sn.UID) *datastore.Key {
	return datastore.IDKey(cachedRootKind, int64(uid), rootKey(id))
}

func newCommitted(id int64) *game {
	return &game{Key: committedKey(id)}
}

func (cl *Client) getCommitted(c *gin.Context) (*game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id, err := getID(c)
	if err != nil {
		return nil, err
	}

	cl.Log.Debugf("getCommitted id: %d", id)
	g := newCommitted(id)
	err = cl.DS.Get(c, g.Key, g)
	cl.Log.Debugf("getCommitted key: %s", g.Key)
	return g, err
}
