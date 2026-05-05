// Package annotator attaches human-readable labels and metadata to diff results
// to provide additional context for auditing and reporting.
package annotator

import (
	"fmt"
	"strings"

	"github.com/your-org/envoy-diff/internal/differ"
)

// Annotation holds extra context attached to a single diff result.
type Annotation struct {
	Label   string
	Summary string
	Tags    []string
}

// AnnotatedResult pairs a diff result with its annotation.
type AnnotatedResult struct {
	Result     differ.DiffResult
	Annotation Annotation
}

// Rule is a function that produces an Annotation for a given DiffResult.
type Rule func(r differ.DiffResult) (Annotation, bool)

// Annotator applies a set of rules to produce annotated results.
type Annotator struct {
	rules []Rule
}

// New creates an Annotator with the provided rules.
func New(rules ...Rule) *Annotator {
	return &Annotator{rules: rules}
}

// Annotate applies all rules to each result, merging annotations.
func (a *Annotator) Annotate(results []differ.DiffResult) []AnnotatedResult {
	out := make([]AnnotatedResult, 0, len(results))
	for _, r := range results {
		ann := Annotation{}
		for _, rule := range a.rules {
			if a2, ok := rule(r); ok {
				ann = merge(ann, a2)
			}
		}
		if ann.Label == "" {
			ann.Label = string(r.Status)
		}
		out = append(out, AnnotatedResult{Result: r, Annotation: ann})
	}
	return out
}

func merge(base, extra Annotation) Annotation {
	if extra.Label != "" {
		base.Label = extra.Label
	}
	if extra.Summary != "" {
		base.Summary = extra.Summary
	}
	base.Tags = append(base.Tags, extra.Tags...)
	return base
}

// DefaultRules returns a standard set of annotation rules.
func DefaultRules() []Rule {
	return []Rule{
		BreakingChangeRule(),
		ResourceTypeRule(),
	}
}

// BreakingChangeRule flags removed resources as breaking changes.
func BreakingChangeRule() Rule {
	return func(r differ.DiffResult) (Annotation, bool) {
		if r.Status == differ.StatusRemoved {
			return Annotation{
				Label:   "breaking",
				Summary: fmt.Sprintf("Resource %q was removed — may cause traffic disruption", r.Name),
				Tags:    []string{"breaking", "removal"},
			}, true
		}
		return Annotation{}, false
	}
}

// ResourceTypeRule adds a tag based on the xDS resource type short name.
func ResourceTypeRule() Rule {
	return func(r differ.DiffResult) (Annotation, bool) {
		parts := strings.Split(r.Type, ".")
		short := strings.ToLower(parts[len(parts)-1])
		return Annotation{
			Tags: []string{"type:" + short},
		}, true
	}
}
