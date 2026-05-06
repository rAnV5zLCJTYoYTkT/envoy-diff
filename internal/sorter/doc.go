// Package sorter orders envoy-diff DiffResult slices by configurable
// criteria — name, resource type, or diff status — without mutating
// the original slice.
//
// Usage:
//
//	sorted := sorter.ByField(results, sorter.ByStatus)
//
//	// or with functional options:
//	sorted := sorter.ApplyWithOptions(results,
//		sorter.WithField(sorter.ByName),
//		sorter.WithDescending(),
//	)
package sorter
