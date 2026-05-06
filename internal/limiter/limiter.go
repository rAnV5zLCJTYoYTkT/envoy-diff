// Package limiter provides result set truncation with configurable caps
// and summary reporting for envoy-diff diff results.
package limiter

import (
	"fmt"
	"io"

	"github.com/envoy-diff/internal/differ"
)

// Options controls how the limiter behaves.
type Options struct {
	// MaxResults is the maximum number of results to keep. 0 means unlimited.
	MaxResults int
	// OnlyChanged, when true, excludes unchanged results before applying the cap.
	OnlyChanged bool
}

// DefaultOptions returns sensible defaults: cap at 100, include all statuses.
func DefaultOptions() Options {
	return Options{
		MaxResults:  100,
		OnlyChanged: false,
	}
}

// Result holds the truncated slice and metadata about the operation.
type Result struct {
	Items     []differ.DiffResult
	Total     int
	Truncated bool
	Dropped   int
}

// Apply truncates results according to opts and returns a Result.
func Apply(results []differ.DiffResult, opts Options) Result {
	working := results
	if opts.OnlyChanged {
		working = filterChanged(working)
	}

	total := len(working)
	if opts.MaxResults <= 0 || total <= opts.MaxResults {
		return Result{Items: working, Total: total}
	}

	return Result{
		Items:     working[:opts.MaxResults],
		Total:     total,
		Truncated: true,
		Dropped:   total - opts.MaxResults,
	}
}

// Write prints a human-readable summary of the limit result to w.
func Write(w io.Writer, r Result) {
	if r.Truncated {
		fmt.Fprintf(w, "Showing %d of %d results (%d dropped due to limit)\n",
			len(r.Items), r.Total, r.Dropped)
	} else {
		fmt.Fprintf(w, "Showing all %d results\n", r.Total)
	}
}

func filterChanged(results []differ.DiffResult) []differ.DiffResult {
	out := make([]differ.DiffResult, 0, len(results))
	for _, r := range results {
		if r.Status != differ.StatusUnchanged {
			out = append(out, r)
		}
	}
	return out
}
