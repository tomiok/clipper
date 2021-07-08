package clipper

import (
	"sync"
	"time"
)

const maxFailures = 5

type Status int

type Clipper struct {
	Name       string
	Failures   int
	TotalFails int
	open       bool
	openedAt   int64

	numOfRuns int
	avgTime   float32

	mutex *sync.Mutex

	paths []string
}

var clippers map[string]*Clipper

func newClipper(name string) *Clipper {
	if clippers == nil {
		clippers = make(map[string]*Clipper)
	}

	return &Clipper{
		Name:  name,
		mutex: &sync.Mutex{},
	}
}

func getClipper(name string) *Clipper {
	cb, ok := clippers[name]

	if !ok {
		c := newClipper(name)
		clippers[name] = c
		return c
	}

	return cb
}

func (c *Clipper) update(err error) {
	if err != nil {
		c.Failures++
		c.TotalFails++
		if c.Failures >= maxFailures {
			c.open = true
			c.openedAt = time.Now().Unix()
			return
		}
	}
	c.open = false
	c.Failures = 0
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
