// Package labeler assigns human-readable labels to diff results based on
// configurable rules, enriching results with contextual metadata for display
// and downstream processing.
package labeler

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/differ"
)

// Rule defines a labeling rule: if the predicate matches a result, the label
// is applied.
type Rule struct {
	Label     string
	Predicate func(r differ.Result) bool
}

// Labeler applies a set of rules to annotate results with labels.
type Labeler struct {
	rules []Rule
}

// New returns a Labeler with the provided rules. If no rules are given,
// DefaultRules are used.
func New(rules ...Rule) *Labeler {
	if len(rules) == 0 {
		rules = DefaultRules()
	}
	return &Labeler{rules: rules}
}

// Apply iterates over results and assigns the first matching rule's label.
// Results that match no rule are left unchanged.
func (l *Labeler) Apply(results []differ.Result) []differ.Result {
	out := make([]differ.Result, len(results))
	for i, r := range results {
		for _, rule := range l.rules {
			if rule.Predicate(r) {
				if r.Tags == nil {
					r.Tags = map[string]string{}
				}
				r.Tags["label"] = rule.Label
				break
			}
		}
		out[i] = r
	}
	return out
}

// DefaultRules returns a set of sensible default labeling rules based on
// resource type and diff status.
func DefaultRules() []Rule {
	return []Rule{
		{
			Label: "listener-change",
			Predicate: func(r differ.Result) bool {
				return strings.EqualFold(r.Type, "listener") && r.Status != differ.StatusUnchanged
			},
		},
		{
			Label: "cluster-change",
			Predicate: func(r differ.Result) bool {
				return strings.EqualFold(r.Type, "cluster") && r.Status != differ.StatusUnchanged
			},
		},
		{
			Label: "route-change",
			Predicate: func(r differ.Result) bool {
				return strings.EqualFold(r.Type, "route") && r.Status != differ.StatusUnchanged
			},
		},
		{
			Label: "added",
			Predicate: func(r differ.Result) bool {
				return r.Status == differ.StatusAdded
			},
		},
		{
			Label: "removed",
			Predicate: func(r differ.Result) bool {
				return r.Status == differ.StatusRemoved
			},
		},
	}
}
