package differ

import (
	"testing"
)

func TestResult_IsChanged(t *testing.T) {
	tests := []struct {
		name   string
		status DiffStatus
		want   bool
	}{
		{"added is changed", StatusAdded, true},
		{"removed is changed", StatusRemoved, true},
		{"modified is changed", StatusModified, true},
		{"unchanged is not changed", StatusUnchanged, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := Result{Status: tc.status}
			if got := r.IsChanged(); got != tc.want {
				t.Errorf("IsChanged() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResult_Clone_IndependentLabels(t *testing.T) {
	orig := Result{
		Type:   "listener",
		Name:   "ingress",
		Status: StatusModified,
		Labels: map[string]string{"env": "prod"},
		Tags:   []string{"breaking"},
		Annotations: []string{"note"},
	}

	cloned := orig.Clone()

	// Mutate clone — original must be unaffected.
	cloned.Labels["env"] = "staging"
	cloned.Tags[0] = "safe"
	cloned.Annotations[0] = "other"

	if orig.Labels["env"] != "prod" {
		t.Errorf("Clone mutated original Labels: got %q", orig.Labels["env"])
	}
	if orig.Tags[0] != "breaking" {
		t.Errorf("Clone mutated original Tags: got %q", orig.Tags[0])
	}
	if orig.Annotations[0] != "note" {
		t.Errorf("Clone mutated original Annotations: got %q", orig.Annotations[0])
	}
}

func TestResult_Clone_NilMaps(t *testing.T) {
	orig := Result{
		Type:   "cluster",
		Name:   "backend",
		Status: StatusAdded,
	}

	cloned := orig.Clone()
	if cloned.Labels != nil {
		t.Errorf("expected nil Labels in clone, got %v", cloned.Labels)
	}
	if cloned.Tags != nil {
		t.Errorf("expected nil Tags in clone, got %v", cloned.Tags)
	}
	if cloned.Annotations != nil {
		t.Errorf("expected nil Annotations in clone, got %v", cloned.Annotations)
	}
}

func TestDiffStatus_Constants(t *testing.T) {
	statuses := []DiffStatus{StatusAdded, StatusRemoved, StatusModified, StatusUnchanged}
	seen := make(map[DiffStatus]bool)
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("duplicate DiffStatus constant: %q", s)
		}
		seen[s] = true
		if s == "" {
			t.Error("DiffStatus constant must not be empty string")
		}
	}
}
