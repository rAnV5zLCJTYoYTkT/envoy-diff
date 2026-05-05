// Package tagger assigns environment-aware tags to diff results based on
// configurable rules, enabling downstream consumers to categorize and route
// diffs by environment, region, or custom metadata.
package tagger

import (
	"strings"

	"github.com/your-org/envoy-diff/internal/differ"
)

// Tag represents a key-value label attached to a diff result.
type Tag struct {
	Key   string
	Value string
}

// Rule defines a matching predicate and the tags to apply when it matches.
type Rule struct {
	Match func(r differ.Result) bool
	Tags  []Tag
}

// Tagger applies a set of rules to diff results and returns tagged copies.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger with the provided rules.
func New(rules []Rule) *Tagger {
	return &Tagger{rules: rules}
}

// TaggedResult wraps a differ.Result with additional tags.
type TaggedResult struct {
	difer.Result
	Tags []Tag
}

// Apply iterates over results and attaches matching tags from all rules.
func (t *Tagger) Apply(results []differ.Result) []TaggedResult {
	out := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		tr := TaggedResult{Result: r}
		for _, rule := range t.rules {
			if rule.Match(r) {
				tr.Tags = append(tr.Tags, rule.Tags...)
			}
		}
		out = append(out, tr)
	}
	return out
}

// DefaultRules returns a baseline set of tagging rules suitable for most
// Envoy xDS auditing workflows.
func DefaultRules() []Rule {
	return []Rule{
		{
			Match: func(r differ.Result) bool { return r.Status == differ.Added },
			Tags:  []Tag{{Key: "change", Value: "added"}},
		},
		{
			Match: func(r differ.Result) bool { return r.Status == differ.Removed },
			Tags:  []Tag{{Key: "change", Value: "removed"}},
		},
		{
			Match: func(r differ.Result) bool { return r.Status == differ.Modified },
			Tags:  []Tag{{Key: "change", Value: "modified"}},
		},
		{
			Match: func(r differ.Result) bool {
				return strings.HasPrefix(r.Type, "type.googleapis.com/envoy.config.listener")
			},
			Tags: []Tag{{Key: "resource-family", Value: "listener"}},
		},
		{
			Match: func(r differ.Result) bool {
				return strings.HasPrefix(r.Type, "type.googleapis.com/envoy.config.cluster")
			},
			Tags: []Tag{{Key: "resource-family", Value: "cluster"}},
		},
	}
}
