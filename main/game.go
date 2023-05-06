package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/firestore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	gameKind   = "Game"
	cachedKind = "Cached"
	rootKind   = "Root"
)

// Game provides a Le Plateau game.
type game struct {
	// Log glog
	sn.Header
	state
}

func gameDocRef(cl *firestore.Client, id string, rev int) *firestore.DocumentRef {
	return gameCollectionRef(cl).Doc(fmt.Sprintf("%s-%d", id, rev))
}

func gameCollectionRef(cl *firestore.Client) *firestore.CollectionRef {
	return cl.Collection(gameKind)
}

func cachedDocRef(cl *firestore.Client, id string, rev int, uid sn.UID) *firestore.DocumentRef {
	return cachedCollectionRef(cl, id).Doc(fmt.Sprintf("%d-%d", rev, uid))
}

func fullyCachedDocRef(cl *firestore.Client, id string, rev int, uid sn.UID) *firestore.DocumentRef {
	return cachedCollectionRef(cl, id).Doc(fmt.Sprintf("%d-%d-0", rev, uid))
}

func cachedCollectionRef(cl *firestore.Client, id string) *firestore.CollectionRef {
	return committedDocRef(cl, id).Collection(cachedKind)
}

func (g game) rev() int {
	return g.Undo.Current
}

func (g *game) start() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Status = sn.Running
	g.Phase = setupPhase

	g.addNewPlayers()

	// g.newEntry(message{"template": "start-game"})
}

func (g game) dealer() *player {
	return pie.First(g.Players)
}

func (g game) forehand() *player {
	return pie.First(pie.DropTop(g.Players, 1))
}

func (g *game) startBidPhase() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = bidPhase
	pie.Each(g.Players, (*player).bidReset)
	g.Bids = nil
	return g.forehand()
}

func (g *game) randomSeats() {
	g.Players = pie.Shuffle(g.Players, myRandomSource)
	g.updateOrder()
}

// Basically a circular shift left of players so dealer is always first element in slice
func (g *game) newDealer() {
	oldDealer := g.dealer()
	rest := g.Players[1:]
	g.Players = append(rest, oldDealer)
	g.updateOrder()
}

// reflect player order game state to header
func (g *game) updateOrder() {
	g.OrderIDS = pie.Map(g.Players, func(p *player) sn.PID { return p.ID })
}

// currentPlayers returns the players whose turn it is.
func (g game) currentPlayers() []*player {
	return pie.Map(g.CPIDS, func(pid sn.PID) *player { return g.playerByPID(pid) })
}

// currentPlayer returns the player whose turn it is.
func (g game) currentPlayer() *player {
	return pie.First(g.currentPlayers())
}

// Returns player asssociated with user if such player is current player
// Otherwise, return nil
func (g game) currentPlayerFor(u sn.User) *player {
	i := g.IndexFor(u.ID())
	if i == -1 {
		return nil
	}

	return g.playerByPID(i.ToPID())
}

func (g *game) setCurrentPlayers(ps ...*player) {
	g.CPIDS = pie.Map(ps, func(p *player) sn.PID { return p.ID })
}

func (cl Client) getGame(ctx *gin.Context, cu sn.User, action stackFunc) (game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	if cu.IsZero() {
		return cl.getCommitted(ctx)
	}

	undo, err := cl.getStack(ctx, cu.ID())

	if status.Code(err) == codes.NotFound {

		//	if err == datastore.ErrNoSuchEntity {
		g, err := cl.getCommitted(ctx)
		if _, ok := err.(*datastore.ErrFieldMismatch); ok {
			cl.Log.Warningf("err: %v", err)
			return g, nil
		}
		if err != nil {
			return game{}, err
		}
		return g, nil
	}

	if err != nil {
		return game{}, err
	}

	// if undo operation does not transistion to different state, pull current state of game
	if !action(&undo) {
		if undo.Current == undo.Committed {
			g, err := cl.getCommitted(ctx)
			if err != nil {
				return game{}, err
			}
			g.Undo = undo
			return g, nil
		}

		g, err := cl.getCached(ctx, undo.Current, cu.ID())
		if err != nil {
			return game{}, err
		}
		g.Undo = undo
		return g, nil
	}

	// Verify current user is current player, which requires
	// getting the commited game state
	gc, err := cl.getCommitted(ctx)
	if err != nil {
		return game{}, err
	}

	_, err = gc.validateCurrentPlayer(cu)
	if err != nil {
		return game{}, err
	}

	// undo.Current revised by above call of action[0](undo)
	if undo.Current == undo.Committed {
		g, err := cl.getCommitted(ctx)
		if err != nil {
			return game{}, err
		}
		g.Undo = undo
		return g, nil
	}

	g, err := cl.getCached(ctx, undo.Current, cu.ID())
	if err != nil {
		return game{}, err
	}
	g.Undo = undo
	return g, nil
}

func (cl Client) getCached(ctx *gin.Context, rev int, uid sn.UID) (game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	snap, err := fullyCachedDocRef(cl.FS, id, rev, uid).Get(ctx)
	var g game
	err = snap.DataTo(&g)
	return g, err
}

func (cl Client) getRev(ctx *gin.Context, rev int) (game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	snap, err := gameDocRef(cl.FS, id, rev).Get(ctx)
	var g game
	err = snap.DataTo(&g)
	return g, err
}

func (cl Client) save(ctx *gin.Context, g game, cu sn.User) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	return cl.FS.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		return cl.saveGameIn(ctx, tx, g, cu)
	})
}

func (cl Client) saveGameIn(ctx *gin.Context, tx *firestore.Transaction, g game, cu sn.User) error {
	g.UpdatedAt = time.Now()
	id := getID(ctx)

	if err := tx.Set(gameDocRef(cl.FS, id, g.rev()), g); err != nil {
		return err
	}

	if err := tx.Set(committedDocRef(cl.FS, id), g); err != nil {
		return err
	}

	for _, p := range g.Players {
		if err := tx.Set(viewDocRef(cl.FS, id, g.uidForPID(p.ID)), g.viewFor(p)); err != nil {
			return err
		}
	}
	return cl.clearCached(ctx, g, id, cu)
}

// remove hand of other players and deck from data viewed by player
func (g game) viewFor(p *player) game {
	g2 := g.copy()
	for _, p2 := range g2.Players {
		if p.ID != p2.ID {
			p2.Hand = nil
		}
	}
	g2.Deck = nil
	return g2
}

// not truly a deep copy, though state is deeply copied.
func (g game) copy() game {
	return game{
		// Log:    g.Log,
		Header: g.Header,
		state:  g.state.copy(),
	}
}

func (cl Client) commit(ctx *gin.Context, g game, cu sn.User) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	g.Undo.Commit()
	return cl.save(ctx, g, cu)
}

func (cl Client) clearCached(ctx context.Context, g game, id string, cu sn.User) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	refs := cachedCollectionRef(cl.FS, id).DocumentRefs(ctx)
	for {
		ref, err := refs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// if current user is admin, clear all cached docs
		// otherwise clear only if cached doc is for current user
		if cu.Admin || docRefFor(ref, cu.ID()) {
			_, err = ref.Delete(ctx)
			if err != nil {
				return err
			}
		}
	}

	_, err := stackDocRef(cl.FS, id, cu.ID()).Delete(ctx)

	return err
}

func docRefFor(ref *firestore.DocumentRef, uid sn.UID) bool {
	ss := pie.Reverse(strings.Split(ref.ID, "-"))
	s := pie.Pop(&ss)
	if *s == "0" {
		s = pie.Pop(&ss)
	}
	return *s == fmt.Sprintf("%d", uid)
}

func (cl Client) putCached(ctx *gin.Context, g game, rev int, uid sn.UID) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	return cl.FS.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		if err := tx.Set(fullyCachedDocRef(cl.FS, id, rev, uid), g); err != nil {
			return err
		}

		if err := tx.Set(cachedDocRef(cl.FS, id, rev, uid), g.viewFor(g.playerByUID(uid))); err != nil {
			return err
		}

		return tx.Set(stackDocRef(cl.FS, id, uid), g.Undo)
	})
}
