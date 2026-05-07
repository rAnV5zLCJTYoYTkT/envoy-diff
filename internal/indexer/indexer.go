// Package indexer builds a lookup index over diff results for fast retrieval
// by type, name, or status without repeated linear scans.
package indexer

import (
	"strings"

	"github.com/envoy-diff/internal/differ"
)

// Index holds pre-built lookup maps over a slice of diff results.
type Index struct {
	byType   map[string][]differ.Result
	byStatus map[string][]differ.Result
	byName   map[string]differ.Result
	all       []differ.Result
}

// Build constructs an Index from the provided results.
func Build(results []differ.Result) *Index {
	idx := &Index{
		byType:   make(map[string][]differ.Result),
		byStatus: make(map[string][]differ.Result),
		byName:   make(map[string]differ.Result),
		all:       results,
	}
	for _, r := range results {
		idx.byType[r.Type] = append(idx.byType[r.Type], r)
		idx.byStatus[string(r.Status)] = append(idx.byStatus[string(r.Status)], r)
		idx.byName[r.Name] = r
	}
	return idx
}

// ByType returns all results matching the given resource type.
func (idx *Index) ByType(resourceType string) []differ.Result {
	return idx.byType[resourceType]
}

// ByStatus returns all results matching the given status string.
func (idx *Index) ByStatus(status string) []differ.Result {
	return idx.byStatus[status]
}

// ByName returns the single result for an exact name match and whether it was found.
func (idx *Index) ByName(name string) (differ.Result, bool) {
	r, ok := idx.byName[name]
	return r, ok
}

// Search returns all results whose name contains the given substring (case-insensitive).
func (idx *Index) Search(substr string) []differ.Result {
	lower := strings.ToLower(substr)
	var out []differ.Result
	for _, r := range idx.all {
		if strings.Contains(strings.ToLower(r.Name), lower) {
			out = append(out, r)
		}
	}
	return out
}

// Types returns the distinct resource types present in the index.
func (idx *Index) Types() []string {
	keys := make([]string, 0, len(idx.byType))
	for k := range idx.byType {
		keys = append(keys, k)
	}
	return keys
}

// All returns every result held by the index.
func (idx *Index) All() []differ.Result {
	return idx.all
}
