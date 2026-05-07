// Package differ compares two Envoy xDS snapshots and produces a list of
// Result values describing what was added, removed, modified, or left
// unchanged between the two configurations.
//
// Usage:
//
//	results := differ.Compare(snapshotA, snapshotB)
//	_ = differ.Render(os.Stdout, results, "text")
package differ
