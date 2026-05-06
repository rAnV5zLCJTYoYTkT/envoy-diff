// Package ranker assigns a numeric rank to diff results based on their
// severity and change status, enabling priority-ordered output.
package ranker

import (
	"sort"

	"github.com/your-org/envoy-diff/internal/differ"
)

// RankedResult wraps a differ.Result with a computed rank.
type RankedResult struct {
	differ.Result
	Rank int
}

// defaultWeights maps diff status strings to integer weights.
// Higher weight = higher priority / more severe.
var defaultWeights = map[string]int{
	"removed":   10,
	"modified":  7,
	"added":     5,
	"unchanged": 1,
}

// Options configures ranking behaviour.
type Options struct {
	// Weights overrides the default status-to-weight mapping.
	Weights map[string]int
	// Descending controls sort order; true = highest rank first.
	Descending bool
}

// DefaultOptions returns sensible ranking defaults.
func DefaultOptions() Options {
	return Options{
		Weights:    defaultWeights,
		Descending: true,
	}
}

// Rank assigns ranks to results and returns them sorted by rank.
func Rank(results []differ.Result, opts Options) []RankedResult {
	weights := opts.Weights
	if len(weights) == 0 {
		weights = defaultWeights
	}

	ranked := make([]RankedResult, len(results))
	for i, r := range results {
		w := weights[r.Status]
		if w == 0 {
			w = 1
		}
		ranked[i] = RankedResult{Result: r, Rank: w}
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if opts.Descending {
			return ranked[i].Rank > ranked[j].Rank
		}
		return ranked[i].Rank < ranked[j].Rank
	})

	return ranked
}

// Unwrap extracts the underlying differ.Result slice from ranked results.
func Unwrap(ranked []RankedResult) []differ.Result {
	out := make([]differ.Result, len(ranked))
	for i, r := range ranked {
		out[i] = r.Result
	}
	return out
}
