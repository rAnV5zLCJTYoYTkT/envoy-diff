package differ

import (
	"encoding/json"
	"fmt"

	"github.com/your-org/envoy-diff/internal/snapshot"
)

// Status represents the diff status of a resource.
type Status string

const (
	Added     Status = "added"
	Removed   Status = "removed"
	Modified  Status = "modified"
	Unchanged Status = "unchanged"
)

// Result holds a single diff result for one resource.
type Result struct {
	Type     string
	Name     string
	Status   Status
	OldValue string
	NewValue string
	Tags     map[string]string
	Labels   []string
	Notes    []string
	Score    float64
}

// Compare produces a list of Results by diffing two snapshots.
func Compare(a, b *snapshot.Snapshot) []Result {
	var results []Result
	for _, t := range unionTypes(a, b) {
		oldNames := a.ResourceNames(t)
		newNames := b.ResourceNames(t)
		for _, name := range unionNames(oldNames, newNames) {
			oldVal, hasOld := a.Get(t, name)
			newVal, hasNew := b.Get(t, name)
			switch {
			case hasOld && !hasNew:
				results = append(results, Result{Type: t, Name: name, Status: Removed, OldValue: oldVal})
			case !hasOld && hasNew:
				results = append(results, Result{Type: t, Name: name, Status: Added, NewValue: newVal})
			case oldVal != newVal:
				results = append(results, Result{Type: t, Name: name, Status: Modified, OldValue: oldVal, NewValue: newVal})
			default:
				results = append(results, Result{Type: t, Name: name, Status: Unchanged, OldValue: oldVal, NewValue: newVal})
			}
		}
	}
	return results
}

func unionTypes(a, b *snapshot.Snapshot) []string {
	return union(a.Types(), b.Types())
}

func unionNames(a, b []string) []string {
	return union(a, b)
}

func union(a, b []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, v := range append(a, b...) {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func keys(m map[string]string) []string {
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func marshalIndent(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
