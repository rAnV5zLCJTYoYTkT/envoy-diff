// Package normalizer provides utilities for normalizing differ.Result slices
// before comparison or display — trimming whitespace from names, lowercasing
// resource types, and canonicalizing JSON payloads.
package normalizer

import (
	"encoding/json"
	"strings"

	"github.com/example/envoy-diff/internal/differ"
)

// Options controls which normalization steps are applied.
type Options struct {
	LowercaseType bool
	TrimName      bool
	CanonicalJSON bool
}

// DefaultOptions returns the recommended normalization settings.
func DefaultOptions() Options {
	return Options{
		LowercaseType: true,
		TrimName:      true,
		CanonicalJSON: true,
	}
}

// Apply normalizes a slice of differ.Result values according to opts.
// It returns a new slice; the originals are not mutated.
func Apply(results []differ.Result, opts Options) []differ.Result {
	out := make([]differ.Result, 0, len(results))
	for _, r := range results {
		nr := r
		if opts.TrimName {
			nr.Name = strings.TrimSpace(nr.Name)
		}
		if opts.LowercaseType {
			nr.Type = strings.ToLower(nr.Type)
		}
		if opts.CanonicalJSON {
			nr.Left = canonicalize(nr.Left)
			nr.Right = canonicalize(nr.Right)
		}
		out = append(out, nr)
	}
	return out
}

// canonicalize round-trips a JSON string through encoding/json so that keys
// are sorted and unnecessary whitespace is removed. If the input is not valid
// JSON it is returned unchanged.
func canonicalize(s string) string {
	if s == "" {
		return s
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return s
	}
	return string(b)
}
