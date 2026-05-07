// Package aggregator combines multiple diff result slices into a single
// unified view, computing per-type and per-status roll-ups.
package aggregator

import (
	"fmt"
	"io"
	"sort"

	"github.com/envoy-diff/internal/differ"
)

// Summary holds aggregated counts across all provided result sets.
type Summary struct {
	Total    int
	ByType   map[string]int
	ByStatus map[string]int
}

// Aggregate merges one or more result slices and returns the combined slice
// together with a roll-up Summary.
func Aggregate(sets ...[]differ.Result) ([]differ.Result, Summary) {
	var combined []differ.Result
	for _, s := range sets {
		combined = append(combined, s...)
	}

	summary := Summary{
		Total:    len(combined),
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	for _, r := range combined {
		summary.ByType[r.Type]++
		summary.ByStatus[r.Status]++
	}

	return combined, summary
}

// Write prints a human-readable aggregation summary to w.
func Write(w io.Writer, s Summary) {
	fmt.Fprintf(w, "Aggregated %d result(s)\n", s.Total)

	fmt.Fprintln(w, "\nBy Type:")
	types := sortedKeys(s.ByType)
	for _, t := range types {
		fmt.Fprintf(w, "  %-30s %d\n", t, s.ByType[t])
	}

	fmt.Fprintln(w, "\nBy Status:")
	statuses := sortedKeys(s.ByStatus)
	for _, st := range statuses {
		fmt.Fprintf(w, "  %-12s %d\n", st, s.ByStatus[st])
	}
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
