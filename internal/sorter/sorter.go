// Package sorter provides utilities for ordering diff results
// by various criteria such as status, type, name, or score.
package sorter

import (
	"sort"

	"github.com/your-org/envoy-diff/internal/differ"
)

// SortBy defines the field used for ordering results.
type SortBy int

const (
	ByName SortBy = iota
	ByType
	ByStatus
)

// Options controls sorting behaviour.
type Options struct {
	Field     SortBy
	Ascending bool
}

// DefaultOptions returns sensible defaults: sort by type ascending.
func DefaultOptions() Options {
	return Options{Field: ByType, Ascending: true}
}

// Apply sorts a copy of results according to opts and returns it.
func Apply(results []differ.DiffResult, opts Options) []differ.DiffResult {
	out := make([]differ.DiffResult, len(results))
	copy(out, results)

	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch opts.Field {
		case ByName:
			less = out[i].Name < out[j].Name
		case ByStatus:
			less = string(out[i].Status) < string(out[j].Status)
		default: // ByType
			less = out[i].Type < out[j].Type
		}
		if opts.Ascending {
			return less
		}
		return !less
	})
	return out
}

// ByField is a convenience wrapper that applies sorting with default ascending order.
func ByField(results []differ.DiffResult, field SortBy) []differ.DiffResult {
	return Apply(results, Options{Field: field, Ascending: true})
}
