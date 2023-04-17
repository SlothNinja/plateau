package main

import (
	"encoding/json"
	"time"

	"github.com/SlothNinja/sn/v3"
)

type glog []*entry

func (g *game) newEntry(m ...message) {
	g.glog = append(g.glog, &entry{
		messages:    append([]message(nil), m...),
		handNumber:  g.handNumber(),
		trickNumber: g.trickNumber(),
		rev:         g.rev(),
		updatedAt:   time.Now(),
	})
}

func (g *game) newEntryFor(pid sn.PID, m ...message) {
	g.glog = append(g.glog, &entry{
		messages:    append([]message(nil), m...),
		pid:         pid,
		handNumber:  g.handNumber(),
		trickNumber: g.trickNumber(),
		rev:         g.rev(),
		updatedAt:   time.Now(),
	})
}

type entry struct {
	messages    []message
	pid         sn.PID
	handNumber  int
	trickNumber int
	rev         int64
	updatedAt   time.Time
}

type jEntry struct {
	Messages    []message `json:"messages"`
	PID         sn.PID    `json:"pid,omitempty"`
	HandNumber  int       `json:"handNumber"`
	TrickNumber int       `json:"trickNumber"`
	Rev         int64     `json:"rev"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (e entry) MarshalJSON() ([]byte, error) {
	return json.Marshal(jEntry{
		Messages:    e.messages,
		PID:         e.pid,
		HandNumber:  e.handNumber,
		TrickNumber: e.trickNumber,
		Rev:         e.rev,
		UpdatedAt:   e.updatedAt,
	})
}

func (e *entry) UnmarshalJSON(bs []byte) error {
	var obj jEntry
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		return err
	}
	e.messages = obj.Messages
	e.pid = obj.PID
	e.handNumber = obj.HandNumber
	e.trickNumber = obj.TrickNumber
	e.rev = obj.Rev
	e.updatedAt = obj.UpdatedAt
	return nil
}

func (g *game) appendEntry(m ...message) {
	e := g.lastEntry()
	e.messages = append(e.messages, m...)
	e.updatedAt = time.Now()
	e.rev = g.rev()
}

func (g game) lastEntry() *entry {
	l := len(g.glog)
	if l == 0 {
		return nil
	}
	return g.glog[l-1]
}

type message map[string]interface{}
