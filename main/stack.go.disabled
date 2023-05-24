package main

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

const stackKind = "Stack"

// func stackKey(id int64, uid sn.UID) *datastore.Key {
// 	return datastore.NameKey(stackKind, "stack", cachedRootKey(id, uid))
// }

// func stackDocRef(cl *firestore.Client, id string, uid sn.UID) *firestore.DocumentRef {
// 	return cl.Collection(stackKind).Doc(fmt.Sprintf("%s-%d", id, uid))
// }

func stackDocRef(cl *firestore.Client, id string, uid sn.UID) *firestore.DocumentRef {
	return stackCollectionRef(cl).Doc(fmt.Sprintf("%s-%d", id, uid))
}

func stackCollectionRef(cl *firestore.Client) *firestore.CollectionRef {
	return cl.Collection(stackKind)
}

func (cl *Client) getStack(ctx *gin.Context, uid sn.UID) (sn.Stack, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	snap, err := stackDocRef(cl.FS, getID(ctx), uid).Get(ctx)
	// gid, err := getID(c)
	if err != nil {
		return sn.Stack{}, err
	}

	var s sn.Stack
	err = snap.DataTo(&s)
	// err = cl.DS.Get(c, stackKey(gid, uid), &s)
	return s, err
}
