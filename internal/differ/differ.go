// Package differ compares two xDS config snapshots and produces structured diff results.
package differ

import (
	"encoding/json"
	"sort"

	"github.com/example/envoy-diff/internal/snapshot"
)

// Status represents the diff status of a resource.
type Status string

const (
	StatusAdded     Status = "added"
	StatusRemoved   Status = "removed"
	StatusModified  Status = "modified"
	StatusUnchanged Status = "unchanged"
)

// Result holds the diff result for a single xDS resource.
type Result struct {
	// Type is the xDS resource type (e.g. "Cluster", "Listener").
	Type string
	// Name is the resource name.
	Name string
	// Status is the diff status.
	Status Status
	// BaselineJSON is the JSON representation from the baseline snapshot (may be empty).
	BaselineJSON string
	// CandidateJSON is the JSON representation from the candidate snapshot (may be empty).
	CandidateJSON string
	// Tags holds arbitrary metadata attached by downstream stages.
	Tags map[string]string
	// Annotations holds human-readable notes attached by downstream stages.
	Annotations []string
	// Score is an optional numeric weight assigned by the scorer.
	Score float64
}

// Compare computes the diff between a baseline and candidate snapshot.
// It returns one Result per resource, across all resource types present in either snapshot.
func Compare(baseline, candidate *snapshot.Snapshot) []Result {
	var results []Result

	types := unionTypes(baseline, candidate)
	sort.Strings(types)

	for _, rtype := range types {
		baseNames := baseline.ResourceNames(rtype)
		candNames := candidate.ResourceNames(rtype)

		allNames := unionNames(baseNames, candNames)
		sort.Strings(allNames)

		for _, name := range allNames {
			baseVal, inBase := baseline.Get(rtype, name)
			candVal, inCand := candidate.Get(rtype, name)

			var status Status
			switch {
			case inBase && !inCand:
				status = StatusRemoved
			case !inBase && inCand:
				status = StatusAdded
			default:
				if jsonEqual(baseVal, candVal) {
					status = StatusUnchanged
				} else {
					status = StatusModified
				}
			}

			results = append(results, Result{
				Type:          rtype,
				Name:          name,
				Status:        status,
				BaselineJSON:  marshalOrEmpty(baseVal),
				CandidateJSON: marshalOrEmpty(candVal),
				Tags:          make(map[string]string),
			})
		}
	}

	return results
}

// unionTypes returns all resource type keys present in either snapshot.
func unionTypes(a, b *snapshot.Snapshot) []string {
	return keys(union(a.Types(), b.Types()))
}

// unionNames returns all names present in either slice.
func unionNames(a, b []string) []string {
	return keys(union(a, b))
}

func union(a, b []string) map[string]struct{} {
	m := make(map[string]struct{}, len(a)+len(b))
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		m[v] = struct{}{}
	}
	return m
}

func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// jsonEqual reports whether two values produce identical JSON representations.
func jsonEqual(a, b interface{}) bool {
	aj, err1 := json.Marshal(a)
	bj, err2 := json.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(aj) == string(bj)
}

// marshalOrEmpty marshals v to a JSON string, returning empty string on error or nil.
func marshalOrEmpty(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}
