// Package exporter handles writing envoy-diff results to external
// destinations. It supports multiple output formats (JSON, CSV, plain text)
// and can target either stdout or a file path.
//
// Usage:
//
//	err := exporter.Export(results, exporter.Options{
//		Format:   exporter.FormatCSV,
//		FilePath: "./diff-report.csv",
//	})
//
// Supported formats:
//   - FormatJSON  — structured JSON array of Result objects
//   - FormatCSV   — comma-separated with header row (type, name, status)
//   - FormatText  — human-readable one-line-per-result output
package exporter
