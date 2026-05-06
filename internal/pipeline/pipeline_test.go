package pipeline_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/filter"
	"github.com/yourorg/envoy-diff/internal/pipeline"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "cluster", Name: "cluster-a", Status: differ.StatusAdded},
		{Type: "listener", Name: "listener-b", Status: differ.StatusRemoved},
		{Type: "cluster", Name: "cluster-c", Status: differ.StatusUnchanged},
	}
}

func TestNew_EmptyPipeline(t *testing.T) {
	p := pipeline.New()
	results := makeResults()
	out := p.Run(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestNew_CustomStage(t *testing.T) {
	called := false
	stage := pipeline.Stage(func(r []differ.DiffResult) []differ.DiffResult {
		called = true
		return r
	})
	p := pipeline.New(stage)
	p.Run(makeResults())
	if !called {
		t.Fatal("expected custom stage to be called")
	}
}

func TestFilterStage_ReducesResults(t *testing.T) {
	opts := filter.Options{Types: []string{"cluster"}}
	p := pipeline.New(pipeline.FilterStage(opts))
	out := p.Run(makeResults())
	for _, r := range out {
		if r.Type != "cluster" {
			t.Errorf("unexpected type %q after filter", r.Type)
		}
	}
}

func TestDefaultPipeline_RunsWithoutPanic(t *testing.T) {
	p := pipeline.DefaultPipeline()
	out := p.Run(makeResults())
	if out == nil {
		t.Fatal("expected non-nil result from default pipeline")
	}
}

func TestNewWithOptions_SkipFilter(t *testing.T) {
	opts := pipeline.Options{
		Filter:     filter.Options{Types: []string{"cluster"}},
		SkipFilter: true,
	}
	p := pipeline.NewWithOptions(opts)
	out := p.Run(makeResults())
	// All three results should survive because filter is skipped.
	if len(out) != 3 {
		t.Fatalf("expected 3 results with SkipFilter, got %d", len(out))
	}
}

func TestNewWithOptions_WithFilter(t *testing.T) {
	opts := pipeline.Options{
		Filter: filter.Options{Types: []string{"listener"}},
	}
	p := pipeline.NewWithOptions(opts)
	out := p.Run(makeResults())
	for _, r := range out {
		if r.Type != "listener" {
			t.Errorf("unexpected type %q; expected only listener", r.Type)
		}
	}
}
