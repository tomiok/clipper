package clipper

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const maxFailures = 5
const defaultTimeout = 10

type Status int

type Clipper struct {
	Name     string
	Failures int64
	open     bool
	openedAt int64

	statistics circuitStats

	mutex *sync.Mutex
}

var clippers map[string]*Clipper

func newClipper(c *Configs) *Clipper {
	if clippers == nil {
		clippers = make(map[string]*Clipper)
	}

	return &Clipper{
		Name:  c.Name,
		mutex: &sync.Mutex{},
	}
}

type Configs struct {
	Name             string
	MaxDurationInSec int
}

func setClipper(cfg *Configs) *Clipper {
	c := newClipper(cfg)
	clippers[cfg.Name] = c
	return c
}

func getClipperWithName(name string) *Clipper {
	if name == "" {
		return nil
	}

	cb, _ok := clippers[name]

	if !_ok {
		return nil
	}

	return cb
}

func getClipper(cfg *Configs) *Clipper {
	if cfg == nil {
		randName := randStr()
		log.Println("empty config with name: " + randName)
		log.Println(fmt.Sprintf("default timeout: %d", defaultTimeout))
		cfg = &Configs{
			Name:             randName,
			MaxDurationInSec: defaultTimeout,
		}
	}

	cb, _ok := clippers[cfg.Name]

	if !_ok {
		return setClipper(cfg)
	}

	return cb
}

func (c *Clipper) update(err error) {
	if err != nil {
		c.Failures++
		if c.Failures >= maxFailures {
			c.open = true
			c.openedAt = time.Now().Unix()
			c.statistics.numOfOpenings++
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
