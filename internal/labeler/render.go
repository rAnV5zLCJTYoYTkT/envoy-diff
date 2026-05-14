package labeler

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/envoy-diff/internal/differ"
)

// WriteText writes a plain-text summary of labeled results to w.
// Only results that carry a "label" tag are included.
func WriteText(w io.Writer, results []differ.Result) {
	fmt.Fprintln(w, "Labeled Results")
	fmt.Fprintln(w, "---------------")
	count := 0
	for _, r := range results {
		label := ""
		if r.Tags != nil {
			label = r.Tags["label"]
		}
		if label == "" {
			continue
		}
		fmt.Fprintf(w, "[%s] %s (%s) — %s\n", label, r.Name, r.Type, r.Status)
		count++
	}
	if count == 0 {
		fmt.Fprintln(w, "(no labeled results)")
	}
}

// labeledEntry is the JSON representation of a single labeled result.
type labeledEntry struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Label  string `json:"label"`
}

// WriteJSON encodes labeled results as a JSON array to w.
// Only results carrying a "label" tag are included.
func WriteJSON(w io.Writer, results []differ.Result) error {
	var entries []labeledEntry
	for _, r := range results {
		label := ""
		if r.Tags != nil {
			label = r.Tags["label"]
		}
		if label == "" {
			continue
		}
		entries = append(entries, labeledEntry{
			Name:   r.Name,
			Type:   r.Type,
			Status: string(r.Status),
			Label:  label,
		})
	}
	if entries == nil {
		entries = []labeledEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
