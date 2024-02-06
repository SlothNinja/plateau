package client

import (
	"encoding/json"
	"time"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type invitation struct {
	sn.Header
}

func getID(ctx *gin.Context) string {
	return ctx.Param("id")
}

type options struct {
	HandsPerPlayer int
}

func encodeOptions(hands int) (string, error) {
	options := options{HandsPerPlayer: hands}

	bs, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func getOptions(encoded string) (*options, error) {
	options := new(options)
	if encoded == "" {
		return options, nil
	}
	err := json.Unmarshal([]byte(encoded), options)
	return options, err
}

func (g *game) finalHand() int {
	opts, err := getOptions(g.Header.OptString)
	if err != nil {
		sn.Warningf("err: %v", err)
		return 0
	}
	return opts.HandsPerPlayer * g.Header.NumPlayers
}

// if field is a serverTimestamp and field is zero value, firestore will auto-timestamp with server time
// updateTime simply returns zero value, which can be used to zero field and cause server to auto-timestamp
func updateTime() (t time.Time) { return }

const (
	minPlayers     = 2
	defaultPlayers = 3
	maxPlayers     = 6

	minRounds     = 1
	defaultRounds = 1
	maxRounds     = 5
)

func (inv *invitation) FromForm(ctx *gin.Context, cu sn.User) (*invitation, []byte, error) {
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
		return nil, nil, err
	}

	inv = new(invitation)
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
		return nil, nil, err
	}

	var hash []byte
	if len(obj.Password) > 0 {
		hash, err = bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}
		inv.Private = true
	}
	inv.AddCreator(cu)
	inv.AddUser(cu)
	inv.Status = sn.Recruiting
	inv.Type = sn.Plateau
	return inv, hash, nil
}
