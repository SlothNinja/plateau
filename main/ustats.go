package main

import (
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

// Player stats for a single game
type stats struct {
	// Number of games played at player count
	GamesPlayed int64 `json:"gamesPlayed"`
	// Number of games won at player count
	Won int64 `json:"won"`
	// Number of moves made by player
	Moves int64 `json:"moves"`
	// Amount of time passed between player moves by player
	Think time.Duration `json:"think"`
	// Position player finished (e.g., 1st, 2nd, etc.)
	Finish int `json:"finish"`
}

const ustatsKind = "UStats"

type ustat struct {
	Key *datastore.Key `json:"key" datastore:"__key__"`
	// Below slices
	// Index 0: Total at all player counts
	// Index 1: Reserved
	// Index 2: Reserved
	// Index 3: Total 3P games
	// Index 4: Total 4P games
	// Index 5: Total 5P games

	// Number of games played
	Played []int64 `json:"played"`

	// Number of games won
	Won []int64 `json:"won"`
	// Number of points scored
	Scored []int64 `json:"scored"`
	// Number of moves made by player
	Moves []int64 `json:"moves"`
	// Amount of time passed between player moves by player
	Think []time.Duration `json:"think"`
	// Average amount of time passed between player moves by player
	ThinkAvg []time.Duration `json:"thinkAvg"`
	// Sum of position finishes
	Finish []int64 `json:"finish"`
	// Average finishing position
	FinishAvg []float32 `json:"finishAvg"`
	// Average Score
	ScoreAvg []float32 `json:"scoreAvg"`
	// Win percentage
	WinPercentage []float32 `json:"winPercentage"`
	// Win percentage
	ExpectedWinPercentage []float32 `json:"expectedWinPercentage"`
	// Time at which stats first created
	CreatedAt time.Time `json:"createdAt"`
	// Time at which stats last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

func newUStat(ukey *datastore.Key) *ustat {
	return &ustat{
		Key:                   newUStatsKey(ukey),
		Played:                make([]int64, maxPlayers+1),
		Won:                   make([]int64, maxPlayers+1),
		Scored:                make([]int64, maxPlayers+1),
		Moves:                 make([]int64, maxPlayers+1),
		Think:                 make([]time.Duration, maxPlayers+1),
		ThinkAvg:              make([]time.Duration, maxPlayers+1),
		Finish:                make([]int64, maxPlayers+1),
		FinishAvg:             make([]float32, maxPlayers+1),
		ScoreAvg:              make([]float32, maxPlayers+1),
		WinPercentage:         make([]float32, maxPlayers+1),
		ExpectedWinPercentage: make([]float32, maxPlayers+1),
	}
}

func newUStatsKey(ukey *datastore.Key) *datastore.Key {
	return datastore.NameKey(ustatsKind, "singleton", ukey)
}

func (g *game) updateUStats(stats []*ustat) {
	for i := range stats {
		g.updateUStat(stats[i], g.UserKeys[i])
	}
}

func (g *game) updateUStat(stat *ustat, ukey *datastore.Key) {
	stat.Played[0]++
	stat.Played[g.NumPlayers]++
	for _, key := range g.WinnerKeys {
		if key.Equal(ukey) {
			stat.Won[0]++
			stat.Won[g.NumPlayers]++
			break
		}
	}

	p := g.playerByUserKey(ukey)
	if p == nil {
		return
	}

	stat.Moves[0] += p.Stats.Moves
	stat.Moves[g.NumPlayers] += p.Stats.Moves

	stat.Think[0] += p.Stats.Think
	stat.Think[g.NumPlayers] += p.Stats.Think

	stat.Scored[0] += int64(p.Score)
	stat.Scored[g.NumPlayers] += int64(p.Score)

	stat.Finish[0] += int64(p.Stats.Finish)
	stat.Finish[g.NumPlayers] += int64(p.Stats.Finish)

	if stat.Played[0] > 0 {
		stat.WinPercentage[0] = float32(stat.Won[0]) / float32(stat.Played[0])
		stat.ExpectedWinPercentage[0] = (float32(stat.Played[3])/3.0 + float32(stat.Played[4])/4.0 +
			float32(stat.Played[5])/5.0) / float32(stat.Played[0])
		stat.FinishAvg[0] = float32(stat.Finish[0]) / float32(stat.Played[0])
		stat.ScoreAvg[0] = float32(stat.Scored[0]) / float32(stat.Played[0])
	}
	stat.ThinkAvg[0] = stat.Think[0] / time.Duration(stat.Moves[0])

	if stat.Played[g.NumPlayers] > 0 {
		stat.WinPercentage[g.NumPlayers] = float32(stat.Won[g.NumPlayers]) / float32(stat.Played[g.NumPlayers])
		stat.FinishAvg[g.NumPlayers] = float32(stat.Finish[g.NumPlayers]) / float32(stat.Played[g.NumPlayers])
		stat.ScoreAvg[g.NumPlayers] = float32(stat.Scored[g.NumPlayers]) / float32(stat.Played[g.NumPlayers])
	}
	stat.ExpectedWinPercentage[g.NumPlayers] = 1.0 / float32(g.NumPlayers)
	stat.ThinkAvg[g.NumPlayers] = stat.Think[g.NumPlayers] / time.Duration(stat.Moves[g.NumPlayers])

}

func (cl *Client) getUStats(c *gin.Context, ukeys ...*datastore.Key) ([]*ustat, error) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	l := len(ukeys)
	ustats := make([]*ustat, l)
	ks := make([]*datastore.Key, l)
	for i, ukey := range ukeys {
		ustats[i] = newUStat(ukey)
		ks[i] = ustats[i].Key
	}

	err := cl.DS.GetMulti(c, ks, ustats)
	if err == nil {
		return ustats, nil
	}

	if merr, ok := err.(datastore.MultiError); ok {
		for i, e := range merr {
			if e == nil {
				continue // no error
			}
			if e == datastore.ErrNoSuchEntity {
				ustats[i] = newUStat(ukeys[i])
				ustats[i].CreatedAt = time.Now()
				continue
			}
			if missErr, ok := e.(*datastore.ErrFieldMismatch); ok {
				if missErr.FieldName == "Encoded" ||
					missErr.FieldName == "Stats2P" ||
					missErr.FieldName == "Stats3P" ||
					missErr.FieldName == "Stats4P" {
					continue // ignore error
				}
			}
			if e != nil {
				return nil, err
			}
		}
	}
	return ustats, nil
}
