package comparator

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/envoy-diff/internal/differ"
)

// WriteText writes a human-readable multi-environment comparison report.
func WriteText(w io.Writer, r Report) {
	if len(r.Pairs) == 0 {
		fmt.Fprintln(w, "no environment pairs to compare")
		return
	}
	for _, p := range r.Pairs {
		fmt.Fprintf(w, "=== %s → %s ===\n", p.Left, p.Right)
		if len(p.Diffs) == 0 {
			fmt.Fprintln(w, "  (no differences)")
			continue
		}
		for _, d := range p.Diffs {
			fmt.Fprintf(w, "  [%s] %s/%s\n", statusLabel(d.Status), d.Type, d.Name)
		}
	}
}

// WriteJSON writes a JSON-encoded multi-environment comparison report.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func statusLabel(s differ.Status) string {
	switch s {
	case differ.Added:
		return "+"
	case differ.Removed:
		return "-"
	case differ.Modified:
		return "~"
	default:
		return "="
	}
}
