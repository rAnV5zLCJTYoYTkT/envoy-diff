package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls output style.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Render writes the diff result to w in the requested format.
func Render(w io.Writer, result *Result, format Format) error {
	switch format {
	case FormatJSON:
		return renderJSON(w, result)
	case FormatText:
		return renderText(w, result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func renderJSON(w io.Writer, result *Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result.Diffs)
}

func renderText(w io.Writer, result *Result) error {
	if len(result.Diffs) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	for _, d := range result.Diffs {
		if d.Status == StatusUnchanged {
			continue
		}
		prefix := statusPrefix(d.Status)
		_, err := fmt.Fprintf(w, "%s [%s] %s\n", prefix, d.Type, d.Name)
		if err != nil {
			return err
		}
		if d.Status == StatusModified {
			_, err = fmt.Fprintf(w, "  - %s\n  + %s\n",
				truncate(d.Left, 120), truncate(d.Right, 120))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func statusPrefix(s DiffStatus) string {
	switch s {
	case StatusAdded:
		return "+"
	case StatusRemoved:
		return "-"
	case StatusModified:
		return "~"
	default:
		return " "
	}
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
