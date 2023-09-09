package client

import (
	"time"

	"github.com/SlothNinja/sn/v3"
)

type glog []*entry

// func (g *game) newEntry(m ...message) {
// 	g.Log = append(g.Log, &entry{
// 		Messages:    append([]message(nil), m...),
// 		HandNumber:  g.currentHand(),
// 		TrickNumber: g.trickNumber(),
// 		Rev:         g.rev(),
// 		UpdatedAt:   time.Now(),
// 	})
// }
//
// func (g *game) newEntryFor(pid sn.PID, m ...message) {
// 	g.Log = append(g.Log, &entry{
// 		Messages:    append([]message(nil), m...),
// 		PID:         pid,
// 		HandNumber:  g.currentHand(),
// 		TrickNumber: g.trickNumber(),
// 		Rev:         g.rev(),
// 		UpdatedAt:   time.Now(),
// 	})
// }

type entry struct {
	Messages    []message
	PID         sn.PID
	HandNumber  int
	TrickNumber int
	Rev         int
	UpdatedAt   time.Time
}

// func (g *game) appendEntry(m ...message) {
// 	e := g.lastEntry()
// 	e.Messages = append(e.Messages, m...)
// 	e.UpdatedAt = time.Now()
// 	e.Rev = g.rev()
// }

// func (g game) lastEntry() *entry {
// 	l := len(g.Log)
// 	if l == 0 {
// 		return nil
// 	}
// 	return g.Log[l-1]
// }

type message map[string]interface{}
