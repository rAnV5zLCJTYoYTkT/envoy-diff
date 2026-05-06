// Package printer provides formatted table output for diff results.
package printer

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/example/envoy-diff/internal/differ"
)

const (
	defaultMinWidth = 0
	defaultTabWidth = 8
	defaultPadding  = 2
)

// Options controls table rendering behaviour.
type Options struct {
	ShowUnchanged bool
	MaxNameWidth  int
}

// DefaultOptions returns sensible defaults for table printing.
func DefaultOptions() Options {
	return Options{
		ShowUnchanged: false,
		MaxNameWidth:  60,
	}
}

// WriteTable writes a human-readable aligned table of diff results to w.
func WriteTable(w io.Writer, results []differ.Result, opts Options) error {
	tw := tabwriter.NewWriter(w, defaultMinWidth, defaultTabWidth, defaultPadding, ' ', 0)

	fmt.Fprintln(tw, "STATUS\tTYPE\tNAME\tNOTES")
	fmt.Fprintln(tw, strings.Repeat("-", 6)+"\t"+strings.Repeat("-", 12)+"\t"+strings.Repeat("-", 40)+"\t"+strings.Repeat("-", 20))

	printed := 0
	for _, r := range results {
		if !opts.ShowUnchanged && r.Status == "unchanged" {
			continue
		}
		name := r.Name
		if opts.MaxNameWidth > 0 && len(name) > opts.MaxNameWidth {
			name = name[:opts.MaxNameWidth-3] + "..."
		}
		notes := buildNotes(r)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Status, r.Type, name, notes)
		printed++
	}

	if printed == 0 {
		fmt.Fprintln(tw, "(no differences found)")
	}

	return tw.Flush()
}

// buildNotes assembles a short annotation string from result metadata.
func buildNotes(r differ.Result) string {
	var parts []string
	if v, ok := r.Tags["label"]; ok && v != "" {
		parts = append(parts, v)
	}
	if v, ok := r.Tags["breaking"]; ok && v == "true" {
		parts = append(parts, "BREAKING")
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, ", ")
}
