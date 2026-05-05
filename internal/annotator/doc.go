// Package annotator enriches diff results with human-readable labels,
// summaries, and tags by applying a configurable set of annotation rules.
//
// Rules are plain functions of type Rule and can be composed freely.
// DefaultRules provides a standard set suitable for most auditing workflows.
package annotator
