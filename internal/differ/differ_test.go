package differ_test

import (
	"testing"

	"github.com/envoy-diff/internal/differ"
	"github.com/envoy-diff/internal/snapshot"
)

func makeSnapshot(resources map[string]map[string]string) *snapshot.Snapshot {
	s := snapshot.New()
	for rType, entries := range resources {
		for name, body := range entries {
			s.AddResource(rType, name, body)
		}
	}
	return s
}

func TestCompare_Added(t *testing.T) {
	left := makeSnapshot(map[string]map[string]string{})
	right := makeSnapshot(map[string]map[string]string{
		"cluster": {"cluster-a": `{"name":"cluster-a"}`},
	})

	result, err := differ.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result.Diffs))
	}
	if result.Diffs[0].Status != differ.StatusAdded {
		t.Errorf("expected added, got %s", result.Diffs[0].Status)
	}
}

func TestCompare_Removed(t *testing.T) {
	left := makeSnapshot(map[string]map[string]string{
		"listener": {"listener-a": `{"name":"listener-a"}`},
	})
	right := makeSnapshot(map[string]map[string]string{})

	result, err := differ.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Diffs[0].Status != differ.StatusRemoved {
		t.Errorf("expected removed, got %s", result.Diffs[0].Status)
	}
}

func TestCompare_Modified(t *testing.T) {
	left := makeSnapshot(map[string]map[string]string{
		"cluster": {"cluster-a": `{"lb":"round_robin"}`},
	})
	right := makeSnapshot(map[string]map[string]string{
		"cluster": {"cluster-a": `{"lb":"least_request"}`},
	})

	result, err := differ.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Diffs[0].Status != differ.StatusModified {
		t.Errorf("expected modified, got %s", result.Diffs[0].Status)
	}
	if !result.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_Unchanged(t *testing.T) {
	body := `{"name":"cluster-a"}`
	left := makeSnapshot(map[string]map[string]string{
		"cluster": {"cluster-a": body},
	})
	right := makeSnapshot(map[string]map[string]string{
		"cluster": {"cluster-a": body},
	})

	result, err := differ.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompare_NilSnapshot(t *testing.T) {
	_, err := differ.Compare(nil, snapshot.New())
	if err == nil {
		t.Error("expected error for nil left snapshot")
	}
}
