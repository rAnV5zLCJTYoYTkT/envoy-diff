// Package classifier assigns severity levels to diff results based on
// configurable rules, enabling downstream consumers to prioritize changes.
package classifier

import (
	"strings"

	"github.com/example/envoy-diff/internal/differ"
)

// Severity represents the impact level of a diff result.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Rule maps a condition to a severity level.
type Rule struct {
	Status      string
	TypeContains string
	NameContains string
	Severity    Severity
}

// Result wraps a differ.Result with an assigned severity.
type Result struct {
	difer.Result
	Severity Severity
}

// DefaultRules provides a baseline set of classification rules.
var DefaultRules = []Rule{
	{Status: "removed", TypeContains: "listener", Severity: SeverityCritical},
	{Status: "removed", TypeContains: "cluster", Severity: SeverityHigh},
	{Status: "removed", Severity: SeverityHigh},
	{Status: "added", Severity: SeverityMedium},
	{Status: "modified", TypeContains: "listener", Severity: SeverityHigh},
	{Status: "modified", Severity: SeverityMedium},
	{Status: "unchanged", Severity: SeverityInfo},
}

// Classify applies rules to each result and returns annotated Results.
// Rules are evaluated in order; the first match wins.
func Classify(results []differ.Result, rules []Rule) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Result:   r,
			Severity: classify(r, rules),
		})
	}
	return out
}

func classify(r differ.Result, rules []Rule) Severity {
	for _, rule := range rules {
		if rule.Status != "" && !strings.EqualFold(rule.Status, r.Status) {
			continue
		}
		if rule.TypeContains != "" && !strings.Contains(strings.ToLower(r.Type), strings.ToLower(rule.TypeContains)) {
			continue
		}
		if rule.NameContains != "" && !strings.Contains(strings.ToLower(r.Name), strings.ToLower(rule.NameContains)) {
			continue
		}
		return rule.Severity
	}
	return SeverityInfo
}
