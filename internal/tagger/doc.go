// Package tagger provides rule-based tagging for envoy-diff results.
//
// Tags are key-value pairs attached to each [TaggedResult] after evaluating
// a prioritised list of [Rule] predicates. Multiple rules may match a single
// result; all matching tags are accumulated.
//
// Usage:
//
//	t := tagger.New(tagger.DefaultRules())
//	tagged := t.Apply(results)
//	for _, tr := range tagged {
//		fmt.Println(tr.Name, tr.Tags)
//	}
package tagger
