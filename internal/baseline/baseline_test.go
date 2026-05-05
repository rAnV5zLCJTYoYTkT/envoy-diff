package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/envoy-diff/internal/baseline"
	"github.com/example/envoy-diff/internal/differ"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "cluster", Name: "api", Status: differ.StatusAdded},
		{Type: "cluster", Name: "db", Status: differ.StatusUnchanged},
		{Type: "listener", Name: "ingress", Status: differ.StatusModified},
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "base.json")
	results := makeResults()

	if err := baseline.Save(path, "test-label", results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	rec, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if rec.Label != "test-label" {
		t.Errorf("label: got %q, want %q", rec.Label, "test-label")
	}
	if len(rec.Results) != len(results) {
		t.Errorf("results len: got %d, want %d", len(rec.Results), len(results))
	}
	if rec.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad*.json")
	_, _ = f.WriteString("{not valid")
	f.Close()
	_, err := baseline.Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompare_DetectsDrift(t *testing.T) {
	base := &baseline.Record{
		CreatedAt: time.Now(),
		Results:   makeResults(),
	}
	newResults := []differ.DiffResult{
		{Type: "cluster", Name: "api", Status: differ.StatusAdded},   // same
		{Type: "cluster", Name: "db", Status: differ.StatusRemoved},  // status changed → drift
		{Type: "listener", Name: "egress", Status: differ.StatusAdded}, // new resource
	}
	drift := baseline.Compare(base, newResults)
	if len(drift) != 2 {
		t.Errorf("drift count: got %d, want 2", len(drift))
	}
}

func TestCompare_NoDrift(t *testing.T) {
	results := makeResults()
	base := &baseline.Record{CreatedAt: time.Now(), Results: results}
	drift := baseline.Compare(base, results)
	if len(drift) != 0 {
		t.Errorf("expected no drift, got %d", len(drift))
	}
}

func TestCompare_EmptyBase(t *testing.T) {
	base := &baseline.Record{CreatedAt: time.Now(), Results: nil}
	drift := baseline.Compare(base, makeResults())
	if len(drift) != len(makeResults()) {
		t.Errorf("all current results should be drift when base is empty")
	}
}
