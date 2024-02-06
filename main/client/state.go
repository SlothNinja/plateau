package client

import (
	"github.com/SlothNinja/sn/v3"
)

// state stores the game state of a Tammany Hall game.
type state struct {
	Deck          []card
	DeclarersTeam []sn.PID
	Tricks        []trick
	Bids          []bid
	LastResults   []lastResult
	Pick          []card
}
