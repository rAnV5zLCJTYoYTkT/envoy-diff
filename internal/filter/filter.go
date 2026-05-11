package filter

import "strings"

// Options holds the filtering criteria for diff results.
type Options struct {
	// Types restricts output to the specified xDS resource types.
	// If empty, all types are included.
	Types []string

	// StatusInclude restricts output to results with the given statuses.
	// Valid values: "added", "removed", "modified", "unchanged".
	// If empty, all statuses are included.
	StatusInclude []string

	// NameContains filters resources whose name contains this substring.
	// Case-insensitive. Empty string disables this filter.
	NameContains string
}

// Result mirrors the differ.Result shape to avoid circular imports.
type Result struct {
	Type   string
	Name   string
	Status string
	Base   string
	Head   string
}

// Apply returns only the results that match all active filter criteria.
func Apply(results []Result, opts Options) []Result {
	typeSet := toSet(opts.Types)
	statusSet := toSet(opts.StatusInclude)

	var out []Result
	for _, r := range results {
		if len(typeSet) > 0 && !typeSet[r.Type] {
			continue
		}
		if len(statusSet) > 0 && !statusSet[r.Status] {
			continue
		}
		if opts.NameContains != "" &&
			!strings.Contains(strings.ToLower(r.Name), strings.ToLower(opts.NameContains)) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// IsEmpty reports whether the options have no active filter criteria,
// meaning Apply would return all results unchanged.
func (o Options) IsEmpty() bool {
	return len(o.Types) == 0 && len(o.StatusInclude) == 0 && o.NameContains == ""
}

func toSet(items []string) map[string]bool {
	if len(items) == 0 {
		return nil
	}
	s := make(map[string]bool, len(items))
	for _, v := range items {
		s[v] = true
	}
	return s
}
