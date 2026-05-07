// Package deduplicator provides utilities for removing duplicate
// differ.Result entries from a slice. Duplicates are identified by
// a configurable key function that defaults to "<type>/<name>".
//
// Usage:
//
//	results = deduplicator.Apply(results, deduplicator.DefaultOptions())
package deduplicator
