// Package exporter provides functionality for exporting diff results
// to various output formats such as JSON files, CSV, or stdout streams.
package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/envoy-diff/internal/differ"
)

// Format represents the output format for export.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
	FormatText Format = "text"
)

// Options configures the export behavior.
type Options struct {
	Format   Format
	FilePath string // empty means stdout
}

// Export writes diff results to the configured destination.
func Export(results []differ.Result, opts Options) error {
	var w io.Writer
	if opts.FilePath == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(opts.FilePath)
		if err != nil {
			return fmt.Errorf("exporter: create file %q: %w", opts.FilePath, err)
		}
		defer f.Close()
		w = f
	}

	switch opts.Format {
	case FormatJSON:
		return writeJSON(w, results)
	case FormatCSV:
		return writeCSV(w, results)
	case FormatText:
		return writeText(w, results)
	default:
		return fmt.Errorf("exporter: unsupported format %q", opts.Format)
	}
}

func writeJSON(w io.Writer, results []differ.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return fmt.Errorf("exporter: encode json: %w", err)
	}
	return nil
}

func writeCSV(w io.Writer, results []differ.Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"type", "name", "status"}); err != nil {
		return fmt.Errorf("exporter: write csv header: %w", err)
	}
	for _, r := range results {
		if err := cw.Write([]string{r.Type, r.Name, string(r.Status)}); err != nil {
			return fmt.Errorf("exporter: write csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeText(w io.Writer, results []differ.Result) error {
	for _, r := range results {
		_, err := fmt.Fprintf(w, "[%s] %s/%s\n", r.Status, r.Type, r.Name)
		if err != nil {
			return fmt.Errorf("exporter: write text: %w", err)
		}
	}
	return nil
}
