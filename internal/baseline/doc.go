// Package baseline provides save/load/compare operations for envoy-diff
// result snapshots. A baseline captures a known-good (or reference) state
// of diff results so that subsequent runs can highlight drift — resources
// that have been added, removed, or modified since the baseline was taken.
//
// Typical workflow:
//
//	// Save today's diff as the baseline.
//	_ = baseline.Save("baseline.json", "v1.2.0", results)
//
//	// On a later run, detect what changed since the baseline.
//	base, _ := baseline.Load("baseline.json")
//	drift   := baseline.Compare(base, newResults)
package baseline
