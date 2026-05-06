package ranker_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/differ"
	"github.com/your-org/envoy-diff/internal/ranker"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "cluster", Name: "c1", Status: "unchanged"},
		{Type: "cluster", Name: "c2", Status: "added"},
		{Type: "listener", Name: "l1", Status: "removed"},
		{Type: "route", Name: "r1", Status: "modified"},
	}
}

func TestRank_DescendingOrder(t *testing.T) {
	results := makeResults()
	ranked := ranker.Rank(results, ranker.DefaultOptions())

	if len(ranked) != 4 {
		t.Fatalf("expected 4 results, got %d", len(ranked))
	}

	// First result should be the highest rank (removed = 10).
	if ranked[0].Status != "removed" {
		t.Errorf("expected first status=removed, got %s", ranked[0].Status)
	}
	if ranked[0].Rank != 10 {
		t.Errorf("expected rank 10, got %d", ranked[0].Rank)
	}

	// Last result should be unchanged (rank 1).
	if ranked[len(ranked)-1].Status != "unchanged" {
		t.Errorf("expected last status=unchanged, got %s", ranked[len(ranked)-1].Status)
	}
}

func TestRank_AscendingOrder(t *testing.T) {
	results := makeResults()
	opts := ranker.DefaultOptions()
	opts.Descending = false
	ranked := ranker.Rank(results, opts)

	if ranked[0].Status != "unchanged" {
		t.Errorf("expected first status=unchanged in ascending, got %s", ranked[0].Status)
	}
}

func TestRank_CustomWeights(t *testing.T) {
	results := makeResults()
	opts := ranker.Options{
		Weights:    map[string]int{"added": 99, "removed": 1, "modified": 1, "unchanged": 1},
		Descending: true,
	}
	ranked := ranker.Rank(results, opts)

	if ranked[0].Status != "added" {
		t.Errorf("expected added to rank first with custom weights, got %s", ranked[0].Status)
	}
}

func TestRank_UnknownStatus_DefaultsToOne(t *testing.T) {
	results := []differ.Result{
		{Type: "endpoint", Name: "e1", Status: "unknown-status"},
	}
	ranked := ranker.Rank(results, ranker.DefaultOptions())
	if ranked[0].Rank != 1 {
		t.Errorf("expected rank 1 for unknown status, got %d", ranked[0].Rank)
	}
}

func TestUnwrap_PreservesOrder(t *testing.T) {
	results := makeResults()
	ranked := ranker.Rank(results, ranker.DefaultOptions())
	unwrapped := ranker.Unwrap(ranked)

	if len(unwrapped) != len(results) {
		t.Fatalf("expected %d results after unwrap, got %d", len(results), len(unwrapped))
	}
	// First unwrapped element should match first ranked element.
	if unwrapped[0].Name != ranked[0].Name {
		t.Errorf("unwrap order mismatch: got %s, want %s", unwrapped[0].Name, ranked[0].Name)
	}
}
