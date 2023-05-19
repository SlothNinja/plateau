package main

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

const (
	committedKind  = "Committed"
	viewKind       = "View"
	cachedRootKind = "CachedRoot"
)

//	func (g game) committedKey(cl *firestore.Client) *firestore.DocumentRef {
//		return committedKey(cl, g.id())
//	}
//
//	func committedKey(cl *firestore.Client, id string) *firestore.DocumentRef {
//		return rootKey(cl, id).Collection(committedKind).Doc(id)
//	}
func committedDocRef(cl *firestore.Client, id string) *firestore.DocumentRef {
	return cl.Collection(committedKind).Doc(id)
}

func viewDocRef(cl *firestore.Client, id string, uid sn.UID) *firestore.DocumentRef {
	return committedDocRef(cl, id).Collection(viewKind).Doc(fmt.Sprintf("%d", uid))
}

// func cachedRootKey(cl *firestore.Client, id string, uid sn.UID) *firestore.DocumentRef {
// 	return rootKey(cl, id).Collection(cachedRootKind).Doc(fmt.Sprintf("%d", uid))
// }

// func newCommitted(id string) game {
// 	return game{ID: id}
// }

func (cl Client) getCommitted(ctx *gin.Context) (game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	snap, err := committedDocRef(cl.FS, id).Get(ctx)
	if err != nil {
		return game{}, err
	}

	var g game
	if err := snap.DataTo(&g); err != nil {
		return game{}, err
	}

	g.id = id
	return g, nil
}
