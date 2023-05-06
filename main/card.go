package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

type card struct {
	Rank     rank
	Suit     suit
	PlayedBy sn.PID
}

func (c card) value() int {
	v := c.Rank.value()
	if c.Suit == trumps && v != noRank.value() {
		v += 100
	}
	return v
}

type suit string

const (
	noSuit   suit = ""
	clubs    suit = "clubs"
	spades   suit = "spades"
	diamonds suit = "diamonds"
	hearts   suit = "hearts"
	trumps   suit = "trumps"
)

type rank string

const (
	noRank        rank = ""
	oneRank       rank = "one"
	twoRank       rank = "two"
	threeRank     rank = "three"
	fourRank      rank = "four"
	fiveRank      rank = "five"
	sixRank       rank = "six"
	sevenRank     rank = "seven"
	eightRank     rank = "eight"
	nineRank      rank = "nine"
	tenRank       rank = "ten"
	elevenRank    rank = "eleven"
	twelveRank    rank = "twelve"
	thirteenRank  rank = "thirteen"
	fourteenRank  rank = "fourteen"
	fifteenRank   rank = "fifteen"
	sixteenRank   rank = "sixteen"
	seventeenRank rank = "seventeen"
	eighteenRank  rank = "eighteen"
	nineteenRank  rank = "nineteen"
	twentyRank    rank = "twenty"
	twentyoneRank rank = "twentyone"
	excuseRank    rank = "excuse"
	valetRank     rank = "valet"
	cavalierRank  rank = "cavalier"
	dameRank      rank = "dame"
	roiRank       rank = "roi"
)

func (r rank) value() int {
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
		// a card should always have a rank so noRank is sort of a catchall invalid value
		sn.Warningf("invalid rank of %s", r)
		return 0
	}
	return v
}

// used for setting trick 'rank' from an index value
func toRank(v int) rank {
	var m = map[int]rank{
		0:  noRank,
		1:  oneRank,
		2:  twoRank,
		3:  threeRank,
		4:  fourRank,
		5:  fiveRank,
		6:  sixRank,
		7:  sevenRank,
		8:  eightRank,
		9:  nineRank,
		10: tenRank,
		11: elevenRank,
		12: twelveRank,
		13: thirteenRank,
		14: fourteenRank,
		15: fifteenRank,
		16: sixteenRank,
	}
	r, ok := m[v]
	if !ok {
		// noRank value = 0
		// a card should always have a rank so noRank is sort of a catchall invalid value
		sn.Warningf("invalid value of %d", v)
		return noRank
	}
	return r
}

func getCards(ctx *gin.Context) ([]card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var obj []card
	err := ctx.ShouldBind(&obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (g game) cardsFor(team []sn.PID) []card {
	var cards []card
	pie.Each(g.tricksFor(team), func(t trick) {
		cards = append(cards, pie.Map(t.Cards, func(c card) card {
			c.PlayedBy = sn.NoPID
			return c
		})...)
	})
	return cards
}

func (c card) toSpace() space {
	return space{c.Rank, kind(c.Suit)}
}
