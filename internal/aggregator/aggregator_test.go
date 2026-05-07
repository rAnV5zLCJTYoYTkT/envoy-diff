package aggregator_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envoy-diff/internal/aggregator"
	"github.com/envoy-diff/internal/differ"
)

func makeResults(typ, status string, n int) []differ.Result {
	results := make([]differ.Result, n)
	for i := range results {
		results[i] = differ.Result{Type: typ, Status: status, Name: typ + "-" + status}
	}
	return results
}

func TestAggregate_Empty(t *testing.T) {
	results, summary := aggregator.Aggregate()
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
	if summary.Total != 0 {
		t.Fatalf("expected total 0, got %d", summary.Total)
	}
}

func TestAggregate_SingleSet(t *testing.T) {
	set := makeResults("cluster", "added", 3)
	results, summary := aggregator.Aggregate(set)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if summary.Total != 3 {
		t.Fatalf("expected total 3, got %d", summary.Total)
	}
	if summary.ByType["cluster"] != 3 {
		t.Fatalf("expected 3 clusters, got %d", summary.ByType["cluster"])
	}
	if summary.ByStatus["added"] != 3 {
		t.Fatalf("expected 3 added, got %d", summary.ByStatus["added"])
	}
}

func TestAggregate_MultipleSets(t *testing.T) {
	setA := makeResults("listener", "added", 2)
	setB := makeResults("cluster", "removed", 1)
	setC := makeResults("route", "modified", 4)

	results, summary := aggregator.Aggregate(setA, setB, setC)
	if len(results) != 7 {
		t.Fatalf("expected 7 results, got %d", len(results))
	}
	if summary.Total != 7 {
		t.Fatalf("expected total 7, got %d", summary.Total)
	}
	if summary.ByType["listener"] != 2 {
		t.Errorf("listener count mismatch")
	}
	if summary.ByType["cluster"] != 1 {
		t.Errorf("cluster count mismatch")
	}
	if summary.ByStatus["modified"] != 4 {
		t.Errorf("modified count mismatch")
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	set := makeResults("endpoint", "unchanged", 2)
	_, summary := aggregator.Aggregate(set)

	var buf bytes.Buffer
	aggregator.Write(&buf, summary)
	out := buf.String()

	for _, want := range []string{"Aggregated", "By Type", "By Status", "endpoint", "unchanged"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}
