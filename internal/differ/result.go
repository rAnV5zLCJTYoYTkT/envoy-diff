package differ

// DiffStatus represents the change status of a resource between two snapshots.
type DiffStatus string

const (
	// StatusAdded indicates the resource exists only in the new snapshot.
	StatusAdded DiffStatus = "added"
	// StatusRemoved indicates the resource exists only in the old snapshot.
	StatusRemoved DiffStatus = "removed"
	// StatusModified indicates the resource exists in both snapshots but differs.
	StatusModified DiffStatus = "modified"
	// StatusUnchanged indicates the resource is identical in both snapshots.
	StatusUnchanged DiffStatus = "unchanged"
)

// Result represents a single diff entry for one xDS resource.
type Result struct {
	// Type is the xDS resource type URL (e.g. "type.googleapis.com/envoy.config.listener.v3.Listener").
	Type string `json:"type"`
	// Name is the resource name within the snapshot.
	Name string `json:"name"`
	// Status describes how the resource changed.
	Status DiffStatus `json:"status"`
	// OldValue holds the JSON-encoded resource from the baseline snapshot, or empty if added.
	OldValue string `json:"old_value,omitempty"`
	// NewValue holds the JSON-encoded resource from the target snapshot, or empty if removed.
	NewValue string `json:"new_value,omitempty"`
	// Labels holds arbitrary key/value metadata attached by downstream stages.
	Labels map[string]string `json:"labels,omitempty"`
	// Tags holds a set of short classifier tags attached by downstream stages.
	Tags []string `json:"tags,omitempty"`
	// Annotations holds human-readable notes produced by the annotator stage.
	Annotations []string `json:"annotations,omitempty"`
}

// IsChanged reports whether the result represents any kind of change
// (i.e. not StatusUnchanged).
func (r Result) IsChanged() bool {
	return r.Status != StatusUnchanged
}

// Clone returns a shallow copy of the result with an independent Labels map.
func (r Result) Clone() Result {
	copy := r
	if r.Labels != nil {
		copy.Labels = make(map[string]string, len(r.Labels))
		for k, v := range r.Labels {
			copy.Labels[k] = v
		}
	}
	if r.Tags != nil {
		copy.Tags = append([]string(nil), r.Tags...)
	}
	if r.Annotations != nil {
		copy.Annotations = append([]string(nil), r.Annotations...)
	}
	return copy
}
