package clipper

import (
	"encoding/json"
	"fmt"
)

type stats struct {
	NumOfCircuits int     `json:"num_of_circuits"`
	IsOpen        bool    `json:"is_open"`
	TotalRuns     int     `json:"total_runs"`
	TotalFails    int     `json:"total_fails"`
	AvgSuccess    float32 `json:"avg_success"`
	Path          string  `json:"path"`
}

func FillStats(name string, print bool) {
	cb := getClipper(name)
	total := cb.numOfRuns
	fails := cb.TotalFails
	var avg float32 = 100
	if fails != 0 {
		avg = 100 - float32((fails*100)/total)
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
