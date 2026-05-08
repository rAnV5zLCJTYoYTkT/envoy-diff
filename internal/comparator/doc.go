// Package comparator extends the core differ to support multi-environment
// snapshot comparisons. It accepts a named map of Snapshots and produces a
// Report containing one Pair entry per compared environment tuple.
//
// By default environments are compared sequentially (dev→staging→prod).
// Set Options.Sequential = false to compare every unique pair.
package comparator
