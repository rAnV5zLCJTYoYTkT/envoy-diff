// Package ranker provides priority ranking for Envoy xDS diff results.
//
// Results are assigned a numeric rank based on their change status
// (removed > modified > added > unchanged by default). Callers may
// supply custom weight maps via Options to adjust severity ordering.
//
// Typical usage:
//
//	ranked := ranker.Rank(results, ranker.DefaultOptions())
//	plain := ranker.Unwrap(ranked)
package ranker
