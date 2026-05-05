package scorer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/scorer"
)

func makeResults(statuses ...differ.DiffStatus) []differ.Result {
	results := make([]differ.Result, 0, len(statuses))
	for i, s := range statuses {
		results = append(results, differ.Result{
			Type:   "cluster",
			Name:   fmt.Sprintf("resource-%d", i),
			Status: s,
		})
	}
	return results
}

func TestCompute_Empty(t *testing.T) {
	s := scorer.Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestCompute_OnlyAdded(t *testing.T) {
	results := []differ.Result{
		{Type: "cluster", Name: "a", Status: differ.StatusAdded},
		{Type: "cluster", Name: "b", Status: differ.StatusAdded},
	}
	s := scorer.Compute(results)
	// 2 added * weight 1 = 2
	if s.Total != 2 {
		t.Errorf("expected total 2, got %d", s.Total)
	}
	if s.Breakdown["added"] != 2 {
		t.Errorf("expected added breakdown 2, got %d", s.Breakdown["added"])
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := []differ.Result{
		{Type: "cluster", Name: "a", Status: differ.StatusAdded},
		{Type: "cluster", Name: "b", Status: differ.StatusRemoved},
		{Type: "cluster", Name: "c", Status: differ.StatusModified},
		{Type: "cluster", Name: "d", Status: differ.StatusUnchanged},
	}
	s := scorer.Compute(results)
	// 1*1 + 1*2 + 1*3 = 6
	if s.Total != 6 {
		t.Errorf("expected total 6, got %d", s.Total)
	}
}

func TestCompute_UnchangedDoesNotContribute(t *testing.T) {
	results := []differ.Result{
		{Type: "cluster", Name: "a", Status: differ.StatusUnchanged},
		{Type: "cluster", Name: "b", Status: differ.StatusUnchanged},
	}
	s := scorer.Compute(results)
	if s.Total != 0 {
		t.Errorf("expected total 0 for all unchanged, got %d", s.Total)
	}
}

func TestWrite_ContainsScore(t *testing.T) {
	results := []differ.Result{
		{Type: "listener", Name: "x", Status: differ.StatusModified},
	}
	s := scorer.Compute(results)
	var buf bytes.Buffer
	scorer.Write(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "Drift Score: 3") {
		t.Errorf("expected 'Drift Score: 3' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "modified") {
		t.Errorf("expected 'modified' in output, got:\n%s", out)
	}
}
