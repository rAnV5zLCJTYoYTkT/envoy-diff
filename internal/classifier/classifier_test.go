package classifier

import (
	"testing"

	"github.com/example/envoy-diff/internal/differ"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "listener", Name: "ingress", Status: "removed"},
		{Type: "cluster", Name: "backend-v1", Status: "removed"},
		{Type: "route", Name: "api-route", Status: "added"},
		{Type: "listener", Name: "egress", Status: "modified"},
		{Type: "endpoint", Name: "svc-a", Status: "modified"},
		{Type: "cluster", Name: "cache", Status: "unchanged"},
	}
}

func TestClassify_DefaultRules_RemovedListener(t *testing.T) {
	results := Classify(makeResults(), DefaultRules)
	for _, r := range results {
		if r.Type == "listener" && r.Status == "removed" {
			if r.Severity != SeverityCritical {
				t.Errorf("expected critical for removed listener, got %s", r.Severity)
			}
			return
		}
	}
	t.Error("removed listener result not found")
}

func TestClassify_DefaultRules_RemovedCluster(t *testing.T) {
	results := Classify(makeResults(), DefaultRules)
	for _, r := range results {
		if r.Type == "cluster" && r.Status == "removed" {
			if r.Severity != SeverityHigh {
				t.Errorf("expected high for removed cluster, got %s", r.Severity)
			}
			return
		}
	}
	t.Error("removed cluster result not found")
}

func TestClassify_DefaultRules_AddedRoute(t *testing.T) {
	results := Classify(makeResults(), DefaultRules)
	for _, r := range results {
		if r.Type == "route" && r.Status == "added" {
			if r.Severity != SeverityMedium {
				t.Errorf("expected medium for added route, got %s", r.Severity)
			}
			return
		}
	}
	t.Error("added route result not found")
}

func TestClassify_DefaultRules_Unchanged(t *testing.T) {
	results := Classify(makeResults(), DefaultRules)
	for _, r := range results {
		if r.Status == "unchanged" {
			if r.Severity != SeverityInfo {
				t.Errorf("expected info for unchanged, got %s", r.Severity)
			}
			return
		}
	}
	t.Error("unchanged result not found")
}

func TestClassify_CustomRule_NameContains(t *testing.T) {
	rules := []Rule{
		{NameContains: "ingress", Severity: SeverityCritical},
		{Severity: SeverityLow},
	}
	input := []differ.Result{
		{Type: "listener", Name: "ingress-main", Status: "modified"},
		{Type: "cluster", Name: "backend", Status: "modified"},
	}
	results := Classify(input, rules)
	if results[0].Severity != SeverityCritical {
		t.Errorf("expected critical for ingress-main, got %s", results[0].Severity)
	}
	if results[1].Severity != SeverityLow {
		t.Errorf("expected low for backend, got %s", results[1].Severity)
	}
}

func TestClassify_EmptyResults(t *testing.T) {
	results := Classify(nil, DefaultRules)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
