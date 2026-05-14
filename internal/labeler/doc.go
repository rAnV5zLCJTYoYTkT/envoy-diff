// Package labeler provides rule-based labeling for envoy-diff results.
//
// Labels are short, human-readable strings stored in the Tags map under the
// key "label". They are used by downstream components (printers, exporters,
// reporters) to group or highlight results.
//
// Usage:
//
//	l := labeler.New()           // use default rules
//	labeled := l.Apply(results)
package labeler
