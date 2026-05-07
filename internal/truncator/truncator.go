// Package truncator limits the number of diff results returned, with
// configurable strategies for which results to prioritize when trimming.
package truncator

import (
	"fmt"
	"io"

	"github.com/yourorg/envoy-diff/internal/differ"
)

// Strategy controls which results are kept when the limit is exceeded.
type Strategy int

const (
	// KeepChanged retains only changed results (added, removed, modified) up to Limit.
	KeepChanged Strategy = iota
	// KeepFirst retains the first N results regardless of status.
	KeepFirst
)

// Options configures the Truncator.
type Options struct {
	Limit    int
	Strategy Strategy
}

// DefaultOptions returns sensible defaults: keep first 100 changed results.
func DefaultOptions() Options {
	return Options{
		Limit:    100,
		Strategy: KeepChanged,
	}
}

// Apply truncates results according to opts. If Limit <= 0, results are returned as-is.
func Apply(results []differ.Result, opts Options) []differ.Result {
	if opts.Limit <= 0 || len(results) <= opts.Limit {
		return results
	}

	switch opts.Strategy {
	case KeepChanged:
		var changed []differ.Result
		for _, r := range results {
			if r.Status != differ.StatusUnchanged {
				changed = append(changed, r)
				if len(changed) >= opts.Limit {
					break
				}
			}
		}
		return changed
	case KeepFirst:
		return results[:opts.Limit]
	default:
		return results[:opts.Limit]
	}
}

// Write prints a summary of truncation to w.
func Write(w io.Writer, original, truncated []differ.Result) {
	if len(original) == len(truncated) {
		fmt.Fprintf(w, "truncator: no truncation applied (%d results)\n", len(original))
		return
	}
	fmt.Fprintf(w, "truncator: showing %d of %d results (omitted %d)\n",
		len(truncated), len(original), len(original)-len(truncated))
}
