package truncator_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/truncator"
)

func makeResults(statuses ...differ.DiffStatus) []differ.Result {
	out := make([]differ.Result, len(statuses))
	for i, s := range statuses {
		out[i] = differ.Result{
			Type:   "listener",
			Name:   fmt.Sprintf("res-%d", i),
			Status: s,
		}
	}
	return out
}

import "fmt"

func TestApply_NoLimit(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusRemoved, differ.StatusModified)
	opts := truncator.DefaultOptions()
	opts.Limit = 0
	out := truncator.Apply(results, opts)
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestApply_KeepChanged_Truncates(t *testing.T) {
	results := makeResults(
		differ.StatusUnchanged, differ.StatusAdded, differ.StatusUnchanged,
		differ.StatusRemoved, differ.StatusModified, differ.StatusAdded,
	)
	opts := truncator.Options{Limit: 2, Strategy: truncator.KeepChanged}
	out := truncator.Apply(results, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	for _, r := range out {
		if r.Status == differ.StatusUnchanged {
			t.Errorf("unexpected unchanged result in output")
		}
	}
}

func TestApply_KeepFirst_Truncates(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusRemoved, differ.StatusModified, differ.StatusUnchanged)
	opts := truncator.Options{Limit: 2, Strategy: truncator.KeepFirst}
	out := truncator.Apply(results, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_UnderLimit_NoTruncation(t *testing.T) {
	results := makeResults(differ.StatusAdded, differ.StatusRemoved)
	opts := truncator.Options{Limit: 10, Strategy: truncator.KeepFirst}
	out := truncator.Apply(results, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestWrite_TruncationSummary(t *testing.T) {
	original := makeResults(differ.StatusAdded, differ.StatusRemoved, differ.StatusModified)
	truncated := original[:1]
	var buf bytes.Buffer
	truncator.Write(&buf, original, truncated)
	if !strings.Contains(buf.String(), "omitted 2") {
		t.Errorf("expected omission count in output, got: %s", buf.String())
	}
}

func TestWrite_NoTruncation(t *testing.T) {
	results := makeResults(differ.StatusAdded)
	var buf bytes.Buffer
	truncator.Write(&buf, results, results)
	if !strings.Contains(buf.String(), "no truncation") {
		t.Errorf("expected no-truncation message, got: %s", buf.String())
	}
}
