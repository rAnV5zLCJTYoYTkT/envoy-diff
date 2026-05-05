// Package summarizer provides a high-level summary of diff results,
// aggregating counts and highlights across resource types and statuses.
package summarizer

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/example/envoy-diff/internal/differ"
)

// Summary holds aggregated diff statistics.
type Summary struct {
	Total    int
	Added    int
	Removed  int
	Modified int
	Unchanged int
	ByType   map[string]TypeSummary
}

// TypeSummary holds per-resource-type counts.
type TypeSummary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
}

// Compute builds a Summary from a slice of DiffResults.
func Compute(results []differ.DiffResult) Summary {
	s := Summary{
		ByType: make(map[string]TypeSummary),
	}
	for _, r := range results {
		s.Total++
		ts := s.ByType[r.Type]
		switch r.Status {
		case differ.StatusAdded:
			s.Added++
			ts.Added++
		case differ.StatusRemoved:
			s.Removed++
			ts.Removed++
		case differ.StatusModified:
			s.Modified++
			ts.Modified++
		case differ.StatusUnchanged:
			s.Unchanged++
			ts.Unchanged++
		}
		s.ByType[r.Type] = ts
	}
	return s
}

// Write renders the Summary as a human-readable table to w.
func Write(w io.Writer, s Summary) error {
	fmt.Fprintf(w, "Diff Summary: %d total resource(s)\n", s.Total)
	fmt.Fprintf(w, "  Added: %d | Removed: %d | Modified: %d | Unchanged: %d\n\n",
		s.Added, s.Removed, s.Modified, s.Unchanged)

	if len(s.ByType) == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tADDED\tREMOVED\tMODIFIED\tUNCHANGED")

	types := make([]string, 0, len(s.ByType))
	for t := range s.ByType {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		ts := s.ByType[t]
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\n", t, ts.Added, ts.Removed, ts.Modified, ts.Unchanged)
	}
	return tw.Flush()
}
