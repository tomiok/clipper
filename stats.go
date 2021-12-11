package clipper

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type Stats struct {
	NumOfCircuits int     `json:"num_of_circuits"`
	IsOpen        bool    `json:"is_open"`
	TotalRuns     int64   `json:"total_runs"`
	TotalFails    int64   `json:"total_fails"`
	AvgSuccess    float64 `json:"avg_success"`
	NumOfOpenings int     `json:"num_of_openings"`
}

type circuitStats struct {
	numOfRuns      int64
	avgTime        float64
	lowestLatency  int64
	highestLatency int64
	numOfOpenings  int
}

func (c *circuitStats) updateRuns(delta int) {
	atomic.AddInt64(&c.numOfRuns, int64(delta))
}

func FillStats(name string, print bool) Stats {
	cb := getClipperWithName(name)
	total := cb.statistics.numOfRuns
	fails := cb.Failures
	var avg float64 = 100
	if fails != 0 {
		avg = float64(100 - (fails * 100 / total))
	}

	s := Stats{
		NumOfCircuits: len(clippers),
		IsOpen:        cb.open,
		TotalRuns:     total,
		TotalFails:    fails,
		AvgSuccess:    avg,
		NumOfOpenings: cb.statistics.numOfOpenings,
	}

	if print {
		printStats(s)
	}

	return s
}

func printStats(s Stats) {
	b, _ := json.MarshalIndent(s, "", "\t")
	fmt.Println(string(b))
}
