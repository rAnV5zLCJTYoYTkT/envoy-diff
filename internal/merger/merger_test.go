package merger_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/merger"
)

func makeResult(rtype, name, status string) differ.Result {
	return differ.Result{
		ResourceType: rtype,
		Name:         name,
		Status:       differ.Status(status),
	}
}

func findResult(results []differ.Result, rtype, name string) (differ.Result, bool) {
	for _, r := range results {
		if r.ResourceType == rtype && r.Name == name {
			return r, true
		}
	}
	return differ.Result{}, false
}

func TestMerge_NoInputs(t *testing.T) {
	out := merger.Merge()
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d results", len(out))
	}
}

func TestMerge_SingleSlice(t *testing.T) {
	input := []differ.Result{
		makeResult("CDS", "cluster-a", "added"),
		makeResult("LDS", "listener-1", "unchanged"),
	}
	out := merger.Merge(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestMerge_DeduplicatesLowerPriority(t *testing.T) {
	a := []differ.Result{makeResult("CDS", "cluster-a", "unchanged")}
	b := []differ.Result{makeResult("CDS", "cluster-a", "modified")}

	out := merger.Merge(a, b)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	r, _ := findResult(out, "CDS", "cluster-a")
	if r.Status != "modified" {
		t.Errorf("expected status 'modified', got %q", r.Status)
	}
}

func TestMerge_KeepsDistinctEntries(t *testing.T) {
	a := []differ.Result{makeResult("CDS", "cluster-a", "added")}
	b := []differ.Result{makeResult("LDS", "listener-1", "removed")}

	out := merger.Merge(a, b)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestMerge_TieKeepsFirst(t *testing.T) {
	a := []differ.Result{makeResult("RDS", "route-x", "added")}
	b := []differ.Result{makeResult("RDS", "route-x", "added")}

	out := merger.Merge(a, b)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestMergeWithLabel_StampsTag(t *testing.T) {
	input := []differ.Result{makeResult("EDS", "endpoint-1", "unchanged")}
	out := merger.MergeWithLabel("env:staging", input)

	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Tags["merged_from"] != "env:staging" {
		t.Errorf("expected tag merged_from=env:staging, got %q", out[0].Tags["merged_from"])
	}
}
