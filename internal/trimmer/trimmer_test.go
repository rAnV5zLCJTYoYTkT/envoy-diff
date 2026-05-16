package trimmer_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/trimmer"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Name: "listener-a", Type: "listener", Status: differ.Added, RawDiff: `{"a":1}`},
		{Name: "cluster-b", Type: "cluster", Status: differ.Modified, RawDiff: `x`},
		{Name: "route-c", Type: "route", Status: differ.Modified, RawDiff: ""},
		{Name: "secret-d", Type: "secret", Status: differ.Removed, RawDiff: `{"key":"v"}`},
		{Name: "endpoint-e", Type: "endpoint", Status: differ.Unchanged, RawDiff: ""},
	}
}

func TestApply_DefaultOptions_DropsEmptyDiff(t *testing.T) {
	results := makeResults()
	opts := trimmer.DefaultOptions()
	out := trimmer.Apply(results, opts)
	for _, r := range out {
		if r.Status != differ.Unchanged && len(r.RawDiff) < opts.MinBodyLen {
			t.Errorf("result %q with empty diff was not trimmed", r.Name)
		}
	}
}

func TestApply_DefaultOptions_KeepsUnchanged(t *testing.T) {
	results := makeResults()
	out := trimmer.Apply(results, trimmer.DefaultOptions())
	found := false
	for _, r := range out {
		if r.Name == "endpoint-e" {
			found = true
		}
	}
	if !found {
		t.Error("expected unchanged result to be kept by default")
	}
}

func TestApply_DropUnchanged(t *testing.T) {
	results := makeResults()
	opts := trimmer.DefaultOptions()
	opts.DropUnchanged = true
	out := trimmer.Apply(results, opts)
	for _, r := range out {
		if r.Status == differ.Unchanged {
			t.Errorf("unchanged result %q should have been dropped", r.Name)
		}
	}
}

func TestApply_MinBodyLen_Threshold(t *testing.T) {
	results := makeResults()
	opts := trimmer.DefaultOptions()
	opts.MinBodyLen = 5
	out := trimmer.Apply(results, opts)
	for _, r := range out {
		if r.Status != differ.Unchanged && len(r.RawDiff) < 5 {
			t.Errorf("result %q has body shorter than MinBodyLen but was kept", r.Name)
		}
	}
}

func TestCount_MatchesDifference(t *testing.T) {
	results := makeResults()
	opts := trimmer.DefaultOptions()
	got := trimmer.Count(results, opts)
	want := len(results) - len(trimmer.Apply(results, opts))
	if got != want {
		t.Errorf("Count() = %d, want %d", got, want)
	}
}

func TestApply_Empty(t *testing.T) {
	out := trimmer.Apply(nil, trimmer.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d results", len(out))
	}
}
