package client

import (
	"testing"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/path"
)

func TestConnected2(t *testing.T) {

	g := newGame()
	g.State.DeclarersTeam = []sn.PID{1}
	g.State.Tricks = []trick{
		trick{[]card{
			card{oneRank, trumps, 1, true},
			card{cavalierRank, hearts, 2, true},
			card{excuseRank, trumps, 3, true},
			card{roiRank, spades, 4, true},
			card{roiRank, hearts, 5, true},
			card{twentyoneRank, trumps, 6, true},
		}, 1},
		trick{[]card{
			card{cavalierRank, diamonds, 1, true},
			card{fourRank, trumps, 2, true},
			card{excuseRank, trumps, 3, true},
			card{tenRank, spades, 4, true},
			card{tenRank, hearts, 5, true},
			card{tenRank, clubs, 6, true},
		}, 1},
	}

	spaces := g.spacesFor(g.State.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	s1 := g.side(graph, 1)
	s2 := g.side(graph, 2)
	s3 := g.side(graph, 3)
	s4 := g.side(graph, 4)
	s5 := g.side(graph, 5)
	s6 := g.side(graph, 6)

	path, connected := isConnected(paths, s1, s2)
	assert.Falsef(t, connected, "path should not be found between side1 and side2")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3)
	assert.Truef(t, connected, "path should be found between side1 and side3")
	assert.NotNilf(t, path, "path not should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s4)
	assert.Truef(t, connected, "path should be found between side1 and side4")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s5)
	assert.Falsef(t, connected, "path should not be found between side1 and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s6)
	assert.Truef(t, connected, "path should not be found between side1 and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)
}

func TestConnected3(t *testing.T) {

	g := newGame()
	g.State.DeclarersTeam = []sn.PID{1}
	g.State.Tricks = []trick{
		trick{[]card{
			card{oneRank, trumps, 1, true},
			card{cavalierRank, hearts, 2, true},
			card{excuseRank, trumps, 3, true},
			card{roiRank, spades, 4, true},
			card{roiRank, hearts, 5, true},
			card{twentyoneRank, trumps, 6, true},
		}, 1},
		trick{[]card{
			card{cavalierRank, diamonds, 1, true},
			card{fourRank, trumps, 2, true},
			card{excuseRank, trumps, 3, true},
			card{tenRank, spades, 4, true},
			card{tenRank, hearts, 5, true},
			card{tenRank, clubs, 6, true},
		}, 1},
	}

	spaces := g.spacesFor(g.State.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	s1 := g.side(graph, 1)
	s2 := g.side(graph, 2)
	s3 := g.side(graph, 3)
	s4 := g.side(graph, 4)
	s5 := g.side(graph, 5)
	s6 := g.side(graph, 6)

	path, connected := isConnected(paths, s1, s2, s3)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side3")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s2, s4)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side4")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s2, s5)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s2, s6)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3, s4)
	assert.Truef(t, connected, "path should be found between side1, side3, and side4")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3, s5)
	assert.Falsef(t, connected, "path should not be found between side1, side3, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3, s6)
	assert.Truef(t, connected, "path should be found between side1, side3, and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s4, s5)
	assert.Falsef(t, connected, "path should not be found between side1, side4, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s4, s6)
	assert.Truef(t, connected, "path should be found between side1, side4, and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s5, s6)
	assert.Falsef(t, connected, "path should not be found between side1, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s2, s3, s4)
	assert.Falsef(t, connected, "path should not be found between side2, side3, and side4")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s2, s3, s5)
	assert.Falsef(t, connected, "path should not be found between side2, side3, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s2, s3, s6)
	assert.Falsef(t, connected, "path should not be found between side2, side3, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s3, s4, s5)
	assert.Falsef(t, connected, "path should not be found between side3, side4, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s3, s5, s6)
	assert.Falsef(t, connected, "path should not be found between side3, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s4, s5, s6)
	assert.Falsef(t, connected, "path should not be found between side4, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)
}

func TestY(t *testing.T) {

	g := newGame()
	g.Header.NumPlayers = 2
	g.State.DeclarersTeam = []sn.PID{1}
	g.State.Bids = []bid{bid{exchangeBid, yBid, noTeamBid, 1}}
	g.State.Tricks = []trick{
		// Trick 0
		trick{
			Cards: []card{
				card{Rank: "roi", Suit: "clubs", PlayedBy: 1, FaceUp: true},
				card{Rank: "ten", Suit: "clubs", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 1
		trick{
			Cards: []card{
				card{Rank: "dame", Suit: "hearts", PlayedBy: 1, FaceUp: true},
				card{Rank: "cavalier", Suit: "hearts", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 2
		trick{
			Cards: []card{
				card{Rank: "twentyone", Suit: "trumps", PlayedBy: 1, FaceUp: true},
				card{Rank: "six", Suit: "trumps", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 3
		trick{
			Cards: []card{
				card{Rank: "twenty", Suit: "trumps", PlayedBy: 1, FaceUp: true},
				card{Rank: "eighteen", Suit: "trumps", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 4
		trick{
			Cards: []card{
				card{Rank: "nineteen", Suit: "trumps", PlayedBy: 1, FaceUp: true},
				card{Rank: "one", Suit: "trumps", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 5
		trick{
			Cards: []card{
				card{Rank: "dame", Suit: "spades", PlayedBy: 1, FaceUp: true},
				card{Rank: "roi", Suit: "spades", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 2,
		},
		// Trick 6
		trick{
			Cards: []card{
				card{Rank: "cavalier", Suit: "spades", PlayedBy: 2, FaceUp: true},
				card{Rank: "five", Suit: "trumps", PlayedBy: 1, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 7
		trick{
			Cards: []card{
				card{Rank: "roi", Suit: "diamonds", PlayedBy: 1, FaceUp: true},
				card{Rank: "valet", Suit: "diamonds", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 8
		trick{
			Cards: []card{
				card{Rank: "cavalier", Suit: "diamonds", PlayedBy: 1, FaceUp: true},
				card{Rank: "two", Suit: "trumps", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 2,
		},
		// Trick 9
		trick{
			Cards: []card{
				card{Rank: "ten", Suit: "spades", PlayedBy: 2, FaceUp: true},
				card{Rank: "four", Suit: "trumps", PlayedBy: 1, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 10
		trick{
			Cards: []card{
				card{Rank: "dame", Suit: "clubs", PlayedBy: 1, FaceUp: true},
				card{Rank: "valet", Suit: "clubs", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// Trick 11
		trick{
			Cards: []card{
				card{Rank: "excuse", Suit: "trumps", PlayedBy: 1, FaceUp: true},
				card{Rank: "seven", Suit: "trumps", PlayedBy: 2, FaceUp: true},
			},
			WonBy: 1,
		},
		// // Trick 12
		// trick{
		// 	Cards: []card{
		// 		card{Rank: "sixteen", Suit: "trumps", PlayedBy: 1, FaceUp: true},
		// 		card{Rank: "three", Suit: "trumps", PlayedBy: 2, FaceUp: true},
		// 	},
		// 	WonBy: 1,
		// },
		// Trick 13
		// trick{
		// 	Cards: []card{
		// 		card{Rank: "dame", Suit: "diamonds", PlayedBy: 1, FaceUp: true},
		// 		card{Rank: "valet", Suit: "spades", PlayedBy: 2, FaceUp: true},
		// 	},
		// 	WonBy: 1,
		// },
		// // Trick 14
		// trick{
		// 	Cards: []card{
		// 		card{Rank: "cavalier", Suit: "clubs", PlayedBy: 1, FaceUp: true},
		// 		card{Rank: "valet", Suit: "hearts", PlayedBy: 2, FaceUp: true},
		// 	},
		// 	WonBy: 1,
		// },
	}
	spaces := g.spacesFor(g.State.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	// _, result := g.y(graph, paths)
	// assert.Equalf(t, dSuccess, result, "graph should include 'y' path")

	// _, connected := isConnected(paths, g.side(graph, 2), g.side(graph, 4), g.side(graph, 6))
	// assert.Truef(t, connected, "side 2, 4, and 6 should be connected")

	fourTrickID := graph.ids[space{fourRank, trickKind}]
	fourTrickNode := graph.Node(fourTrickID).(node)
	assert.Truef(t, pie.Contains(g.side(graph, 2), fourTrickNode), "four trick space should be on side 2")

	cavSpadeID := graph.ids[space{cavalierRank, spadeKind}]
	assert.Truef(t, pie.Contains(g.side(graph, 2), fourTrickNode), "four trick space should be on side 2")
	p1, _, _ := paths.Between(fourTrickID, cavSpadeID)
	sn.Debugf("p1: %#v", p1)

	roiClubID := graph.ids[space{roiRank, clubKind}]
	p2, _, _ := paths.Between(fourTrickID, roiClubID)
	sn.Debugf("p2: %#v", p2)

	twentyoneTrumpID := graph.ids[space{twentyoneRank, trumpKind}]
	p3, _, _ := paths.Between(fourTrickID, twentyoneTrumpID)
	sn.Debugf("p3: %#v", p3)

	twentyoneTrumpNode := graph.Node(twentyoneTrumpID).(node)
	neighbors := g.neighbors()[twentyoneTrumpNode.Space]
	sn.Debugf("neighbors: %#v", neighbors)

	tenTrickID := graph.ids[space{tenRank, trickKind}]
	p4, _, _ := paths.Between(fourTrickID, tenTrickID)
	sn.Debugf("p4: %#v", p4)

	_, connected := isConnected(paths, g.side(graph, 2), g.side(graph, 4))
	assert.Truef(t, connected, "side 2 and 4 should be connected")
}
