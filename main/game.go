package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

const (
	gameKind   = "Game"
	cachedKind = "Cached"
	rootKind   = "Root"
)

// Game provides a Tammany Hall game.
type game struct {
	Key          *datastore.Key `datastore:"__key__"`
	EncodedState string         `datastore:",noindex"`
	EncodedLog   string         `datastore:",noindex"`
	Header
	glog
	state
}

func rootKey(id int64) *datastore.Key {
	return datastore.IDKey(rootKind, id, nil)
}

func newGame(id, rev int64) *game {
	return &game{Key: newGameKey(id, rev)}
}

func (g *game) gameKey() *datastore.Key {
	return datastore.NameKey(gameKind, fmt.Sprintf("%d-%d", g.id(), g.Undo.Committed), rootKey(g.id()))
}

func newGameKey(id, rev int64) *datastore.Key {
	return datastore.NameKey(gameKind, fmt.Sprintf("%d-%d", id, rev), rootKey(id))
}

func cachedKey(id, rev, uid int64) *datastore.Key {
	return datastore.IDKey(gameKind, rev, cachedRootKey(id, uid))
}

func (g *game) id() int64 {
	if g.Key == nil || g.Key.Parent == nil {
		return 0
	}
	return g.Key.Parent.ID
}

func (g *game) rev() int64 {
	if g.Key == nil {
		return 0
	}
	s := strings.Split(g.Key.Name, "-")
	if len(s) != 2 {
		return g.Undo.Current
	}
	rev, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		sn.Warningf(err.Error())
		return 0
	}
	return rev
}

func (g *game) Load(ps []datastore.Property) error {
	err := datastore.LoadStruct(g, ps)
	if err != nil {
		return err
	}

	var s state
	err = json.Unmarshal([]byte(g.EncodedState), &s)
	if err != nil {
		return err
	}
	g.state = s

	var l glog
	err = json.Unmarshal([]byte(g.EncodedLog), &l)
	if err != nil {
		return err
	}
	g.glog = l
	return nil
}

func (g *game) Save() ([]datastore.Property, error) {

	encodedState, err := json.Marshal(g.state)
	if err != nil {
		return nil, err
	}
	g.EncodedState = string(encodedState)

	encodedLog, err := json.Marshal(g.glog)
	if err != nil {
		return nil, err
	}
	g.EncodedLog = string(encodedLog)

	t := time.Now()
	if g.CreatedAt.IsZero() {
		g.CreatedAt = t
	}
	g.UpdatedAt = t

	return datastore.SaveStruct(g)
}

func (g *game) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	dm := g.Header.data()
	dm = g.state.dmap(dm)

	dm["key"] = g.Key
	dm["id"] = g.id()
	dm["log"] = g.glog
	dm["rev"] = g.rev()

	return json.Marshal(dm)
}

// // CurrentWard returns the ward currently conducting an election.
// func (g *game) CurrentWard() *ward {
// 	return g.wardByID(g.currentWardID)
// }
//
// func (g *game) wardByID(wid wardID) *ward {
// 	return g.wards[wid]
// }
//
// func (g *game) setCurrentWard(w *ward) {
// 	wid := noWardID
// 	if w != nil {
// 		wid = w.ID
// 	}
// 	g.currentWardID = wid
// }
//
// func (g *game) moveFromWard() *ward {
// 	return g.wardByID(g.moveFromWardID)
// }
//
// func (g *game) setMoveFromWard(w *ward) {
// 	wid := noWardID
// 	if w != nil {
// 		wid = w.ID
// 	}
// 	g.moveFromWardID = wid
// }
//
// // Term provides the current game term.
// func (g *game) Term() int {
// 	return (g.Round + 3) / 4
// }
//
// // Year provides the current game year.
// func (g *game) Year() int {
// 	return g.Round
// }

func (g *game) setYear(y int) {
	g.Round = y
}

// Games provides a slice of Games.
type Games []*game

func (g *game) start() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Status = sn.Running
	g.Phase = setup

	g.addNewPlayers()
	g.randomTurnOrder()

	// g.setYear(1)
	// g.wards = newWards()
	// g.setMoveFromWard(nil)
	// g.setCurrentWard(nil)
	// g.bag = defaultBag()
	// g.castleGarden = make(nationals)
	// g.immigration()

	g.newEntry(message{
		"template": "start-game",
		"pids":     pids(g.players),
	})

	// g.castleGardenPhase(g.currentPlayer())
	g.startActionsPhase()
}

// func (g *game) startNextTerm(cp *player) {
// 	g.setYear(g.Year() + 1)
//
// 	for _, p := range g.players {
// 		p.reset()
// 		p.LockedUp = 0
// 	}
//
// 	g.beginningOfTurnReset()
// 	g.unlockWards()
//
// 	g.immigration()
// 	g.castleGardenPhase(cp)
// 	g.startActionsPhase()
// }
//
// func (g *game) unlockWards() {
// 	for _, w := range g.activeWards() {
// 		w.LockedUp = false
// 	}
// }

func (g *game) startActionsPhase() {
	g.Phase = actions
}

// func (g *game) startElections(c *gin.Context, cp *player) bool {
// 	sn.Debugf(msgEnter)
// 	defer sn.Debugf(msgExit)
//
// 	g.Phase = elections
// 	g.emptyGarden()
// 	g.beginningOfTurnReset()
//
// 	for _, p := range g.players {
// 		p.reset()
// 	}
//
// 	for _, w := range g.activeWards() {
// 		w.Resolved = false
// 	}
//
// 	return g.continueElections(c)
// }
//
// func (g *game) candidates() []*player {
// 	var cs []*player
// 	for _, p := range g.players {
// 		if p.Candidate {
// 			cs = append(cs, p)
// 		}
// 	}
// 	return cs
// }

func (g *game) newTurnOrder(c *gin.Context) {
	//	if g.mayor() != nil {
	//		index, _ := g.indexFor(g.mayor())
	//		playersTwice := append(g.players, g.players...)
	//		newOrder := playersTwice[index : index+g.NumPlayers]
	//		g.players = newOrder
	//	}
}

func (g *game) setCurrentPlayer(p *player) {
	if len(g.players) < 1 {
		return
	}
	g.CPIDS = []sn.PID{g.players[0].ID}
}

func (g *game) randomTurnOrder() {
	rand.Shuffle(len(g.players), func(i, j int) {
		g.players[i], g.players[j] = g.players[j], g.players[i]
	})
	g.setCurrentPlayer(g.players[0])

	g.OrderIDS = make([]sn.PID, len(g.players))
	for i, p := range g.players {
		g.OrderIDS[i] = p.ID
	}
}

// currentPlayers returns the players whose turn it is.
func (g *game) currentPlayers() []*player {
	l := len(g.CPIDS)
	if l < 1 {
		return nil
	}

	ps := make([]*player, l)
	for i, id := range g.CPIDS {
		ps[i] = g.playerByPID(id)
	}
	return ps
}

// currentPlayer returns the player whose turn it is.
func (g *game) currentPlayer() *player {
	cps := g.currentPlayers()
	if len(cps) == 0 {
		return nil
	}
	return cps[0]
}

// Returns player asssociated with user if such player is current player
// Otherwise, return nil
func (g *game) currentPlayerFor(u *sn.User) *player {
	if u == nil {
		return nil
	}

	i := g.IndexFor(u.ID())
	if i == -1 {
		return nil
	}

	return g.playerByPID(i.ToPID())
}

func (g *game) setCurrentPlayers(ps ...*player) {
	g.CPIDS = nil
	if len(ps) < 1 {
		return
	}
	for _, p := range ps {
		g.CPIDS = append(g.CPIDS, p.ID)
	}
}

func (g *game) removeCurrentPlayer(p *player) {
	if p == nil {
		return
	}

	for i, pid := range g.CPIDS {
		if pid == p.ID {
			g.CPIDS = append(g.CPIDS[:i], g.CPIDS[i+1:]...)
			return
		}
	}
}

func (cl *Client) getGame(c *gin.Context, cu *sn.User, action stackFunc) (*game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	if cu == nil {
		return cl.getCommitted(c)
	}

	undo, err := cl.getStack(c, cu.ID())

	if err == datastore.ErrNoSuchEntity {
		g, err := cl.getCommitted(c)
		if _, ok := err.(*datastore.ErrFieldMismatch); ok {
			cl.Log.Warningf("err: %v", err)
			return g, nil
		}
		if err != nil {
			return nil, err
		}
		return g, nil
	}

	if err != nil {
		return nil, err
	}

	// if undo operation does not transistion to different state, pull current state of game
	if !action(undo) {
		if undo.Current == undo.Committed {
			g, err := cl.getCommitted(c)
			if err != nil {
				return nil, err
			}
			g.Undo = *undo
			return g, nil
		}

		g, err := cl.getCached(c, undo.Current, cu.ID())
		if err != nil {
			return nil, err
		}
		g.Undo = *undo
		return g, nil
	}

	// Verify current user is current player or admin, which requires
	// getting the commited game state
	gc, err := cl.getCommitted(c)
	if err != nil {
		return nil, err
	}

	_, err = gc.validateCPorAdmin(cu)
	if err != nil {
		return nil, err
	}

	// undo.Current revised by above call of action[0](undo)
	if undo.Current == undo.Committed {
		g, err := cl.getCommitted(c)
		if err != nil {
			return nil, err
		}
		g.Undo = *undo
		return g, nil
	}

	g, err := cl.getCached(c, undo.Current, cu.ID())
	if err != nil {
		return nil, err
	}
	g.Undo = *undo
	return g, nil
}

func (cl *Client) getCached(c *gin.Context, rev, uid int64) (*game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id, err := getID(c)
	if err != nil {
		return nil, err
	}

	g := newGame(id, rev)
	err = cl.DS.Get(c, cachedKey(id, rev, uid), g)
	return g, err
}

func (cl *Client) getRev(c *gin.Context, rev int64) (*game, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id, err := getID(c)
	if err != nil {
		return nil, err
	}

	g := newGame(id, rev)
	err = cl.DS.Get(c, g.Key, g)
	return g, err
}

func (cl *Client) save(c *gin.Context, g *game, uid int64) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	_, err := cl.DS.RunInTransaction(c, func(tx *datastore.Transaction) error {
		h := g.Header
		_, err := tx.PutMulti([]*datastore.Key{g.headerKey(), g.gameKey(), g.committedKey()},
			[]interface{}{&h, g, g})
		if err != nil {
			return err
		}
		return cl.clearCached(c, g, uid)
	})
	return err
}

func (cl *Client) commit(c *gin.Context, g *game, uid int64) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	g.Undo.Commit()
	return cl.save(c, g, uid)
}

func (cl *Client) clearCached(c *gin.Context, g *game, cuid int64) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	ks, err := cl.DS.GetAll(c, datastore.NewQuery("").Ancestor(cachedRootKey(g.id(), cuid)).KeysOnly(), nil)
	if err != nil {
		return err
	}

	if len(ks) == 0 {
		return nil
	}
	return cl.DS.DeleteMulti(c, ks)
}

func (cl *Client) putCached(c *gin.Context, g *game, rev, uid int64) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	_, err := cl.DS.RunInTransaction(c, func(tx *datastore.Transaction) error {
		undo, gid := g.Undo, g.id()
		_, err := tx.PutMulti([]*datastore.Key{cachedKey(gid, rev, uid), stackKey(gid, uid)}, []interface{}{g, &undo})
		return err
	})
	return err
}
