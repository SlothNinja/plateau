package main

import (
	"encoding/json"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
)

const headerKind = "Header"

// Header provides game/invitation header data
type Header struct {
	Key *datastore.Key `datastore:"__key__" json:"-"`
	sn.Header
}

// sn.Header.Turn used to track card played in current hand
// zero based to align with indices of trick slice
func (g *game) trickNumber() int {
	return g.Turn
}

func (g *game) incTrickNumber() int {
	g.Turn++
	return g.Turn
}

func (g *game) resetTrickNumber() int {
	g.Turn = 0
	return g.Turn
}

// sn.Header.Round used to track # of hands
func (g *game) handNumber() int {
	return g.Round
}

func (g *game) incHandNumber() int {
	g.Round++
	return g.Round
}

func (h Header) id() int64 {
	if h.Key == nil {
		return 0
	}
	return h.Key.ID
}

// type dmap map[string]interface{}
//
// func (h Header) data() dmap {
// 	dm := make(dmap)
//
// 	dm["type"] = h.Type
// 	dm["title"] = h.Title
// 	dm["turn"] = h.Turn
// 	dm["phase"] = h.Phase
// 	dm["subPhase"] = h.SubPhase
// 	dm["round"] = h.Round
// 	dm["numPlayers"] = h.NumPlayers
// 	dm["password"] = h.Password
// 	dm["passwordHash"] = h.PasswordHash
// 	dm["creatorId"] = h.CreatorID
// 	dm["creatorKey"] = h.CreatorKey
// 	dm["creatorSId"] = h.CreatorSID
// 	dm["creatorName"] = h.CreatorName
// 	dm["creatorEmail"] = h.CreatorEmail
// 	dm["creatorEmailNotifications"] = h.CreatorEmailNotifications
// 	dm["creatorEmailHash"] = h.CreatorEmailHash
// 	dm["creatorGravType"] = h.CreatorGravType
// 	dm["userIds"] = h.UserIDS
// 	dm["userKeys"] = h.UserKeys
// 	dm["userSIds"] = h.UserSIDS
// 	dm["userNames"] = h.UserNames
// 	dm["userEmails"] = h.UserEmails
// 	dm["userEmailHashes"] = h.UserEmailHashes
// 	dm["userEmailNotifications"] = h.UserEmailNotifications
// 	dm["userGravTypes"] = h.UserGravTypes
// 	dm["cpids"] = h.CPIDS
// 	dm["winnerIndices"] = h.WinnerIDS
// 	dm["winnerKeys"] = h.WinnerKeys
// 	dm["status"] = h.Status
// 	dm["undo"] = h.Undo
// 	dm["progress"] = h.Progress
// 	dm["options"] = h.Options
// 	dm["optString"] = h.OptString
// 	dm["startedAt"] = h.StartedAt
// 	dm["createdAt"] = h.CreatedAt
// 	dm["updatedAt"] = h.UpdatedAt
// 	dm["endedAt"] = h.EndedAt
// 	dm["phase"] = h.Phase
// 	sn.Debugf("phase: %s", h.Phase)
// 	return dm
// }

// // Load of PropertyLoadSaver Interface
// func (h *Header) Load(ps []datastore.Property) error {
// 	return datastore.LoadStruct(h, ps)
// }
//
// // Save of PropertyLoadSaver Interface
// func (h *Header) Save() ([]datastore.Property, error) {
// 	t := time.Now()
// 	if h.CreatedAt.IsZero() {
// 		h.CreatedAt = t
// 	}
// 	h.UpdatedAt = t
// 	return datastore.SaveStruct(h)
// }

type jHeader struct {
	ID int64 `json:"id"`
	sn.Header
}

// MarshalJSON implements json.Marshaler interface
func (h Header) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return json.Marshal(jHeader{
		ID:     h.id(),
		Header: h.Header,
	})
}

func (g *game) headerKey() *datastore.Key {
	return datastore.IDKey(headerKind, g.id(), rootKey(g.id()))
}

// GHeader stores game headers with associate game data.
type GHeader struct {
	Key *datastore.Key `datastore:"__key__"`
	Header
}

func (gh GHeader) id() int64 {
	if gh.Key == nil {
		return 0
	}
	return gh.Key.ID
}

// MarshalJSON implements json.Marshaler interface
func (gh GHeader) MarshalJSON() ([]byte, error) {
	inv := invitation(gh)
	return inv.MarshalJSON()
}

// // Load implements datastore.PropertyLoadSaver interface
// func (gh *GHeader) Load(ps []datastore.Property) error {
// 	return datastore.LoadStruct(gh, ps)
// }
//
// // Save implements datastore.PropertyLoadSaver interface
// func (gh *GHeader) Save() ([]datastore.Property, error) {
// 	t := time.Now()
// 	if gh.CreatedAt.IsZero() {
// 		gh.CreatedAt = t
// 	}
// 	gh.UpdatedAt = t
// 	return datastore.SaveStruct(gh)
// }

// LoadKey implements datastore.LoadKey interface
// func (gh *GHeader) LoadKey(k *datastore.Key) error {
// 	gh.Key = k
// 	return nil
// }
