package indexer_test

import (
	"testing"

	"github.com/envoy-diff/internal/differ"
	"github.com/envoy-diff/internal/indexer"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "listener", Name: "listener-prod-a", Status: differ.Added},
		{Type: "listener", Name: "listener-prod-b", Status: differ.Removed},
		{Type: "cluster", Name: "cluster-prod-a", Status: differ.Modified},
		{Type: "cluster", Name: "cluster-staging-b", Status: differ.Unchanged},
		{Type: "route", Name: "route-prod-c", Status: differ.Added},
	}
}

func TestByType_ReturnsMatchingResults(t *testing.T) {
	idx := indexer.Build(makeResults())
	listeners := idx.ByType("listener")
	if len(listeners) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(listeners))
	}
	for _, r := range listeners {
		if r.Type != "listener" {
			t.Errorf("unexpected type %q", r.Type)
		}
	}
}

func TestByType_UnknownType_ReturnsEmpty(t *testing.T) {
	idx := indexer.Build(makeResults())
	if got := idx.ByType("endpoint"); len(got) != 0 {
		t.Errorf("expected empty, got %d results", len(got))
	}
}

func TestByStatus_ReturnsMatchingResults(t *testing.T) {
	idx := indexer.Build(makeResults())
	added := idx.ByStatus("added")
	if len(added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(added))
	}
}

func TestByName_ExactMatch(t *testing.T) {
	idx := indexer.Build(makeResults())
	r, ok := idx.ByName("cluster-prod-a")
	if !ok {
		t.Fatal("expected result to be found")
	}
	if r.Name != "cluster-prod-a" {
		t.Errorf("unexpected name %q", r.Name)
	}
}

func TestByName_Missing_ReturnsFalse(t *testing.T) {
	idx := indexer.Build(makeResults())
	_, ok := idx.ByName("does-not-exist")
	if ok {
		t.Error("expected not found")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	idx := indexer.Build(makeResults())
	hits := idx.Search("PROD")
	if len(hits) != 4 {
		t.Errorf("expected 4 prod results, got %d", len(hits))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	idx := indexer.Build(makeResults())
	if got := idx.Search("zzz"); len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}

func TestTypes_ReturnsDistinctTypes(t *testing.T) {
	idx := indexer.Build(makeResults())
	types := idx.Types()
	if len(types) != 3 {
		t.Errorf("expected 3 distinct types, got %d: %v", len(types), types)
	}
}

func TestAll_ReturnsAllResults(t *testing.T) {
	results := makeResults()
	idx := indexer.Build(results)
	if len(idx.All()) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(idx.All()))
	}
}
