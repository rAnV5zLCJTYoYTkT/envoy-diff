package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/envoy-diff/internal/snapshot"
)

// ResourceMap is a raw map of resource type to list of JSON objects.
type ResourceMap map[string][]json.RawMessage

// FromFile reads a JSON snapshot file from disk and returns a populated Snapshot.
// The file is expected to be a JSON object where keys are xDS resource types
// and values are arrays of resource JSON objects (each must have a "name" field).
func FromFile(path string) (*snapshot.Snapshot, error) {
	path = filepath.Clean(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loader: read file %q: %w", path, err)
	}
	return FromBytes(data)
}

// FromBytes parses raw JSON bytes into a Snapshot.
func FromBytes(data []byte) (*snapshot.Snapshot, error) {
	var rm ResourceMap
	if err := json.Unmarshal(data, &rm); err != nil {
		return nil, fmt.Errorf("loader: unmarshal snapshot: %w", err)
	}

	snap := snapshot.New()
	for rtype, resources := range rm {
		for _, raw := range resources {
			var meta struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(raw, &meta); err != nil {
				return nil, fmt.Errorf("loader: resource in type %q missing name: %w", rtype, err)
			}
			if meta.Name == "" {
				return nil, fmt.Errorf("loader: resource in type %q has empty name", rtype)
			}
			snap.AddResource(rtype, meta.Name, raw)
		}
	}
	return snap, nil
}
