// Package aggregator provides utilities for combining multiple sets of
// differ.Result values into a single unified slice and producing roll-up
// statistics (totals, per-type counts, per-status counts).
//
// Typical usage:
//
//	results, summary := aggregator.Aggregate(setA, setB, setC)
//	aggregator.Write(os.Stdout, summary)
package aggregator
