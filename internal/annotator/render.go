package annotator

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteText writes annotated results as human-readable text to w.
func WriteText(w io.Writer, results []AnnotatedResult) error {
	for _, ar := range results {
		tags := ""
		if len(ar.Annotation.Tags) > 0 {
			tags = " [" + strings.Join(ar.Annotation.Tags, ", ") + "]"
		}
		summary := ar.Annotation.Summary
		if summary == "" {
			summary = "-"
		}
		_, err := fmt.Fprintf(w, "%-12s %-40s %-20s%s\n  %s\n",
			ar.Annotation.Label,
			ar.Result.Name,
			ar.Result.Type,
			tags,
			summary,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// jsonAnnotatedResult is the JSON-serialisable form of an AnnotatedResult.
type jsonAnnotatedResult struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Label   string   `json:"label"`
	Summary string   `json:"summary,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// WriteJSON writes annotated results as a JSON array to w.
func WriteJSON(w io.Writer, results []AnnotatedResult) error {
	records := make([]jsonAnnotatedResult, len(results))
	for i, ar := range results {
		records[i] = jsonAnnotatedResult{
			Type:    ar.Result.Type,
			Name:    ar.Result.Name,
			Status:  string(ar.Result.Status),
			Label:   ar.Annotation.Label,
			Summary: ar.Annotation.Summary,
			Tags:    ar.Annotation.Tags,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
