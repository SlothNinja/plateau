package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v3"
)

// state stores the game state of a Tammany Hall game.
type state struct {
	players       []*player
	deck          []card
	declarersTeam []sn.PID
	tricks        []trick
	bids          []bid
}

type jState struct {
	Players       []*player `json:"players"`
	Deck          []card    `json:"deck"`
	DeclarersTeam []sn.PID  `json:"declarersTeam"`
	Tricks        []trick   `json:"tricks"`
	Bids          []bid     `json:"bids"`
}

func (s state) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return json.Marshal(jState{
		Players:       s.players,
		Deck:          s.deck,
		DeclarersTeam: s.declarersTeam,
		Tricks:        s.tricks,
		Bids:          s.bids,
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
	s.declarersTeam = obj.DeclarersTeam
	s.tricks = obj.Tricks
	s.bids = obj.Bids
	return nil
}
