package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v2"
)

// state stores the game state of a Tammany Hall game.
type state struct {
	players   []*player
	deck      deck
	dealerPID sn.PID
	bids      []bid
}

type jState struct {
	Players   []*player `json:"players"`
	Deck      deck      `json:"deck"`
	DealerPID sn.PID    `json:"dealerPID"`
	Bids      []bid     `json:"bids"`
}

func (s state) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return json.Marshal(jState{
		Players:   s.players,
		Deck:      s.deck,
		DealerPID: s.dealerPID,
		Bids:      s.bids,
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
	s.deck = obj.Deck
	s.dealerPID = obj.DealerPID
	s.bids = obj.Bids
	return nil
}
