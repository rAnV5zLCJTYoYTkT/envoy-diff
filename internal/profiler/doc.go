// Package profiler provides lightweight stage-level timing instrumentation
// for the envoy-diff pipeline. Each pipeline stage can record its name,
// elapsed duration, and the number of diff results it processed. The
// resulting Profile can be written to any io.Writer for human inspection
// or captured programmatically for metrics export.
package profiler
