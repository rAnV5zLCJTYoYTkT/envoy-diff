package differ

import (
	"fmt"
	"sort"

	"github.com/envoy-diff/internal/snapshot"
)

// ResourceDiff represents a diff entry for a single xDS resource.
type ResourceDiff struct {
	Type   string
	Name   string
	Status DiffStatus
	Left   string
	Right  string
}

// DiffStatus indicates whether a resource was added, removed, or modified.
type DiffStatus string

const (
	StatusAdded    DiffStatus = "added"
	StatusRemoved  DiffStatus = "removed"
	StatusModified DiffStatus = "modified"
	StatusUnchanged DiffStatus = "unchanged"
)

// Result holds the full diff between two snapshots.
type Result struct {
	Diffs []ResourceDiff
}

// HasChanges returns true if there are any non-unchanged diffs.
func (r *Result) HasChanges() bool {
	for _, d := range r.Diffs {
		if d.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

// Compare computes the diff between two snapshots.
func Compare(left, right *snapshot.Snapshot) (*Result, error) {
	if left == nil || right == nil {
		return nil, fmt.Errorf("both snapshots must be non-nil")
	}

	result := &Result{}

	allTypes := unionTypes(left.ResourceTypes(), right.ResourceTypes())

	for _, rType := range allTypes {
		leftResources := left.Resources(rType)
		rightResources := right.Resources(rType)

		allNames := unionNames(keys(leftResources), keys(rightResources))

		for _, name := range allNames {
			lVal, lOk := leftResources[name]
			rVal, rOk := rightResources[name]

			var status DiffStatus
			switch {
			case lOk && !rOk:
				status = StatusRemoved
			case !lOk && rOk:
				status = StatusAdded
			case lVal == rVal:
				status = StatusUnchanged
			default:
				status = StatusModified
			}

			result.Diffs = append(result.Diffs, ResourceDiff{
				Type:   rType,
				Name:   name,
				Status: status,
				Left:   lVal,
				Right:  rVal,
			})
		}
	}

	return result, nil
}

func unionTypes(a, b []string) []string {
	return unionNames(a, b)
}

func unionNames(a, b []string) []string {
	seen := make(map[string]struct{})
	for _, v := range a {
		seen[v] = struct{}{}
	}
	for _, v := range b {
		seen[v] = struct{}{}
	}
	result := make([]string, 0, len(seen))
	for k := range seen {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func keys(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
