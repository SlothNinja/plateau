package client

import (
	"time"

	"github.com/SlothNinja/sn/v3"
)

type glog []*entry

type entry struct {
	Messages    []message
	PID         sn.PID
	HandNumber  int
	TrickNumber int
	Rev         int
	UpdatedAt   time.Time
}

type message map[string]interface{}
