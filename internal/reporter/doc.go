// Package reporter aggregates and summarises the results produced by the
// differ package into human-readable and machine-consumable reports.
//
// Usage:
//
//	results := differ.Compare(snapA, snapB)
//	reporter.Write(os.Stdout, results)
//
// The package exposes three entry points:
//
//   - Compute – returns an overall Summary for a slice of DiffResults.
//   - ByType  – returns per-resource-type summaries.
//   - Write   – writes a formatted summary to any io.Writer.
package reporter
