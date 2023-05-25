package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
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

func (inv invitation) Start() (*game, sn.PID, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var g game
	g.Header = inv.Header
	g.Status = sn.Running
	g.Phase = setupPhase
	g.StartedAt = updateTime()

	g.addNewPlayers()

	cp := g.startHand()
	g.setCurrentPlayers(cp)
	// g.newEntry(message{"template": "start-game"})
	return &g, cp.ID, nil
}

func (inv invitation) Start2() (*game, sn.PID, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var g game
	g.Header = inv.Header
	g.Status = sn.Running
	g.Phase = setupPhase
	g.StartedAt = updateTime()

	g.addNewPlayers()

	cp := g.startHand()
	g.setCurrentPlayers(cp)
	// g.newEntry(message{"template": "start-game"})
	return &g, cp.ID, nil
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

// func (cl Client) getCached(ctx *gin.Context, rev int, uid sn.UID) (game, error) {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	id := getID(ctx)
// 	snap, err := cl.FullyCachedDocRef(id, rev, uid).Get(ctx)
// 	if err != nil {
// 		return game{}, err
// 	}
//
// 	var g game
// 	if err := snap.DataTo(&g); err != nil {
// 		return game{}, err
// 	}
//
// 	g.ID = id
// 	return g, nil
// }
//
// func (cl Client) getRev(ctx *gin.Context, rev int) (game, error) {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	id := getID(ctx)
// 	snap, err := cl.GameDocRef(id, rev).Get(ctx)
// 	if err != nil {
// 		return game{}, err
// 	}
//
// 	var g game
// 	if err := snap.DataTo(&g); err != nil {
// 		return game{}, err
// 	}
// 	g.ID = id
// 	return g, nil
// }
//
// func (cl Client) save(ctx context.Context, g game, cu sn.User) error {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	return cl.FS.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
// 		return cl.saveGameIn(ctx, tx, g, cu)
// 	})
// }
//
// func (cl Client) saveGameIn(ctx context.Context, tx *firestore.Transaction, g game, cu sn.User) error {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	g.UpdatedAt = updateTime()
//
// 	if err := tx.Set(cl.GameDocRef(g.ID, g.Rev()), g); err != nil {
// 		return err
// 	}
//
// 	if err := tx.Set(cl.CommittedDocRef(g.ID), g); err != nil {
// 		return err
// 	}
//
// 	for _, p := range g.Players {
// 		if err := tx.Set(cl.ViewDocRef(g.ID, g.uidForPID(p.ID)), g.viewFor(p)); err != nil {
// 			return err
// 		}
// 	}
// 	return cl.clearCached(ctx, g, cu)
// }

func (g *game) New() *game {
	return new(game)
}

func (g game) Views() ([]sn.UID, []*game) {
	uids, games := make([]sn.UID, g.NumPlayers), make([]*game, g.NumPlayers)
	for i, p := range g.Players {
		uids[i] = g.uidForPID(p.ID)
		games[i] = g.viewFor(p)
	}
	return uids, games
}

// remove hand of other players and deck from data viewed by player
func (g game) viewFor(p *player) *game {
	g2 := g.copy()
	for _, p2 := range g2.Players {
		if p.ID != p2.ID {
			p2.Hand = nil
		}
		stacksView(p2)
	}
	g2.Deck = nil
	return &g2
}

func stacksView(p *player) {
	stackView(p.Stack0)
	stackView(p.Stack1)
	stackView(p.Stack2)
	stackView(p.Stack3)
	stackView(p.Stack4)
}

func stackView(stack []card) {
	for i, c := range stack {
		if !c.FaceUp {
			stack[i].Rank = noRank
			stack[i].Suit = noSuit
		}
	}
}

// not truly a deep copy, though state is deeply copied.
func (g game) copy() game {
	return game{
		// Log:    g.Log,
		Header: g.Header,
		state:  g.state.copy(),
	}
}

// func (cl Client) commit(ctx context.Context, g game, cu sn.User) error {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	g.Undo.Commit()
// 	return cl.save(ctx, g, cu)
// }
//
// func (cl Client) clearCached(ctx context.Context, g game, cu sn.User) error {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	refs := cl.CachedCollectionRef(g.ID).DocumentRefs(ctx)
// 	for {
// 		ref, err := refs.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
//
// 		// if current user is admin, clear all cached docs
// 		// otherwise clear only if cached doc is for current user
// 		if cu.Admin || docRefFor(ref, cu.ID()) {
// 			_, err = ref.Delete(ctx)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	_, err := cl.StackDocRef(g.ID, cu.ID()).Delete(ctx)
//
// 	return err
// }

// func docRefFor(ref *firestore.DocumentRef, uid sn.UID) bool {
// 	ss := pie.Reverse(strings.Split(ref.ID, "-"))
// 	s := pie.Pop(&ss)
// 	if *s == "0" {
// 		s = pie.Pop(&ss)
// 	}
// 	return *s == fmt.Sprintf("%d", uid)
// }

// func (cl Client) putCached(ctx *gin.Context, g game, rev int, uid sn.UID) error {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	return cl.FS.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
// 		if err := tx.Set(cl.FullyCachedDocRef(g.ID, rev, uid), g); err != nil {
// 			return err
// 		}
//
// 		if err := tx.Set(cl.CachedDocRef(g.ID, rev, uid), g.viewFor(g.playerByUID(uid))); err != nil {
// 			return err
// 		}
//
// 		return tx.Set(cl.StackDocRef(g.ID, uid), g.Undo)
// 	})
// }
