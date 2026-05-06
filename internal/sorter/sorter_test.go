package sorter_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/differ"
	"github.com/your-org/envoy-diff/internal/sorter"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "listener", Name: "zebra", Status: differ.StatusAdded},
		{Type: "cluster", Name: "apple", Status: differ.StatusRemoved},
		{Type: "route", Name: "mango", Status: differ.StatusModified},
		{Type: "cluster", Name: "banana", Status: differ.StatusUnchanged},
	}
}

func TestApply_ByName_Ascending(t *testing.T) {
	res := sorter.ByField(makeResults(), sorter.ByName)
	names := []string{res[0].Name, res[1].Name, res[2].Name, res[3].Name}
	want := []string{"apple", "banana", "mango", "zebra"}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("pos %d: got %q want %q", i, names[i], want[i])
		}
	}
}

func TestApply_ByType_Ascending(t *testing.T) {
	res := sorter.ByField(makeResults(), sorter.ByType)
	if res[0].Type != "cluster" {
		t.Errorf("expected first type=cluster, got %q", res[0].Type)
	}
}

func TestApply_ByStatus_Ascending(t *testing.T) {
	res := sorter.ByField(makeResults(), sorter.ByStatus)
	// statuses sorted lexicographically: added < modified < removed < unchanged
	if res[0].Status != differ.StatusAdded {
		t.Errorf("expected first status=added, got %q", res[0].Status)
	}
}

func TestApply_Descending(t *testing.T) {
	res := sorter.ApplyWithOptions(makeResults(),
		sorter.WithField(sorter.ByName),
		sorter.WithDescending(),
	)
	if res[0].Name != "zebra" {
		t.Errorf("expected first name=zebra, got %q", res[0].Name)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	orig := makeResults()
	firstName := orig[0].Name
	sorter.ByField(orig, sorter.ByName)
	if orig[0].Name != firstName {
		t.Error("original slice was mutated")
	}
}

func TestApply_Empty(t *testing.T) {
	res := sorter.ByField(nil, sorter.ByName)
	if len(res) != 0 {
		t.Errorf("expected empty result, got %d items", len(res))
	}
}
