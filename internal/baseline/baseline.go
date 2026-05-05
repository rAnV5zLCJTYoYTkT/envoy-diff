// Package baseline provides functionality to save and load diff results
// as a baseline for future comparisons, enabling drift detection over time.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/example/envoy-diff/internal/differ"
)

// Record wraps a set of diff results with metadata.
type Record struct {
	CreatedAt time.Time          `json:"created_at"`
	Label     string             `json:"label,omitempty"`
	Results   []differ.DiffResult `json:"results"`
}

// Save writes the given results to a JSON file at path.
func Save(path, label string, results []differ.DiffResult) error {
	rec := Record{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %q: %w", path, err)
	}
	return nil
}

// Load reads a baseline record from the JSON file at path.
func Load(path string) (*Record, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %q: %w", path, err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &rec, nil
}

// Compare returns the results that differ between current and baseline.
// A result is considered new if it does not appear (by type+name+status) in base.
func Compare(base *Record, current []differ.DiffResult) []differ.DiffResult {
	type key struct{ rtype, name, status string }
	seen := make(map[key]struct{}, len(base.Results))
	for _, r := range base.Results {
		seen[key{r.Type, r.Name, string(r.Status)}] = struct{}{}
	}
	var drift []differ.DiffResult
	for _, r := range current {
		if _, ok := seen[key{r.Type, r.Name, string(r.Status)}]; !ok {
			drift = append(drift, r)
		}
	}
	return drift
}
