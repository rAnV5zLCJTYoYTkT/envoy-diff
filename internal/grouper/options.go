package grouper

import "github.com/example/envoy-diff/internal/differ"

// Options configures grouper behaviour.
type Options struct {
	KeyFn    KeyFunc
	MinCount int // groups with fewer results than MinCount are omitted (0 = keep all)
}

// DefaultOptions returns Options using ByType with no minimum count.
func DefaultOptions() Options {
	return Options{KeyFn: ByType, MinCount: 0}
}

// ApplyWithOptions partitions results according to opts, filtering out
// groups that contain fewer than opts.MinCount results.
func ApplyWithOptions(results []differ.DiffResult, opts Options) []Group {
	keyFn := opts.KeyFn
	if keyFn == nil {
		keyFn = ByType
	}

	all := Apply(results, keyFn)

	if opts.MinCount <= 0 {
		return all
	}

	filtered := all[:0]
	for _, g := range all {
		if len(g.Results) >= opts.MinCount {
			filtered = append(filtered, g)
		}
	}
	return filtered
}
