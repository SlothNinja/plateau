package main

import (
	"encoding/json"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Pallinder/go-randomdata"
	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

const invitationKind = "Invitation"

type invitation struct {
	Key *datastore.Key `json:"-" datastore:"__key__"`
	Header
}

func (inv *invitation) MarshalJSON() ([]byte, error) {
	opt, err := getOptions(inv.OptString)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		ID               int64            `json:"id"`
		Type             sn.Type          `json:"type"`
		Title            string           `json:"title"`
		NumPlayers       int              `json:"numPlayers"`
		CreatorID        sn.UID           `json:"creatorId"`
		CreatorKey       *datastore.Key   `json:"creatorKey"`
		CreatorName      string           `json:"creatorName"`
		CreatorEmailHash string           `json:"creatorEmailHash"`
		CreatorGravType  string           `json:"creatorGravType"`
		Status           sn.Status        `json:"status"`
		CPIDS            []sn.PID         `json:"cpids"`
		UserIDS          []sn.UID         `json:"userIds"`
		UserKeys         []*datastore.Key `json:"userKeys"`
		UserNames        []string         `json:"userNames"`
		UserEmailHashes  []string         `json:"userEmailHashes"`
		UserGravTypes    []string         `json:"userGravTypes"`
		WinnerIDS        []sn.UID         `json:"winnerIds"`
		HandsPerPlayer   int              `json:"handsPerPlayer"`
		CreatedAt        time.Time        `json:"createdAt"`
		UpdatedAt        time.Time        `json:"updatedAt"`
		LastUpdated      string           `json:"lastUpdated"`
		Public           bool             `json:"public"`
	}{
		ID:               inv.id(),
		Type:             inv.Type,
		Title:            inv.Title,
		NumPlayers:       inv.NumPlayers,
		CreatorID:        inv.CreatorID,
		CreatorKey:       inv.CreatorKey,
		CreatorName:      inv.CreatorName,
		CreatorEmailHash: inv.CreatorEmailHash,
		CreatorGravType:  inv.CreatorGravType,
		Status:           inv.Status,
		CPIDS:            inv.CPIDS,
		UserIDS:          inv.UserIDS,
		UserKeys:         inv.UserKeys,
		UserNames:        inv.UserNames,
		UserEmailHashes:  inv.UserEmailHashes,
		UserGravTypes:    inv.UserGravTypes,
		WinnerIDS:        inv.WinnerIDS,
		HandsPerPlayer:   opt.HandsPerPlayer,
		CreatedAt:        inv.CreatedAt,
		UpdatedAt:        inv.UpdatedAt,
		LastUpdated:      sn.LastUpdated(inv.UpdatedAt),
		Public:           len(inv.PasswordHash) == 0,
	})
}

func (inv *invitation) id() int64 {
	if inv == nil || inv.Key == nil {
		return 0
	}
	return inv.Key.ID
}

func newInvitation(id int64) *invitation {
	return &invitation{Key: newInvitationKey(id)}
}

func newInvitationKey(id int64) *datastore.Key {
	return datastore.IDKey(invitationKind, id, rootKey(id))
}

func defaultInvitation() (*invitation, error) {
	opt, err := options(1)
	if err != nil {
		return nil, err
	}

	inv := newInvitation(0)
	// Default Values
	inv.Type = sn.Plateau
	inv.Title = randomdata.SillyName()
	inv.NumPlayers = defaultPlayers
	inv.OptString = opt
	return inv, nil
}

func getID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func (cl *Client) getInvitation(c *gin.Context) (*invitation, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	id, err := getID(c)
	if err != nil {
		return nil, err
	}

	inv := newInvitation(id)
	err = cl.DS.Get(c, inv.Key, inv)
	return inv, err
}

type Options struct {
	HandsPerPlayer int `json:"handsPerPlayer"`
}

func options(hands int) (string, error) {
	options := Options{HandsPerPlayer: hands}

	bs, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func getOptions(encoded string) (*Options, error) {
	options := new(Options)
	if encoded == "" {
		return options, nil
	}
	err := json.Unmarshal([]byte(encoded), options)
	return options, err
}
