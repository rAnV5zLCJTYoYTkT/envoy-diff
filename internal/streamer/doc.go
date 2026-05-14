// Package streamer provides incremental, streaming output of differ.Result
// values to an io.Writer.
//
// Unlike batch renderers (e.g. exporter), streamer emits each result as soon
// as it is available, making it suitable for long-running watch pipelines or
// large snapshots where buffering all results in memory is undesirable.
//
// Supported output formats are text (human-readable) and JSON (newline-
// delimited or array form).
package streamer
