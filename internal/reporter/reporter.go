// Package reporter provides summary reporting for xDS diff results.
package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envoy-diff/internal/differ"
)

// Summary holds aggregated counts of diff results.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
	Total    int
}

// Compute builds a Summary from a slice of DiffResults.
func Compute(results []differ.DiffResult) Summary {
	var s Summary
	for _, r := range results {
		s.Total++
		switch r.Status {
		case differ.StatusAdded:
			s.Added++
		case differ.StatusRemoved:
			s.Removed++
		case differ.StatusModified:
			s.Modified++
		case differ.StatusUnchanged:
			s.Unchanged++
		}
	}
	return s
}

// ByType returns a map of resource type to Summary for the given results.
func ByType(results []differ.DiffResult) map[string]Summary {
	groups := make(map[string][]differ.DiffResult)
	for _, r := range results {
		groups[r.Type] = append(groups[r.Type], r)
	}
	out := make(map[string]Summary, len(groups))
	for t, rs := range groups {
		out[t] = Compute(rs)
	}
	return out
}

// Write writes a human-readable summary report to w.
func Write(w io.Writer, results []differ.DiffResult) {
	overall := Compute(results)
	fmt.Fprintf(w, "=== Diff Summary ===\n")
	fmt.Fprintf(w, "Total: %d  Added: %d  Removed: %d  Modified: %d  Unchanged: %d\n\n",
		overall.Total, overall.Added, overall.Removed, overall.Modified, overall.Unchanged)

	byType := ByType(results)
	types := make([]string, 0, len(byType))
	for t := range byType {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		s := byType[t]
		fmt.Fprintf(w, "  [%s] total=%d added=%d removed=%d modified=%d unchanged=%d\n",
			t, s.Total, s.Added, s.Removed, s.Modified, s.Unchanged)
	}
}
