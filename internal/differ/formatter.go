package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Render writes diff results to w in the given format ("text" or "json").
func Render(w io.Writer, results []Result, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return renderJSON(w, results)
	default:
		return renderText(w, results)
	}
}

func renderJSON(w io.Writer, results []Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func renderText(w io.Writer, results []Result) error {
	for _, r := range results {
		prefix := statusPrefix(r.Status)
		line := fmt.Sprintf("%s [%s] %s/%s", prefix, r.Type, r.Name, truncate(r.NewValue, 60))
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No diffs found.")
		return err
	}
	return nil
}

func statusPrefix(s Status) string {
	switch s {
	case Added:
		return "+"
	case Removed:
		return "-"
	case Modified:
		return "~"
	default:
		return " "
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
