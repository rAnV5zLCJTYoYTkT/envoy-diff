// Package trimmer removes results whose diff body is below a minimum size
// threshold, helping to suppress noise from cosmetic-only changes.
package trimmer

import (
	"github.com/yourorg/envoy-diff/internal/differ"
)

// Options controls trimmer behaviour.
type Options struct {
	// MinBodyLen is the minimum number of characters the raw diff body must
	// contain for a result to be kept. Results with a body shorter than this
	// (including empty bodies) are dropped. Unchanged results are always kept
	// unless DropUnchanged is set.
	MinBodyLen int

	// DropUnchanged removes results whose status is Unchanged regardless of
	// body length.
	DropUnchanged bool
}

// DefaultOptions returns sensible defaults: drop results whose diff body is
// empty, and keep unchanged results.
func DefaultOptions() Options {
	return Options{
		MinBodyLen:    1,
		DropUnchanged: false,
	}
}

// Apply filters results according to opts, returning the surviving slice.
func Apply(results []differ.Result, opts Options) []differ.Result {
	out := make([]differ.Result, 0, len(results))
	for _, r := range results {
		if opts.DropUnchanged && r.Status == differ.Unchanged {
			continue
		}
		if r.Status != differ.Unchanged && len(r.RawDiff) < opts.MinBodyLen {
			continue
		}
		out = append(out, r)
	}
	return out
}

// Count returns the number of results that would be removed by Apply.
func Count(results []differ.Result, opts Options) int {
	return len(results) - len(Apply(results, opts))
}
