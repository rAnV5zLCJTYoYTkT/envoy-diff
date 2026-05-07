package grouper_test

import (
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/grouper"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "listener", Name: "l1", Status: differ.Added},
		{Type: "listener", Name: "l2", Status: differ.Removed},
		{Type: "cluster", Name: "c1", Status: differ.Modified},
		{Type: "cluster", Name: "c2", Status: differ.Unchanged},
		{Type: "route", Name: "r1", Status: differ.Added},
	}
}

func TestApply_ByType(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.ByType)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	// groups are sorted by key
	if groups[0].Key != "cluster" {
		t.Errorf("expected first group 'cluster', got %q", groups[0].Key)
	}
	if len(groups[0].Results) != 2 {
		t.Errorf("expected 2 cluster results, got %d", len(groups[0].Results))
	}
}

func TestApply_ByStatus(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.ByStatus)
	counts := grouper.Counts(groups)
	if counts["added"] != 2 {
		t.Errorf("expected 2 added, got %d", counts["added"])
	}
	if counts["removed"] != 1 {
		t.Errorf("expected 1 removed, got %d", counts["removed"])
	}
	if counts["modified"] != 1 {
		t.Errorf("expected 1 modified, got %d", counts["modified"])
	}
}

func TestApply_CustomKeyFn(t *testing.T) {
	// group by first character of name
	keyFn := func(r differ.DiffResult) string {
		if len(r.Name) == 0 {
			return "unknown"
		}
		return string(r.Name[0])
	}
	groups := grouper.Apply(makeResults(), keyFn)
	counts := grouper.Counts(groups)
	if counts["l"] != 2 {
		t.Errorf("expected 2 in group 'l', got %d", counts["l"])
	}
	if counts["c"] != 2 {
		t.Errorf("expected 2 in group 'c', got %d", counts["c"])
	}
	if counts["r"] != 1 {
		t.Errorf("expected 1 in group 'r', got %d", counts["r"])
	}
}

func TestApply_Empty(t *testing.T) {
	groups := grouper.Apply(nil, grouper.ByType)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestCounts_Empty(t *testing.T) {
	counts := grouper.Counts(nil)
	if len(counts) != 0 {
		t.Errorf("expected empty counts map, got %v", counts)
	}
}

func TestApply_SortedKeys(t *testing.T) {
	groups := grouper.Apply(makeResults(), grouper.ByType)
	keys := make([]string, len(groups))
	for i, g := range groups {
		keys[i] = g.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("groups not sorted: %v", keys)
		}
	}
}

// TestApply_ResultsPreserveOrder verifies that results within each group
// maintain the same relative order as they appeared in the input slice.
func TestApply_ResultsPreserveOrder(t *testing.T) {
	results := []differ.DiffResult{
		{Type: "cluster", Name: "c2", Status: differ.Unchanged},
		{Type: "cluster", Name: "c1", Status: differ.Modified},
		{Type: "cluster", Name: "c3", Status: differ.Added},
	}
	groups := grouper.Apply(results, grouper.ByType)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	got := groups[0].Results
	expected := []string{"c2", "c1", "c3"}
	for i, name := range expected {
		if got[i].Name != name {
			t.Errorf("position %d: expected name %q, got %q", i, name, got[i].Name)
		}
	}
}
