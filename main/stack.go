package main

import (
	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

const stackKind = "Stack"

func stackKey(id int64, uid sn.UID) *datastore.Key {
	return datastore.NameKey(stackKind, "stack", cachedRootKey(id, uid))
}

func (cl *Client) getStack(c *gin.Context, uid sn.UID) (*sn.Stack, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	gid, err := getID(c)
	if err != nil {
		return nil, err
	}

	s := new(sn.Stack)
	err = cl.DS.Get(c, stackKey(gid, uid), s)
	return s, err
}
