// Package merger combines multiple diff results slices into a single unified
// result set, deduplicating by (type, name) and resolving conflicts by status
// priority: Modified > Added > Removed > Unchanged.
package merger

import "github.com/yourorg/envoy-diff/internal/differ"

// Priority defines the resolution order when the same (type, name) key appears
// in more than one input slice.
var priority = map[string]int{
	"modified":  4,
	"added":     3,
	"removed":   2,
	"unchanged": 1,
}

func rankOf(status string) int {
	if r, ok := priority[status]; ok {
		return r
	}
	return 0
}

// key uniquely identifies a diff result.
type key struct {
	resourceType string
	name         string
}

// Merge combines one or more slices of differ.Result into a single slice.
// When the same (ResourceType, Name) appears in multiple inputs the entry with
// the highest status priority is kept. Ties are broken by keeping the first
// occurrence.
func Merge(inputs ...[]differ.Result) []differ.Result {
	seen := make(map[key]differ.Result)

	for _, results := range inputs {
		for _, r := range results {
			k := key{resourceType: r.ResourceType, name: r.Name}
			existing, found := seen[k]
			if !found || rankOf(string(r.Status)) > rankOf(string(existing.Status)) {
				seen[k] = r
			}
		}
	}

	out := make([]differ.Result, 0, len(seen))
	for _, r := range seen {
		out = append(out, r)
	}
	return out
}

// MergeWithLabel is like Merge but stamps each result with a label tag so the
// origin of each entry is traceable in downstream rendering.
func MergeWithLabel(label string, inputs ...[]differ.Result) []differ.Result {
	results := Merge(inputs...)
	for i := range results {
		if results[i].Tags == nil {
			results[i].Tags = make(map[string]string)
		}
		results[i].Tags["merged_from"] = label
	}
	return results
}
