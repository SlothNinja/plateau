package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v2"
)

type card struct {
	suit suit
	rank rank
}

func (c card) value() int {
	v := rankValue(c.rank)
	if c.suit == trumps && v != rankValue(noRank) {
		v += 100
	}
	return v
}

type jCard struct {
	Suit suit `json:"suit"`
	Rank rank `json:"rank"`
}

func (c card) MarshalJSON() ([]byte, error) {
	return json.Marshal(jCard{
		Suit: c.suit,
		Rank: c.rank,
	})
}

func (c *card) UnmarshalJSON(bs []byte) error {
	obj := new(jCard)
	err := json.Unmarshal(bs, obj)
	if err != nil {
		return err
	}

	c.suit = obj.Suit
	c.rank = obj.Rank
	return nil
}

type suit string

const (
	noSuit   = ""
	clubs    = "clubs"
	spades   = "spades"
	diamonds = "diamonds"
	hearts   = "hearts"
	trumps   = "trumps"
)

type rank string

const (
	noRank        = ""
	oneRank       = "one"
	twoRank       = "two"
	threeRank     = "three"
	fourRank      = "four"
	fiveRank      = "five"
	sixRank       = "six"
	sevenRank     = "seven"
	eightRank     = "eight"
	nineRank      = "nine"
	tenRank       = "ten"
	elevenRank    = "eleven"
	twelveRank    = "twelve"
	thirteenRank  = "thirteen"
	fourteenRank  = "fourteen"
	fifteenRank   = "fifteen"
	sixteenRank   = "sixteen"
	seventeenRank = "seventeen"
	eighteenRank  = "eighteen"
	nineteenRank  = "nineteen"
	twentyRank    = "twenty"
	twentyoneRank = "twentyone"
	excuseRank    = "excuse"
	valetRank     = "valet"
	cavalierRank  = "cavalier"
	dameRank      = "dame"
	roiRank       = "roi"
)

func rankValue(r rank) int {
	var m = map[rank]int{
		noRank:        0,
		oneRank:       1,
		twoRank:       2,
		threeRank:     3,
		fourRank:      4,
		fiveRank:      5,
		sixRank:       6,
		sevenRank:     7,
		eightRank:     8,
		nineRank:      9,
		tenRank:       10,
		elevenRank:    11,
		twelveRank:    12,
		thirteenRank:  13,
		fourteenRank:  14,
		fifteenRank:   15,
		sixteenRank:   16,
		seventeenRank: 17,
		eighteenRank:  18,
		nineteenRank:  19,
		twentyRank:    20,
		twentyoneRank: 21,
		excuseRank:    22,
		valetRank:     11,
		cavalierRank:  12,
		dameRank:      13,
		roiRank:       14,
	}
	v, ok := m[r]
	if !ok {
		// noRank value = 0
		// a card show always have a rank so noRank is sort of a catchall invalid value
		return 0
		sn.Warningf("invalid rank of %s", r)
	}
	return v
}
