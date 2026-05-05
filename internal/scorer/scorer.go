// Package scorer assigns a numeric drift score to a set of diff results,
// providing a single comparable metric for snapshot divergence.
package scorer

import (
	"fmt"
	"io"

	"github.com/example/envoy-diff/internal/differ"
)

// Weights control how much each status contributes to the drift score.
const (
	WeightAdded    = 1
	WeightRemoved  = 2
	WeightModified = 3
	WeightUnchanged = 0
)

// Score holds the computed drift score and its breakdown.
type Score struct {
	Total     int            `json:"total"`
	Breakdown map[string]int `json:"breakdown"`
}

// Compute calculates a drift score from the provided diff results.
// Higher scores indicate greater divergence between snapshots.
func Compute(results []differ.Result) Score {
	breakdown := map[string]int{
		"added":     0,
		"removed":   0,
		"modified":  0,
		"unchanged": 0,
	}

	for _, r := range results {
		switch r.Status {
		case differ.StatusAdded:
			breakdown["added"] += WeightAdded
		case differ.StatusRemoved:
			breakdown["removed"] += WeightRemoved
		case differ.StatusModified:
			breakdown["modified"] += WeightModified
		case differ.StatusUnchanged:
			breakdown["unchanged"] += WeightUnchanged
		}
	}

	total := breakdown["added"] + breakdown["removed"] + breakdown["modified"]

	return Score{
		Total:     total,
		Breakdown: breakdown,
	}
}

// Write renders a human-readable score summary to w.
func Write(w io.Writer, s Score) {
	fmt.Fprintf(w, "Drift Score: %d\n", s.Total)
	fmt.Fprintf(w, "  added:     %d (x%d weight)\n", s.Breakdown["added"]/weightOrOne(WeightAdded), WeightAdded)
	fmt.Fprintf(w, "  removed:   %d (x%d weight)\n", s.Breakdown["removed"]/weightOrOne(WeightRemoved), WeightRemoved)
	fmt.Fprintf(w, "  modified:  %d (x%d weight)\n", s.Breakdown["modified"]/weightOrOne(WeightModified), WeightModified)
	fmt.Fprintf(w, "  unchanged: %d\n", s.Breakdown["unchanged"])
}

func weightOrOne(w int) int {
	if w == 0 {
		return 1
	}
	return w
}
