// Package patcher applies a set of patch rules to diff results,
// allowing field-level overrides or suppression of known differences.
package patcher

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/differ"
)

// Rule describes a single patch operation to apply to matching diff results.
type Rule struct {
	// Type filters by resource type (empty means match all).
	Type string
	// NameContains filters by substring match on resource name.
	NameContains string
	// SuppressStatus, when set, marks matching results as suppressed.
	SuppressStatus bool
	// OverrideLabel replaces the result label with this value.
	OverrideLabel string
}

// Patcher holds a collection of rules and applies them to diff results.
type Patcher struct {
	rules []Rule
}

// New creates a Patcher with the provided rules.
func New(rules []Rule) *Patcher {
	return &Patcher{rules: rules}
}

// Apply iterates over results and applies matching patch rules.
// Results that are suppressed have their Status set to "suppressed".
func (p *Patcher) Apply(results []differ.Result) []differ.Result {
	out := make([]differ.Result, 0, len(results))
	for _, r := range results {
		r = p.patch(r)
		out = append(out, r)
	}
	return out
}

func (p *Patcher) patch(r differ.Result) differ.Result {
	for _, rule := range p.rules {
		if !p.matches(rule, r) {
			continue
		}
		if rule.SuppressStatus {
			r.Status = "suppressed"
		}
		if rule.OverrideLabel != "" {
			if r.Tags == nil {
				r.Tags = map[string]string{}
			}
			r.Tags["label"] = rule.OverrideLabel
		}
	}
	return r
}

func (p *Patcher) matches(rule Rule, r differ.Result) bool {
	if rule.Type != "" && rule.Type != r.Type {
		return false
	}
	if rule.NameContains != "" && !strings.Contains(r.Name, rule.NameContains) {
		return false
	}
	return true
}
