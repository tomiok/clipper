package clipper

import (
	"sync"
	"time"
)

const maxFailures = 5

type Status int

type Clipper struct {
	Name        string
	MaxFailures int
	open        bool
	openedAt    int64

	numOfRuns int
	avgTime   float32

	mutex sync.Mutex
}

func NewClipper(name string) *Clipper {
	return &Clipper{
		Name:        name,
		MaxFailures: 0,
	}
}

var clippers map[string]*Clipper

func getClipper(name string) *Clipper {
	cb, ok := clippers[name]

	if !ok {
		c := NewClipper(name)
		clippers[name] = c
		return c
	}

	return cb
}

func (c *Clipper) update(err error) {
	if err != nil {
		c.MaxFailures++
		if c.MaxFailures >= maxFailures {
			c.open = true
			c.openedAt = time.Now().Unix()
			return
		}
	}
	c.open = false
	c.MaxFailures = 0
}

func (c *Clipper) isOpen() bool {
	if c.open {
		now := time.Now().Unix()
		// 3 minutes
		if (now - c.openedAt) > 180 {
			c.open = false
			return false
		} else {
			return true
		}
	}
	return false
}
