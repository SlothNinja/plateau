package main

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (cl *Client) newInvitationHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	inv, err := defaultInvitation()
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"invitation": inv})
}

func (cl *Client) createHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	inv := newInvitation(0)
	err = inv.fromForm(c, cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	ks, err := cl.DS.AllocateIDs(c, []*datastore.Key{rootKey(0)})
	if err != nil {
		sn.JErr(c, err)
		return
	}
	t := time.Now()
	inv.Key = newInvitationKey(ks[0].ID)
	inv.CreatedAt, inv.UpdatedAt = t, t

	_, err = cl.DS.RunInTransaction(c, func(tx *datastore.Transaction) error {
		m := sn.NewMLog(inv.Key.ID)
		m.UpdatedAt, m.CreatedAt = t, t
		ks := []*datastore.Key{inv.Key, m.Key}
		es := []interface{}{inv, m}

		_, err := tx.PutMulti(ks, es)
		return err
	})
	if err != nil {
		sn.JErr(c, err)
		return
	}

	inv2, err := defaultInvitation()
	if err != nil {
		sn.JErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"invitation": inv2,
		"message":    fmt.Sprintf("%s created game %q", cu.Name, inv.Title),
	})
}

const (
	minPlayers     = 2
	defaultPlayers = 3
	maxPlayers     = 6

	minRounds     = 1
	defaultRounds = 1
	maxRounds     = 5
)

func (inv *invitation) fromForm(c *gin.Context, cu *sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	obj := struct {
		Title           string `json:"title"`
		NumPlayers      int    `json:"numPlayers"`
		RoundsPerPlayer int    `json:"roundsPerPlayer"`
		Password        string `json:"password"`
	}{}

	err := c.ShouldBind(&obj)
	if err != nil {
		return err
	}

	inv.Title = cu.Name + "'s Game"
	if obj.Title != "" {
		inv.Title = obj.Title
	}

	inv.NumPlayers = defaultPlayers
	if obj.NumPlayers >= minPlayers && obj.NumPlayers <= maxPlayers {
		inv.NumPlayers = obj.NumPlayers
	}

	rounds := defaultRounds
	if obj.RoundsPerPlayer >= minRounds && obj.RoundsPerPlayer <= maxRounds {
		rounds = obj.RoundsPerPlayer
	}
	inv.OptString, err = options(rounds)
	if err != nil {
		return err
	}

	if len(obj.Password) > 0 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		inv.PasswordHash = hashed
	}
	inv.AddCreator(cu)
	inv.AddUser(cu)
	inv.Status = sn.Recruiting
	inv.Type = sn.Plateau
	return nil
}

func (cl *Client) invitationsIndexHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// obj := struct {
	// 	Options struct {
	// 		ItemsPerPage int `json:"itemsPerPage"`
	// 	} `json:"options"`
	// 	Forward string `json:"forward"`
	// }{}

	// err := c.ShouldBind(&obj)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// cu, err := cl.User.Current(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// forward, err := datastore.DecodeCursor(obj.Forward)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	q := datastore.
		NewQuery(invitationKind).
		// FilterField("Status", "=", sn.Recruiting.String()).
		FilterField("Status", "=", sn.Recruiting.String()).
		Order("-UpdatedAt")

	var es []*invitation
	_, err := cl.DS.GetAll(c, q, &es)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	for _, inv := range es {
		cl.Log.Debugf("inv: %#v", inv)
	}

	// cnt, err := cl.DS.Count(c, q)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// cl.Log.Debugf("cnt: %d", cnt)

	// items := obj.Options.ItemsPerPage
	// if obj.Options.ItemsPerPage == -1 {
	// 	items = cnt
	// }

	// var es []*invitation
	// it := cl.DS.Run(c, q.Start(forward))
	// for i := 0; i < items; i++ {
	// 	var inv invitation
	// 	_, err := it.Next(&inv)
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		sn.JErr(c, err)
	// 		return
	// 	}
	// 	es = append(es, &inv)
	// }

	// forward, err = it.Cursor()
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"invitations": es,
		// "totalItems":  cnt,
		// "forward":     forward.String(),
		// "cu":          cu,
	})
}

type detail struct {
	ID     int64   `json:"id"`
	ELO    int     `json:"elo"`
	Played int64   `json:"played"`
	Won    int64   `json:"won"`
	WP     float32 `json:"wp"`
}

func (cl *Client) detailsHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	inv, err := cl.getInvitation(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cu, err := cl.requireLogin(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	uids := make([]sn.UID, len(inv.UserIDS))
	copy(uids, inv.UserIDS)

	if hasUID := pie.Any(inv.UserIDS, func(id sn.UID) bool { return id == cu.ID() }); !hasUID {
		uids = append(uids, cu.ID())
	}

	elos, err := cl.Elo.GetMulti(c, uids)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	ustats, err := cl.getUStats(c, uids...)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	details := make([]detail, len(elos))
	for i := range elos {
		played, won, wp := ustats[i].Played, ustats[i].Won, ustats[i].WinPercentage
		details[i] = detail{
			ID:     elos[i].Key.Parent.ID,
			ELO:    elos[i].Rating,
			Played: played[0],
			Won:    won[0],
			WP:     wp[0],
		}
	}

	for i, d := range details {
		cl.Log.Debugf("details[%d]: %#v", i, d)
	}

	c.JSON(http.StatusOK, gin.H{"details": details})
}

func (cl *Client) acceptHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	inv, err := cl.getInvitation(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cu, err := cl.requireLogin(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	obj := struct {
		Password string `json:"password"`
	}{}

	err = c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	start, err := inv.AcceptWith(cu, []byte(obj.Password))
	if err != nil {
		cl.Log.Debugf("accept err: %v", err)
		sn.JErr(c, err)
		return
	}

	if !start {
		_, err = cl.DS.Put(c, inv.Key, inv)
		if err != nil {
			sn.JErr(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"invitation": inv,
			"message":    fmt.Sprintf("%s joined game: %d", cu.Name, inv.Key.ID),
		})
		return
	}

	g := newGame(inv.Key.ID, 0)
	g.Header = inv.Header.Header
	cp := g.start()
	g.setCurrentPlayers(cp)

	_, err = cl.DS.RunInTransaction(c, func(tx *datastore.Transaction) error {
		cl.Log.Debugf("inv.Key: %s", inv.Key)
		err = tx.Delete(inv.Key)
		if err != nil {
			return err
		}

		g.StartedAt = time.Now()
		h := g.Header
		_, err = tx.PutMulti([]*datastore.Key{g.headerKey(), committedKey(g.id()), g.Key},
			[]interface{}{&h, g, g})
		return err
	})
	if err != nil {
		sn.JErr(c, err)
		return
	}

	// cl.sendTurnNotificationsTo(c, g, cp)
	err = cl.sendNotifications(c, g)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"invitation": g.Header,
		"message": fmt.Sprintf(
			`<div>Game: %d has started.</div>
			<div></div>
			<div><strong>%s</strong> is start player.</div>`,
			g.id(), g.NameFor(cp.id)),
	})
}

func (cl *Client) dropHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	inv, err := cl.getInvitation(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cu, err := cl.requireLogin(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = inv.Drop(cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	if len(inv.UserKeys) == 0 {
		inv.Status = sn.Aborted
	}

	_, err = cl.DS.Put(c, inv.Key, inv)
	if err != nil {
		sn.JErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"invitation": inv,
		"message":    fmt.Sprintf("%s dropped from game invitation: %d", cu.Name, inv.Key.ID),
	})
}
