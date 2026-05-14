// Package differ provides primitives for comparing Envoy xDS config snapshots.
//
// The central type is [Result], which captures the diff outcome for a single
// named resource within a given xDS type. A slice of Results is the primary
// data structure passed between pipeline stages such as the filter, annotator,
// tagger, labeler, scorer, and exporter.
//
// # Status values
//
// Each Result carries a [DiffStatus] that describes the relationship between
// the baseline and target snapshots:
//
//   - [StatusAdded]     – resource is present only in the target snapshot.
//   - [StatusRemoved]   – resource is present only in the baseline snapshot.
//   - [StatusModified]  – resource is present in both but the JSON differs.
//   - [StatusUnchanged] – resource is identical in both snapshots.
//
// # Metadata
//
// Results may carry additional metadata populated by downstream stages:
//
//   - Labels      – arbitrary key/value pairs (see internal/labeler).
//   - Tags        – short classifier strings (see internal/tagger).
//   - Annotations – human-readable notes (see internal/annotator).
package differ
