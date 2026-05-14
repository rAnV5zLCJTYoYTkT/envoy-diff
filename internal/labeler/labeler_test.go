package labeler_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/labeler"
)

func makeResults(specs []struct {
	name   string
	typ    string
	status differ.DiffStatus
}) []differ.Result {
	out := make([]differ.Result, len(specs))
	for i, s := range specs {
		out[i] = differ.Result{Name: s.name, Type: s.typ, Status: s.status}
	}
	return out
}

func findLabel(results []differ.Result, name string) string {
	for _, r := range results {
		if r.Name == name && r.Tags != nil {
			return r.Tags["label"]
		}
	}
	return ""
}

func TestApply_DefaultRules_ListenerChange(t *testing.T) {
	results := makeResults([]struct {
		name   string
		typ    string
		status differ.DiffStatus
	}{
		{"ingress", "listener", differ.StatusModified},
	})
	l := labeler.New()
	out := l.Apply(results)
	if got := findLabel(out, "ingress"); got != "listener-change" {
		t.Errorf("expected listener-change, got %q", got)
	}
}

func TestApply_DefaultRules_AddedResource(t *testing.T) {
	results := makeResults([]struct {
		name   string
		typ    string
		status differ.DiffStatus
	}{
		{"new-cluster", "endpoint", differ.StatusAdded},
	})
	l := labeler.New()
	out := l.Apply(results)
	if got := findLabel(out, "new-cluster"); got != "added" {
		t.Errorf("expected added, got %q", got)
	}
}

func TestApply_DefaultRules_Unchanged_NoLabel(t *testing.T) {
	results := makeResults([]struct {
		name   string
		typ    string
		status differ.DiffStatus
	}{
		{"stable", "cluster", differ.StatusUnchanged},
	})
	l := labeler.New()
	out := l.Apply(results)
	if got := findLabel(out, "stable"); got != "" {
		t.Errorf("expected no label for unchanged, got %q", got)
	}
}

func TestApply_CustomRule_Overrides(t *testing.T) {
	customRule := labeler.Rule{
		Label: "critical",
		Predicate: func(r differ.Result) bool {
			return r.Name == "prod-listener"
		},
	}
	results := makeResults([]struct {
		name   string
		typ    string
		status differ.DiffStatus
	}{
		{"prod-listener", "listener", differ.StatusRemoved},
	})
	l := labeler.NewWithOptions(labeler.WithRules([]labeler.Rule{customRule}))
	out := l.Apply(results)
	if got := findLabel(out, "prod-listener"); got != "critical" {
		t.Errorf("expected critical, got %q", got)
	}
}

func TestApply_WithExtraRules_AppendsBehavior(t *testing.T) {
	extra := labeler.Rule{
		Label: "secret-route",
		Predicate: func(r differ.Result) bool {
			return r.Type == "route" && r.Status == differ.StatusAdded
		},
	}
	results := makeResults([]struct {
		name   string
		typ    string
		status differ.DiffStatus
	}{
		{"new-route", "route", differ.StatusAdded},
	})
	// route-change rule fires first for modified; added route should match route-change
	// but we replace rules entirely to test extra rule path
	l := labeler.NewWithOptions(
		labeler.WithRules([]labeler.Rule{}),
		labeler.WithExtraRules(extra),
	)
	out := l.Apply(results)
	if got := findLabel(out, "new-route"); got != "secret-route" {
		t.Errorf("expected secret-route, got %q", got)
	}
}

func TestApply_EmptyResults(t *testing.T) {
	l := labeler.New()
	out := l.Apply(nil)
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d results", len(out))
	}
}
