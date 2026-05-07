// Package highlighter provides diff highlighting for xDS resource changes,
// annotating changed fields within a resource's JSON representation.
package highlighter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/envoy-diff/internal/differ"
)

// FieldDiff represents a single field-level change within a resource.
type FieldDiff struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

// Highlight holds the field-level diffs for a single diff result.
type Highlight struct {
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Fields []FieldDiff `json:"fields,omitempty"`
}

// Apply computes field-level diffs for each modified result and returns
// a slice of Highlight values. Unchanged and added/removed results receive
// an empty Fields slice.
func Apply(results []differ.Result) []Highlight {
	highlights := make([]Highlight, 0, len(results))
	for _, r := range results {
		h := Highlight{
			Type:   r.Type,
			Name:   r.Name,
			Status: r.Status,
		}
		if r.Status == "modified" {
			h.Fields = diffJSON(r.Before, r.After)
		}
		highlights = append(highlights, h)
	}
	return highlights
}

// diffJSON compares two JSON strings and returns field-level differences.
func diffJSON(before, after string) []FieldDiff {
	var bMap, aMap map[string]interface{}
	if err := json.Unmarshal([]byte(before), &bMap); err != nil {
		bMap = map[string]interface{}{}
	}
	if err := json.Unmarshal([]byte(after), &aMap); err != nil {
		aMap = map[string]interface{}{}
	}

	keys := unionKeys(bMap, aMap)
	sort.Strings(keys)

	var diffs []FieldDiff
	for _, k := range keys {
		bVal := stringify(bMap[k])
		aVal := stringify(aMap[k])
		if bVal != aVal {
			diffs = append(diffs, FieldDiff{Field: k, Before: bVal, After: aVal})
		}
	}
	return diffs
}

func unionKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", v))
	}
}
