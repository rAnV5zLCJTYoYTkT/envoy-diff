// Package grouper groups diff results by a user-defined key function,
// enabling aggregation of DiffResults by resource type, status, or tag.
package grouper

import (
	"sort"

	"github.com/example/envoy-diff/internal/differ"
)

// KeyFunc extracts a grouping key from a DiffResult.
type KeyFunc func(r differ.DiffResult) string

// Group holds all DiffResults sharing the same key.
type Group struct {
	Key     string
	Results []differ.DiffResult
}

// ByType groups results by their resource type.
func ByType(r differ.DiffResult) string {
	return r.Type
}

// ByStatus groups results by their diff status string.
func ByStatus(r differ.DiffResult) string {
	return string(r.Status)
}

// Apply partitions results using keyFn and returns groups sorted by key.
func Apply(results []differ.DiffResult, keyFn KeyFunc) []Group {
	index := make(map[string]*Group)
	order := []string{}

	for _, r := range results {
		k := keyFn(r)
		if _, ok := index[k]; !ok {
			index[k] = &Group{Key: k}
			order = append(order, k)
		}
		index[k].Results = append(index[k].Results, r)
	}

	sort.Strings(order)

	groups := make([]Group, 0, len(order))
	for _, k := range order {
		groups = append(groups, *index[k])
	}
	return groups
}

// Counts returns a map of key -> count of results in each group.
func Counts(groups []Group) map[string]int {
	m := make(map[string]int, len(groups))
	for _, g := range groups {
		m[g.Key] = len(g.Results)
	}
	return m
}
