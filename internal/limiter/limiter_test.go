package limiter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envoy-diff/internal/differ"
	"github.com/envoy-diff/internal/limiter"
)

func makeResults(statuses ...differ.DiffStatus) []differ.DiffResult {
	out := make([]differ.DiffResult, len(statuses))
	for i, s := range statuses {
		out[i] = differ.DiffResult{
			Type:   "Cluster",
			Name:   fmt.Sprintf("resource-%d", i),
			Status: s,
		}
	}
	return out
}

import "fmt"

func TestApply_NoLimit(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusRemoved, differ.StatusUnchanged)
	opts := limiter.DefaultOptions()
	opts.MaxResults = 0
	r := limiter.Apply(results, opts)
	if r.Truncated {
		t.Error("expected not truncated")
	}
	if len(r.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(r.Items))
	}
}

func TestApply_Truncates(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusRemoved, differ.StatusModified, differ.StatusUnchanged)
	opts := limiter.Options{MaxResults: 2, OnlyChanged: false}
	r := limiter.Apply(results, opts)
	if !r.Truncated {
		t.Error("expected truncation")
	}
	if len(r.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(r.Items))
	}
	if r.Dropped != 2 {
		t.Errorf("expected 2 dropped, got %d", r.Dropped)
	}
}

func TestApply_OnlyChanged(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusUnchanged, differ.StatusUnchanged, differ.StatusRemoved)
	opts := limiter.Options{MaxResults: 10, OnlyChanged: true}
	r := limiter.Apply(results, opts)
	if r.Truncated {
		t.Error("expected not truncated")
	}
	if len(r.Items) != 2 {
		t.Errorf("expected 2 changed items, got %d", len(r.Items))
	}
}

func TestApply_OnlyChanged_WithCap(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusModified, differ.StatusRemoved, differ.StatusUnchanged)
	opts := limiter.Options{MaxResults: 2, OnlyChanged: true}
	r := limiter.Apply(results, opts)
	if !r.Truncated {
		t.Error("expected truncation after filtering")
	}
	if r.Total != 3 {
		t.Errorf("expected total 3 (changed), got %d", r.Total)
	}
	if r.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", r.Dropped)
	}
}

func TestWrite_Truncated(t *testing.T) {
	r := limiter.Result{Items: make([]differ.DiffResult, 5), Total: 10, Truncated: true, Dropped: 5}
	var buf bytes.Buffer
	limiter.Write(&buf, r)
	if !strings.Contains(buf.String(), "dropped") {
		t.Errorf("expected 'dropped' in output, got: %s", buf.String())
	}
}

func TestWrite_NotTruncated(t *testing.T) {
	r := limiter.Result{Items: make([]differ.DiffResult, 3), Total: 3}
	var buf bytes.Buffer
	limiter.Write(&buf, r)
	if !strings.Contains(buf.String(), "all 3") {
		t.Errorf("expected 'all 3' in output, got: %s", buf.String())
	}
}
