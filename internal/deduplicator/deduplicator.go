// Package deduplicator removes duplicate diff results based on configurable key functions.
package deduplicator

import "github.com/your-org/envoy-diff/internal/differ"

// Options controls deduplication behaviour.
type Options struct {
	// KeyFn derives a deduplication key from a result.
	// If nil, the default key (Type + "/" + Name) is used.
	KeyFn func(r differ.Result) string

	// KeepFirst retains the first occurrence when true;
	// otherwise the last occurrence is kept.
	KeepFirst bool
}

// DefaultOptions returns sensible deduplication defaults.
func DefaultOptions() Options {
	return Options{
		KeyFn:     defaultKey,
		KeepFirst: true,
	}
}

func defaultKey(r differ.Result) string {
	return r.Type + "/" + r.Name
}

// Apply deduplicates results using the provided options.
func Apply(results []differ.Result, opts Options) []differ.Result {
	if opts.KeyFn == nil {
		opts.KeyFn = defaultKey
	}

	seen := make(map[string]int, len(results)) // key -> index in out
	out := make([]differ.Result, 0, len(results))

	for _, r := range results {
		k := opts.KeyFn(r)
		if idx, exists := seen[k]; exists {
			if !opts.KeepFirst {
				out[idx] = r
			}
			continue
		}
		seen[k] = len(out)
		out = append(out, r)
	}

	return out
}

// Count returns the number of duplicates that would be removed.
func Count(results []differ.Result, opts Options) int {
	if opts.KeyFn == nil {
		opts.KeyFn = defaultKey
	}
	seen := make(map[string]struct{}, len(results))
	dupes := 0
	for _, r := range results {
		k := opts.KeyFn(r)
		if _, exists := seen[k]; exists {
			dupes++
			continue
		}
		seen[k] = struct{}{}
	}
	return dupes
}
