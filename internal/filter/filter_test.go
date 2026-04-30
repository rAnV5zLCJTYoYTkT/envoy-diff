package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/filter"
)

func makeResults() []filter.Result {
	return []filter.Result{
		{Type: "listeners", Name: "ingress-http", Status: "added"},
		{Type: "listeners", Name: "egress-grpc", Status: "modified"},
		{Type: "clusters", Name: "backend-cluster", Status: "removed"},
		{Type: "clusters", Name: "sidecar-cluster", Status: "unchanged"},
		{Type: "routes", Name: "default-route", Status: "modified"},
	}
}

func TestApply_NoFilter(t *testing.T) {
	results := makeResults()
	out := filter.Apply(results, filter.Options{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_FilterByType(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{Types: []string{"clusters"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Type != "clusters" {
			t.Errorf("unexpected type %q", r.Type)
		}
	}
}

func TestApply_FilterByStatus(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{StatusInclude: []string{"modified"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 modified results, got %d", len(out))
	}
}

func TestApply_FilterByNameContains(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{NameContains: "grpc"})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Name != "egress-grpc" {
		t.Errorf("unexpected name %q", out[0].Name)
	}
}

func TestApply_FilterByNameContainsCaseInsensitive(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{NameContains: "CLUSTER"})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{
		Types:         []string{"listeners"},
		StatusInclude: []string{"added"},
	})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Name != "ingress-http" {
		t.Errorf("unexpected name %q", out[0].Name)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := filter.Apply(nil, filter.Options{Types: []string{"listeners"}})
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}
