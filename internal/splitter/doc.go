// Package splitter partitions diff results into named buckets using a
// configurable key function.
//
// Built-in key functions are provided for common dimensions:
//
//   - ByType    — splits by xDS resource type (e.g. "Listener", "Cluster")
//   - ByStatus  — splits by diff status (added / removed / modified / unchanged)
//   - ByLabel   — splits by a named annotation label attached to each result
//
// Custom key functions can be supplied directly to Apply for domain-specific
// partitioning strategies.
package splitter
