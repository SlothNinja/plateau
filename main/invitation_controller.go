package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// if field is a serverTimestamp and field is zero value, firestore will auto-timestamp with server time
// updateTime simply returns zero value, which can be used to zero field and cause server to auto-timestamp
func updateTime() (t time.Time) { return }

func (cl Client) newInvitationHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	inv, err := defaultInvitation()
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Invitation": inv})
}

func (cl Client) createHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	var inv invitation
	hash, err := inv.fromForm(ctx, cu)
	sn.Debugf("hash: %v", hash)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	if _, _, err := invitationCollectionRef(cl.FS).Add(ctx, &inv); err != nil {
		sn.JErr(ctx, err)
		return
	}

	inv2, err := defaultInvitation()
	if err != nil {
		sn.JErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Invitation": inv2,
		"Message":    fmt.Sprintf("%s created game %q", cu.Name, inv.Title),
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

func (inv *invitation) fromForm(ctx *gin.Context, cu sn.User) ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	obj := struct {
		Title           string
		NumPlayers      int
		RoundsPerPlayer int
		Password        string
	}{}

	err := ctx.ShouldBind(&obj)
	if err != nil {
		return nil, err
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
	inv.OptString, err = encodeOptions(rounds)
	if err != nil {
		return nil, err
	}

	var hash []byte
	if len(obj.Password) > 0 {
		hash, err = bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		inv.Private = true
	}
	inv.AddCreator(cu)
	inv.AddUser(cu)
	inv.Status = sn.Recruiting
	inv.Type = sn.Plateau
	return hash, nil
}

type detail struct {
	ID     int64
	ELO    int
	Played int64
	Won    int64
	WP     float32
}

func (cl Client) detailsHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	inv, err := cl.getInvitation(ctx, getID(ctx))
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	uids := make([]sn.UID, len(inv.UserIDS))
	copy(uids, inv.UserIDS)

	if hasUID := pie.Any(inv.UserIDS, func(id sn.UID) bool { return id == cu.ID() }); !hasUID {
		uids = append(uids, cu.ID())
	}

	// elos, err := cl.Elo.GetMulti(c, uids)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// ustats, err := cl.getUStats(c, uids...)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// details := make([]detail, len(elos))
	// for i := range elos {
	// 	played, won, wp := ustats[i].Played, ustats[i].Won, ustats[i].WinPercentage
	// 	details[i] = detail{
	// 		ID:     elos[i].ID,
	// 		ELO:    elos[i].Rating,
	// 		Played: played[0],
	// 		Won:    won[0],
	// 		WP:     wp[0],
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{"details": details})
}

func (cl Client) acceptHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	inv, err := cl.getInvitation(ctx, id)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	obj := struct {
		Password string
	}{}

	err = ctx.ShouldBind(&obj)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	start, err := inv.AcceptWith(cu, []byte(obj.Password), nil)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	if !start {
		inv.UpdatedAt = updateTime()
		_, err = invitationDocRef(cl.FS, id).Set(ctx, &inv)
		if err != nil {
			sn.JErr(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("%s joined game: %s", cu.Name, inv.Title)})
		return
	}

	var g game
	g.Header = inv.Header
	g.start()
	cp := g.startHand()
	g.setCurrentPlayers(cp)

	err = cl.save(ctx, g, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	_, err = invitationDocRef(cl.FS, id).Delete(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	// cl.sendTurnNotificationsTo(c, g, cp)
	// 	err = cl.sendNotifications(c, g)
	// 	if err != nil {
	// 		cl.Log.Warningf(err.Error())
	// 	}
	//
	ctx.JSON(http.StatusOK, gin.H{"Message": g.startGameMessage(cp.ID)})
}

func (g game) startGameMessage(pid sn.PID) string {
	return fmt.Sprintf("<div>Game: %s has started.</div><div></div><div><strong>%s</strong> is start player.</div>",
		g.Title, g.NameFor(pid))
}

func (cl Client) dropHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id := getID(ctx)
	inv, err := cl.getInvitation(ctx, id)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	cu, err := cl.requireLogin(ctx)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = inv.Drop(cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	if len(inv.UserIDS) == 0 {
		inv.Status = sn.Aborted
	}

	inv.UpdatedAt = updateTime()
	_, err = invitationDocRef(cl.FS, id).Set(ctx, &inv)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("%s dropped from game invitation: %s", cu.Name, inv.Title)})
}
