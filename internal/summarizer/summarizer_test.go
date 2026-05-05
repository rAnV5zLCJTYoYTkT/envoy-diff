package summarizer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/summarizer"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "cluster", Name: "c1", Status: differ.StatusAdded},
		{Type: "cluster", Name: "c2", Status: differ.StatusRemoved},
		{Type: "cluster", Name: "c3", Status: differ.StatusModified},
		{Type: "listener", Name: "l1", Status: differ.StatusUnchanged},
		{Type: "listener", Name: "l2", Status: differ.StatusAdded},
	}
}

func TestCompute_Totals(t *testing.T) {
	s := summarizer.Compute(makeResults())
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", s.Unchanged)
	}
}

func TestCompute_ByType(t *testing.T) {
	s := summarizer.Compute(makeResults())
	cluster, ok := s.ByType["cluster"]
	if !ok {
		t.Fatal("expected cluster type in ByType")
	}
	if cluster.Added != 1 || cluster.Removed != 1 || cluster.Modified != 1 {
		t.Errorf("unexpected cluster counts: %+v", cluster)
	}
	listener := s.ByType["listener"]
	if listener.Added != 1 || listener.Unchanged != 1 {
		t.Errorf("unexpected listener counts: %+v", listener)
	}
}

func TestCompute_Empty(t *testing.T) {
	s := summarizer.Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0 for empty input")
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	s := summarizer.Compute(makeResults())
	var buf bytes.Buffer
	if err := summarizer.Write(&buf, s); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Diff Summary") {
		t.Error("expected 'Diff Summary' in output")
	}
	if !strings.Contains(out, "cluster") {
		t.Error("expected 'cluster' in output")
	}
	if !strings.Contains(out, "listener") {
		t.Error("expected 'listener' in output")
	}
}

func TestWrite_Empty(t *testing.T) {
	s := summarizer.Compute(nil)
	var buf bytes.Buffer
	if err := summarizer.Write(&buf, s); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 total") {
		t.Error("expected '0 total' in empty summary output")
	}
}
