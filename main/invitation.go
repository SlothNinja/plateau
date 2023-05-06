package main

import (
	"encoding/json"

	"cloud.google.com/go/firestore"
	"github.com/Pallinder/go-randomdata"
	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

const invitationKind = "Invitation"

type invitation struct {
	sn.Header
}

// func (inv invitation) MarshalJSON() ([]byte, error) {
// 	opt, err := getOptions(inv.OptString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return json.Marshal(&struct {
// 		ID               string           `json:"id"`
// 		Type             sn.Type          `json:"type"`
// 		Title            string           `json:"title"`
// 		NumPlayers       int              `json:"numPlayers"`
// 		CreatorID        sn.UID           `json:"creatorId"`
// 		CreatorKey       *datastore.Key   `json:"creatorKey"`
// 		CreatorName      string           `json:"creatorName"`
// 		CreatorEmailHash string           `json:"creatorEmailHash"`
// 		CreatorGravType  string           `json:"creatorGravType"`
// 		Status           sn.Status        `json:"status"`
// 		CPIDS            []sn.PID         `json:"cpIds"`
// 		UserIDS          []sn.UID         `json:"userIds"`
// 		UserKeys         []*datastore.Key `json:"userKeys"`
// 		UserNames        []string         `json:"userNames"`
// 		UserEmailHashes  []string         `json:"userEmailHashes"`
// 		UserGravTypes    []string         `json:"userGravTypes"`
// 		WinnerIDS        []sn.UID         `json:"winnerIds"`
// 		HandsPerPlayer   int              `json:"handsPerPlayer"`
// 		CreatedAt        time.Time        `json:"createdAt"`
// 		UpdatedAt        time.Time        `json:"updatedAt"`
// 		LastUpdated      string           `json:"lastUpdated"`
// 		Public           bool             `json:"public"`
// 	}{
// 		ID:               inv.ID,
// 		Type:             inv.Type,
// 		Title:            inv.Title,
// 		NumPlayers:       inv.NumPlayers,
// 		CreatorID:        inv.CreatorID,
// 		CreatorName:      inv.CreatorName,
// 		CreatorEmailHash: inv.CreatorEmailHash,
// 		CreatorGravType:  inv.CreatorGravType,
// 		Status:           inv.Status,
// 		CPIDS:            inv.CPIDS,
// 		UserIDS:          inv.UserIDS,
// 		UserNames:        inv.UserNames,
// 		UserEmailHashes:  inv.UserEmailHashes,
// 		UserGravTypes:    inv.UserGravTypes,
// 		WinnerIDS:        inv.WinnerIDS,
// 		HandsPerPlayer:   opt.HandsPerPlayer,
// 		CreatedAt:        inv.CreatedAt,
// 		UpdatedAt:        inv.UpdatedAt,
// 		LastUpdated:      sn.LastUpdated(inv.UpdatedAt),
// 		Public:           len(inv.PasswordHash) == 0,
// 	})
// }

// func (inv invitation) id() string {
// 	if inv.Ref == nil {
// 		return ""
// 	}
// 	return inv.Ref.ID
// }
//
// func newInvitation(id string) invitation {
// 	return invitation{ID: id}
// }

// func newInvitationKey(id string) *datastore.Key {
// 	return datastore.IDKey(invitationKind, id, rootKey(id))
// }

func invitationDocRef(cl *firestore.Client, id string) *firestore.DocumentRef {
	return invitationCollectionRef(cl).Doc(id)
}

func invitationCollectionRef(cl *firestore.Client) *firestore.CollectionRef {
	sn.Debugf("cl: %#v", cl)
	return cl.Collection(invitationKind)
}

func defaultInvitation() (invitation, error) {
	opt, err := encodeOptions(1)
	if err != nil {
		return invitation{}, err
	}

	var inv invitation
	// Default Values
	inv.Type = sn.Plateau
	inv.Title = randomdata.SillyName()
	inv.NumPlayers = defaultPlayers
	inv.OptString = opt
	return inv, nil
}

func getID(ctx *gin.Context) string {
	return ctx.Param("id")
}

func (cl Client) getInvitation(ctx *gin.Context, id string) (invitation, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	snap, err := invitationDocRef(cl.FS, id).Get(ctx)
	if err != nil {
		return invitation{}, err
	}

	var inv invitation
	err = snap.DataTo(&inv)
	return inv, err
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

func (g game) finalHand() int {
	opts, err := getOptions(g.OptString)
	if err != nil {
		sn.Warningf("err: %v", err)
		return 0
	}
	return opts.HandsPerPlayer * g.NumPlayers
}
