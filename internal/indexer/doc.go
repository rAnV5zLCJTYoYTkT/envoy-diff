// Package indexer provides fast, pre-built lookup indexes over slices of
// differ.Result values. After a diff is computed, callers can wrap the
// results in an Index to perform O(1) lookups by type, status, or exact
// name, and substring searches without repeated linear scans.
//
// Example usage:
//
//	idx := indexer.Build(results)
//	listeners := idx.ByType("listener")
//	added := idx.ByStatus("added")
//	r, ok := idx.ByName("my-cluster")
//	hits := idx.Search("prod")
package indexer
