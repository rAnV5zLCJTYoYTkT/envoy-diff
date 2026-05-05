package patcher_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/patcher"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "cluster", Name: "healthcheck-cluster", Status: "modified"},
		{Type: "listener", Name: "public-listener", Status: "added"},
		{Type: "cluster", Name: "stats-cluster", Status: "removed"},
		{Type: "route", Name: "api-route", Status: "unchanged"},
	}
}

func TestApply_NoRules(t *testing.T) {
	p := patcher.New(nil)
	results := makeResults()
	out := p.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
	if out[0].Status != "modified" {
		t.Errorf("expected status modified, got %s", out[0].Status)
	}
}

func TestApply_SuppressByNameContains(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{NameContains: "healthcheck", SuppressStatus: true},
	})
	out := p.Apply(makeResults())
	if out[0].Status != "suppressed" {
		t.Errorf("expected suppressed, got %s", out[0].Status)
	}
	if out[1].Status != "added" {
		t.Errorf("expected added unchanged, got %s", out[1].Status)
	}
}

func TestApply_SuppressByType(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Type: "listener", SuppressStatus: true},
	})
	out := p.Apply(makeResults())
	if out[1].Status != "suppressed" {
		t.Errorf("expected listener suppressed, got %s", out[1].Status)
	}
	if out[0].Status == "suppressed" {
		t.Error("cluster should not be suppressed")
	}
}

func TestApply_OverrideLabel(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{NameContains: "api-route", OverrideLabel: "known-safe"},
	})
	out := p.Apply(makeResults())
	if out[3].Tags["label"] != "known-safe" {
		t.Errorf("expected label 'known-safe', got %v", out[3].Tags)
	}
}

func TestApply_DefaultSuppressRules(t *testing.T) {
	p := patcher.New(patcher.DefaultSuppressRules())
	out := p.Apply(makeResults())
	// healthcheck-cluster and stats-cluster should be suppressed
	for _, r := range out {
		if (r.Name == "healthcheck-cluster" || r.Name == "stats-cluster") && r.Status != "suppressed" {
			t.Errorf("%s should be suppressed, got %s", r.Name, r.Status)
		}
	}
}

func TestNewWithOptions(t *testing.T) {
	p := patcher.NewWithOptions(
		patcher.WithRules(patcher.Rule{NameContains: "public", SuppressStatus: true}),
	)
	out := p.Apply(makeResults())
	if out[1].Status != "suppressed" {
		t.Errorf("expected public-listener suppressed, got %s", out[1].Status)
	}
}
