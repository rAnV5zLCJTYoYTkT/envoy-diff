// Package summarizer aggregates differ.DiffResult slices into
// structured Summary values, providing total and per-type counts
// for added, removed, modified, and unchanged resources.
//
// Usage:
//
//	results := differ.Compare(snapA, snapB)
//	summary := summarizer.Compute(results)
//	summarizer.Write(os.Stdout, summary)
package summarizer
