package validator_test

import (
	"strings"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/validator"
)

func makeResults(statuses ...differ.Status) []differ.Result {
	results := make([]differ.Result, len(statuses))
	for i, s := range statuses {
		results[i] = differ.Result{
			Type:   "Cluster",
			Name:   fmt.Sprintf("res-%d", i),
			Status: s,
		}
	}
	return results
}

func TestValidate_NoViolations(t *testing.T) {
	results := []differ.Result{
		{Type: "Cluster", Name: "a", Status: differ.Unchanged},
		{Type: "Cluster", Name: "b", Status: differ.Added},
	}
	rules := validator.DefaultRules(50.0, []string{"Cluster"})
	rep := validator.Validate(results, rules)
	if !rep.IsClean() {
		t.Errorf("expected clean report, got violations: %v", rep.Violations)
	}
	if len(rep.Passed) != 2 {
		t.Errorf("expected 2 passed rules, got %d", len(rep.Passed))
	}
}

func TestValidate_RemovalThresholdExceeded(t *testing.T) {
	results := []differ.Result{
		{Type: "Cluster", Name: "a", Status: differ.Removed},
		{Type: "Cluster", Name: "b", Status: differ.Removed},
		{Type: "Cluster", Name: "c", Status: differ.Unchanged},
	}
	rules := validator.DefaultRules(50.0, nil)
	rep := validator.Validate(results, rules)
	if rep.IsClean() {
		t.Error("expected violations for removal threshold")
	}
	found := false
	for _, v := range rep.Violations {
		if strings.Contains(v, "max-removal-percentage") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected max-removal-percentage violation, got: %v", rep.Violations)
	}
}

func TestValidate_UnexpectedType(t *testing.T) {
	results := []differ.Result{
		{Type: "Cluster", Name: "a", Status: differ.Added},
		{Type: "RouteConfiguration", Name: "b", Status: differ.Added},
	}
	rules := validator.DefaultRules(100.0, []string{"Cluster"})
	rep := validator.Validate(results, rules)
	if rep.IsClean() {
		t.Error("expected violation for unexpected type")
	}
	if !strings.Contains(rep.Violations[0], "RouteConfiguration") {
		t.Errorf("expected RouteConfiguration in violation, got: %s", rep.Violations[0])
	}
}

func TestValidate_EmptyResults(t *testing.T) {
	rules := validator.DefaultRules(50.0, []string{"Cluster"})
	rep := validator.Validate(nil, rules)
	if !rep.IsClean() {
		t.Errorf("expected clean report for empty results, got: %v", rep.Violations)
	}
}

func TestReport_Write(t *testing.T) {
	rep := &validator.Report{
		Passed:     []string{"rule-a"},
		Violations: []string{"rule-b: something wrong"},
	}
	out := rep.Write()
	if !strings.Contains(out, "[FAIL]") {
		t.Error("expected [FAIL] in output")
	}
	if !strings.Contains(out, "[PASS]") {
		t.Error("expected [PASS] in output")
	}
	if !strings.Contains(out, "1 violation") {
		t.Error("expected violation count in output")
	}
}
