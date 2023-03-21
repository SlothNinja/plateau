package main

import (
	"encoding/json"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Pallinder/go-randomdata"
	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

const invitationKind = "Invitation"

type invitation struct {
	Key *datastore.Key `datastore:"__key__"`
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
		CreatorID        int64            `json:"creatorId"`
		CreatorKey       *datastore.Key   `json:"creatorKey"`
		CreatorName      string           `json:"creatorName"`
		CreatorEmailHash string           `json:"creatorEmailHash"`
		CreatorGravType  string           `json:"creatorGravType"`
		UserIDS          []int64          `json:"userIds"`
		UserKeys         []*datastore.Key `json:"userKeys"`
		UserNames        []string         `json:"userNames"`
		UserEmailHashes  []string         `json:"userEmailHashes"`
		UserGravTypes    []string         `json:"userGravTypes"`
		RoundsPerPlayer  int              `json:"roundsPerPlayer"`
		CreatedAt        time.Time        `json:"createdAt"`
		UpdatedAt        time.Time        `json:"updatedAt"`
		LastUpdated      string           `json:"lastUpdated"`
		Public           bool             `json:"public"`
	}{
		ID:               inv.ID(),
		Type:             inv.Type,
		Title:            inv.Title,
		NumPlayers:       inv.NumPlayers,
		CreatorID:        inv.CreatorID,
		CreatorKey:       inv.CreatorKey,
		CreatorName:      inv.CreatorName,
		CreatorEmailHash: inv.CreatorEmailHash,
		CreatorGravType:  inv.CreatorGravType,
		UserIDS:          inv.UserIDS,
		UserKeys:         inv.UserKeys,
		UserNames:        inv.UserNames,
		UserEmailHashes:  inv.UserEmailHashes,
		UserGravTypes:    inv.UserGravTypes,
		RoundsPerPlayer:  opt.RoundsPerPlayer,
		CreatedAt:        inv.CreatedAt,
		UpdatedAt:        inv.UpdatedAt,
		LastUpdated:      sn.LastUpdated(inv.UpdatedAt),
		Public:           len(inv.PasswordHash) == 0,
	})
}

func (inv *invitation) ID() int64 {
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

func (inv *invitation) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(inv, ps)
}

func (inv *invitation) Save() ([]datastore.Property, error) {
	t := time.Now()
	if inv.CreatedAt.IsZero() {
		inv.CreatedAt = t
	}
	inv.UpdatedAt = t
	return datastore.SaveStruct(inv)
}

func (inv *invitation) LoadKey(k *datastore.Key) error {
	inv.Key = k
	return nil
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

// func toInvitation(h *sn.Header) (*invitation, error) {
// 	opts, err := getOptions(h.OptString)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &invitation{
// 		ID:               h.ID(),
// 		AdmiralVariant:   opts.AdmiralVariant,
// 		BasicGame:        opts.BasicGame,
// 		Title:            h.Title,
// 		Public:           len(h.PasswordHash) == 0,
// 		Status:           h.Status,
// 		CreatorID:        h.CreatorID,
// 		CreatorName:      h.CreatorName,
// 		CreatorEmailHash: h.CreatorEmailHash,
// 		CreatorGravType:  h.CreatorGravType,
// 		NumPlayers:       h.NumPlayers,
// 		UserIDS:          h.UserIDS,
// 		UserNames:        h.UserNames,
// 		UserEmailHashes:  h.UserEmailHashes,
// 		UserGravTypes:    h.UserGravTypes,
// 		LastUpdated:      sn.LastUpdated(h.UpdatedAt),
// 	}, nil
// }

type Options struct {
	RoundsPerPlayer int `json:"roundsPerPlayer"`
}

func options(rounds int) (string, error) {
	options := Options{RoundsPerPlayer: rounds}

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
