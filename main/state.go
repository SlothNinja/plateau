package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v2"
)

// state stores the game state of a Tammany Hall game.
type state struct {
	players []*player
	// wards              wards
	// castleGarden       nationals
	// bag                nationals
	// currentWardID      wardID
	// moveFromWardID     wardID
	// immigrantInTransit nationality
	// slanderedPlayerID  sn.PID
	// slanderNationality nationality
}

func (s state) dmap(dms ...map[string]interface{}) map[string]interface{} {
	var dm map[string]interface{}
	if len(dms) == 1 {
		dm = dms[0]
	} else {
		dm = make(map[string]interface{})
	}

	dm["players"] = s.players
	// dm["wards"] = s.wards
	// dm["castleGarden"] = s.castleGarden
	// dm["bag"] = s.bag
	// dm["currentWardID"] = s.currentWardID
	// dm["moveFromWardID"] = s.moveFromWardID
	// dm["immigrantInTransit"] = s.immigrantInTransit
	// dm["slanderedPlayerID"] = s.slanderedPlayerID
	// dm["slanderNationality"] = s.slanderNationality

	return dm
}

type jState struct {
	Players []*player `json:"players"`
	// Wards              wards       `json:"wards"`
	// CastleGarden       nationals   `json:"castleGarden"`
	// Bag                nationals   `json:"bag"`
	// CurrentWardID      wardID      `json:"currentWardID"`
	// MoveFromWardID     wardID      `json:"moveFromWardID"`
	// ImmigrantInTransit nationality `json:"immigrantInTransit"`
	// SlanderedPlayerID  sn.PID      `json:"slanderedPlayerID"`
	// SlanderNationality nationality `json:"slanderNationality"`
}

func (s state) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return json.Marshal(jState{
		Players: s.players,
		// Wards:              s.wards,
		// CastleGarden:       s.castleGarden,
		// Bag:                s.bag,
		// CurrentWardID:      s.currentWardID,
		// MoveFromWardID:     s.moveFromWardID,
		// ImmigrantInTransit: s.immigrantInTransit,
		// SlanderedPlayerID:  s.slanderedPlayerID,
		// SlanderNationality: s.slanderNationality,
	})
}

func (s *state) UnmarshalJSON(bs []byte) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var obj jState
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		return err
	}
	s.players = obj.Players
	// s.wards = obj.Wards
	// s.castleGarden = obj.CastleGarden
	// s.bag = obj.Bag
	// s.currentWardID = obj.CurrentWardID
	// s.moveFromWardID = obj.MoveFromWardID
	// s.immigrantInTransit = obj.ImmigrantInTransit
	// s.slanderedPlayerID = obj.SlanderedPlayerID
	// s.slanderNationality = obj.SlanderNationality
	return nil
}
