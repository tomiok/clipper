package clipper

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type stats struct {
	NumOfCircuits int   `json:"num_of_circuits"`
	IsOpen        bool    `json:"is_open"`
	TotalRuns     int64   `json:"total_runs"`
	TotalFails    int64   `json:"total_fails"`
	AvgSuccess    float64 `json:"avg_success"`
}

type circuitStats struct {
	numOfRuns      int64
	avgTime        float64
	lowestLatency  int64
	highestLatency int64
}

func (c *circuitStats) updateRuns(delta int) {
	atomic.AddInt64(&c.numOfRuns, int64(delta))
}

func FillStats(name string, print bool) {
	cb := getClipper(name)
	total := cb.statistics.numOfRuns
	fails := cb.TotalFails
	var avg float64 = 100
	if fails != 0 {
		avg = float64(100 - (fails * 100 / total))
	}

	s := stats{
		NumOfCircuits: len(clippers),
		IsOpen:        cb.open,
		TotalRuns:     total,
		TotalFails:    fails,
		AvgSuccess:    avg,
	}

	if print {
		printStats(s)
	}
}

func printStats(s stats) {
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
}
