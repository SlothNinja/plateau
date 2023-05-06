package main

import (
	"testing"

	"github.com/SlothNinja/sn/v3"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/path"
)

func TestConnected2(t *testing.T) {

	var g game
	g.DeclarersTeam = []sn.PID{1}
	g.Tricks = []trick{
		trick{[]card{
			card{oneRank, trumps, 1},
			card{cavalierRank, hearts, 2},
			card{excuseRank, trumps, 3},
			card{roiRank, spades, 4},
			card{roiRank, hearts, 5},
			card{twentyoneRank, trumps, 6},
		}, 1},
		trick{[]card{
			card{cavalierRank, diamonds, 1},
			card{fourRank, trumps, 2},
			card{excuseRank, trumps, 3},
			card{tenRank, spades, 4},
			card{tenRank, hearts, 5},
			card{tenRank, clubs, 6},
		}, 1},
	}

	spaces := g.spacesFor(g.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	s1 := graph.side(1)
	s2 := graph.side(2)
	s3 := graph.side(3)
	s4 := graph.side(4)
	s5 := graph.side(5)
	s6 := graph.side(6)

	path, found := connected(paths, s1, s2)
	assert.False(t, found, "path should not be found between side1 and side2")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s3)
	assert.True(t, found, "path should be found between side1 and side3")
	assert.NotNilf(t, path, "path not should be nil: %#v", path)

	path, found = connected(paths, s1, s4)
	assert.True(t, found, "path should be found between side1 and side4")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, found = connected(paths, s1, s5)
	assert.False(t, found, "path should not be found between side1 and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s6)
	assert.True(t, found, "path should not be found between side1 and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)
}

func TestConnected3(t *testing.T) {

	var g game
	g.DeclarersTeam = []sn.PID{1}
	g.Tricks = []trick{
		trick{[]card{
			card{oneRank, trumps, 1},
			card{cavalierRank, hearts, 2},
			card{excuseRank, trumps, 3},
			card{roiRank, spades, 4},
			card{roiRank, hearts, 5},
			card{twentyoneRank, trumps, 6},
		}, 1},
		trick{[]card{
			card{cavalierRank, diamonds, 1},
			card{fourRank, trumps, 2},
			card{excuseRank, trumps, 3},
			card{tenRank, spades, 4},
			card{tenRank, hearts, 5},
			card{tenRank, clubs, 6},
		}, 1},
	}

	spaces := g.spacesFor(g.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	s1 := graph.side(1)
	s2 := graph.side(2)
	s3 := graph.side(3)
	s4 := graph.side(4)
	s5 := graph.side(5)
	s6 := graph.side(6)

	path, found := connected(paths, s1, s2, s3)
	assert.False(t, found, "path should not be found between side1, side2, and side3")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s2, s4)
	assert.False(t, found, "path should not be found between side1, side2, and side4")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s2, s5)
	assert.False(t, found, "path should not be found between side1, side2, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s2, s6)
	assert.False(t, found, "path should not be found between side1, side2, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s3, s4)
	assert.True(t, found, "path should be found between side1, side3, and side4")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, found = connected(paths, s1, s3, s5)
	assert.False(t, found, "path should not be found between side1, side3, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s3, s6)
	assert.True(t, found, "path should be found between side1, side3, and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, found = connected(paths, s1, s4, s5)
	assert.False(t, found, "path should not be found between side1, side4, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s1, s4, s6)
	assert.True(t, found, "path should be found between side1, side4, and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, found = connected(paths, s1, s5, s6)
	assert.False(t, found, "path should not be found between side1, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s2, s3, s4)
	assert.False(t, found, "path should not be found between side2, side3, and side4")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s2, s3, s5)
	assert.False(t, found, "path should not be found between side2, side3, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s2, s3, s6)
	assert.False(t, found, "path should not be found between side2, side3, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s3, s4, s5)
	assert.False(t, found, "path should not be found between side3, side4, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s3, s5, s6)
	assert.False(t, found, "path should not be found between side3, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, found = connected(paths, s4, s5, s6)
	assert.False(t, found, "path should not be found between side4, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)
}
