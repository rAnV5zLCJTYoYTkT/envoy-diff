// Package validator provides validation logic for xDS snapshot diff results,
// checking for anomalies such as excessive removals or unexpected resource types.
package validator

import (
	"fmt"
	"strings"

	"github.com/example/envoy-diff/internal/differ"
)

// Rule defines a single validation rule applied to diff results.
type Rule struct {
	Name    string
	Check   func(results []differ.Result) []string
}

// Report holds the outcome of running all validation rules.
type Report struct {
	Passed   []string
	Violations []string
}

// IsClean returns true when no violations were found.
func (r *Report) IsClean() bool {
	return len(r.Violations) == 0
}

// Write formats the report to a string.
func (r *Report) Write() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Validation: %d passed, %d violation(s)\n", len(r.Passed), len(r.Violations)))
	for _, v := range r.Violations {
		sb.WriteString(fmt.Sprintf("  [FAIL] %s\n", v))
	}
	for _, p := range r.Passed {
		sb.WriteString(fmt.Sprintf("  [PASS] %s\n", p))
	}
	return sb.String()
}

// DefaultRules returns the built-in set of validation rules.
func DefaultRules(maxRemovalPct float64, allowedTypes []string) []Rule {
	return []Rule{
		{
			Name: "max-removal-percentage",
			Check: func(results []differ.Result) []string {
				if len(results) == 0 {
					return nil
				}
				removed := 0
				for _, r := range results {
					if r.Status == differ.Removed {
						removed++
					}
				}
				pct := float64(removed) / float64(len(results)) * 100
				if pct > maxRemovalPct {
					return []string{fmt.Sprintf("%.1f%% of resources removed (threshold: %.1f%%)", pct, maxRemovalPct)}
				}
				return nil
			},
		},
		{
			Name: "allowed-resource-types",
			Check: func(results []differ.Result) []string {
				if len(allowedTypes) == 0 {
					return nil
				}
				allowed := make(map[string]struct{}, len(allowedTypes))
				for _, t := range allowedTypes {
					allowed[t] = struct{}{}
				}
				var violations []string
				seen := map[string]bool{}
				for _, r := range results {
					if _, ok := allowed[r.Type]; !ok && !seen[r.Type] {
						violations = append(violations, fmt.Sprintf("unexpected resource type: %s", r.Type))
						seen[r.Type] = true
					}
				}
				return violations
			},
		},
	}
}

// Validate runs all provided rules against the results and returns a Report.
func Validate(results []differ.Result, rules []Rule) *Report {
	rep := &Report{}
	for _, rule := range rules {
		violations := rule.Check(results)
		if len(violations) == 0 {
			rep.Passed = append(rep.Passed, rule.Name)
		} else {
			for _, v := range violations {
				rep.Violations = append(rep.Violations, fmt.Sprintf("%s: %s", rule.Name, v))
			}
		}
	}
	return rep
}
