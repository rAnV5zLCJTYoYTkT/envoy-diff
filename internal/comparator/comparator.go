// Package comparator provides multi-environment diff comparison,
// allowing results from more than two snapshots to be compared pairwise.
package comparator

import (
	"fmt"
	"sort"

	"github.com/your-org/envoy-diff/internal/differ"
	"github.com/your-org/envoy-diff/internal/snapshot"
)

// Pair holds two environment labels and their diff results.
type Pair struct {
	Left  string
	Right string
	Diffs []differ.Result
}

// Report is the full output of a multi-environment comparison.
type Report struct {
	Pairs []Pair
}

// Options controls which environment pairs are compared.
type Options struct {
	// Sequential compares adjacent environments in order (A→B, B→C, …).
	// When false, all unique pairs are compared.
	Sequential bool
}

// DefaultOptions returns the default comparator options.
func DefaultOptions() Options {
	return Options{Sequential: true}
}

// Compare runs pairwise diffs across the provided named snapshots.
func Compare(envs map[string]*snapshot.Snapshot, opts Options) (Report, error) {
	if len(envs) < 2 {
		return Report{}, fmt.Errorf("comparator: at least two environments required, got %d", len(envs))
	}

	names := sortedKeys(envs)
	pairs := buildPairs(names, opts.Sequential)

	var report Report
	for _, p := range pairs {
		results := differ.Compare(envs[p[0]], envs[p[1]])
		report.Pairs = append(report.Pairs, Pair{
			Left:  p[0],
			Right: p[1],
			Diffs: results,
		})
	}
	return report, nil
}

func buildPairs(names []string, sequential bool) [][2]string {
	var out [][2]string
	if sequential {
		for i := 0; i+1 < len(names); i++ {
			out = append(out, [2]string{names[i], names[i+1]})
		}
		return out
	}
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			out = append(out, [2]string{names[i], names[j]})
		}
	}
	return out
}

func sortedKeys(m map[string]*snapshot.Snapshot) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
