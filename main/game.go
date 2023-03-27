package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v2"
	"github.com/elliotchance/pie/v2"
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
	sn.Header
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

func cachedKey(id, rev int64, uid sn.UID) *datastore.Key {
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
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	err := datastore.LoadStruct(g, ps)
	if err != nil {
		return err
	}

	sn.Debugf("g.EncodedState: %#v", g.EncodedState)
	var s state
	err = json.Unmarshal([]byte(g.EncodedState), &s)
	if err != nil {
		return err
	}
	g.state = s
	sn.Debugf("g.state: %#v", g.state)

	var l glog
	err = json.Unmarshal([]byte(g.EncodedLog), &l)
	if err != nil {
		return err
	}
	g.glog = l
	return nil
}

func (g *game) Save() ([]datastore.Property, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

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

type jGame struct {
	ID     int64     `json:"id"`
	Rev    int64     `json:"rev"`
	Hands  int       `json:"hands"`
	Header sn.Header `json:"header"`
	State  state     `json:"state"`
	GLog   glog      `json:"glog"`
}

func (g *game) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	opt, err := getOptions(g.OptString)
	if err != nil {
		return nil, err
	}

	return json.Marshal(jGame{
		ID:     g.id(),
		Rev:    g.rev(),
		Hands:  opt.HandsPerPlayer * g.NumPlayers,
		Header: g.Header,
		State:  g.state,
		GLog:   g.glog,
	})
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
	g.Phase = setupPhase

	g.addNewPlayers()
	g.randomSeats()

	// set initial dealer
	if p := pie.First(g.players); p != nil {
		g.dealerPID = p.id
	}

	g.deck = deckFor(g.NumPlayers)

	g.newEntry(message{
		"template": "start-game",
		"pids":     pids(g.players),
	})

	g.Phase = dealPhase

	// Standard Round field of header used to track hands/rounds of game
	g.Round = 1

	g.deal()
	g.startBidPhase()
}

func (g *game) dealer() *player {
	return g.playerByPID(g.dealerPID)
}

func (g *game) forehand() *player {
	return g.nextPlayer(g.dealer())
}

func (g *game) startBidPhase() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = bidPhase
	g.setCurrentPlayers(g.forehand())
}

func (g *game) randomSeats() {
	g.players = pie.Shuffle(g.players, myRandomSource)

	// reflect player order game state to header
	g.OrderIDS = pie.Map(g.players, func(p *player) sn.PID { return p.id })
}

// currentPlayers returns the players whose turn it is.
func (g *game) currentPlayers() []*player {
	return pie.Map(g.CPIDS, func(pid sn.PID) *player { return g.playerByPID(pid) })
}

// currentPlayer returns the player whose turn it is.
func (g *game) currentPlayer() *player {
	return pie.First(g.currentPlayers())
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
	sn.Debugf("ps: %#v", ps)
	g.CPIDS = pie.Map(ps, func(p *player) sn.PID { return p.id })
}

// func (g *game) removeCurrentPlayer(p *player) {
// 	if p == nil {
// 		return
// 	}
//
// 	for i, pid := range g.CPIDS {
// 		if pid == p.id {
// 			g.CPIDS = append(g.CPIDS[:i], g.CPIDS[i+1:]...)
// 			return
// 		}
// 	}
// }

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

func (cl *Client) getCached(c *gin.Context, rev int64, uid sn.UID) (*game, error) {
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

func (cl *Client) save(c *gin.Context, g *game, uid sn.UID) error {
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

func (cl *Client) commit(c *gin.Context, g *game, uid sn.UID) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	g.Undo.Commit()
	return cl.save(c, g, uid)
}

func (cl *Client) clearCached(c *gin.Context, g *game, cuid sn.UID) error {
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

func (cl *Client) putCached(c *gin.Context, g *game, rev int64, uid sn.UID) error {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	_, err := cl.DS.RunInTransaction(c, func(tx *datastore.Transaction) error {
		undo, gid := g.Undo, g.id()
		_, err := tx.PutMulti([]*datastore.Key{cachedKey(gid, rev, uid), stackKey(gid, uid)}, []interface{}{g, &undo})
		return err
	})
	return err
}
