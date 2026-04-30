package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/reporter"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "cluster", Name: "c1", Status: differ.StatusAdded},
		{Type: "cluster", Name: "c2", Status: differ.StatusRemoved},
		{Type: "cluster", Name: "c3", Status: differ.StatusModified},
		{Type: "cluster", Name: "c4", Status: differ.StatusUnchanged},
		{Type: "listener", Name: "l1", Status: differ.StatusAdded},
		{Type: "listener", Name: "l2", Status: differ.StatusUnchanged},
	}
}

func TestCompute(t *testing.T) {
	results := makeResults()
	s := reporter.Compute(results)

	if s.Total != 6 {
		t.Errorf("Total: want 6, got %d", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("Added: want 2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("Modified: want 1, got %d", s.Modified)
	}
	if s.Unchanged != 2 {
		t.Errorf("Unchanged: want 2, got %d", s.Unchanged)
	}
}

func TestByType(t *testing.T) {
	results := makeResults()
	byType := reporter.ByType(results)

	if len(byType) != 2 {
		t.Fatalf("expected 2 types, got %d", len(byType))
	}
	cs := byType["cluster"]
	if cs.Total != 4 || cs.Added != 1 || cs.Removed != 1 || cs.Modified != 1 || cs.Unchanged != 1 {
		t.Errorf("cluster summary mismatch: %+v", cs)
	}
	ls := byType["listener"]
	if ls.Total != 2 || ls.Added != 1 || ls.Unchanged != 1 {
		t.Errorf("listener summary mismatch: %+v", ls)
	}
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	reporter.Write(&buf, makeResults())
	out := buf.String()

	for _, want := range []string{
		"Diff Summary",
		"Total: 6",
		"Added: 2",
		"[cluster]",
		"[listener]",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
	}
}
