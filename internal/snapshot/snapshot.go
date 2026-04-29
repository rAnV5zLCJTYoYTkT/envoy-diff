package snapshot

import "time"

// ResourceType represents the type of xDS resource.
type ResourceType string

const (
	ResourceTypeClusters  ResourceType = "clusters"
	ResourceTypeListeners ResourceType = "listeners"
	ResourceTypeRoutes    ResourceType = "routes"
	ResourceTypeEndpoints ResourceType = "endpoints"
)

// Snapshot represents an Envoy xDS config snapshot from a single environment.
type Snapshot struct {
	// Environment identifies the source environment (e.g., "staging", "production").
	Environment string `json:"environment"`

	// CapturedAt is the time the snapshot was taken.
	CapturedAt time.Time `json:"captured_at"`

	// Resources holds the raw xDS resource data keyed by resource type.
	Resources map[ResourceType][]Resource `json:"resources"`
}

// Resource represents a single xDS resource entry.
type Resource struct {
	// Name is the unique identifier of the resource.
	Name string `json:"name"`

	// Version is the xDS version string for this resource.
	Version string `json:"version"`

	// Body holds the raw JSON-encoded resource payload.
	Body map[string]interface{} `json:"body"`
}

// New creates a new empty Snapshot for the given environment.
func New(env string) *Snapshot {
	return &Snapshot{
		Environment: env,
		CapturedAt:  time.Now().UTC(),
		Resources:   make(map[ResourceType][]Resource),
	}
}

// AddResource appends a resource to the snapshot under the given type.
func (s *Snapshot) AddResource(rt ResourceType, r Resource) {
	s.Resources[rt] = append(s.Resources[rt], r)
}

// ResourceNames returns the set of resource names for a given type.
func (s *Snapshot) ResourceNames(rt ResourceType) map[string]struct{} {
	names := make(map[string]struct{})
	for _, r := range s.Resources[rt] {
		names[r.Name] = struct{}{}
	}
	return names
}
