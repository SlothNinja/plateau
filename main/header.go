package main

import (
	"cloud.google.com/go/firestore"
	"github.com/SlothNinja/sn/v3"
)

const headerKind = "Header"

func (g game) currentHand() int {
	return g.Round
}

func (g *game) nextHand() int {
	g.Round++
	return g.Round
}

func headerDocRef(cl *firestore.Client, id string) *firestore.DocumentRef {
	return cl.Collection(headerKind).Doc(id)
}

// GHeader stores game headers with associate game data.
type GHeader struct {
	sn.Header
}
