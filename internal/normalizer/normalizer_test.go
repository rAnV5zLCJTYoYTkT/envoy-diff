package normalizer_test

import (
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/normalizer"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Name: "  listener-a  ", Type: "Listener", Status: differ.Modified,
			Left:  `{"z":1,"a":2}`,
			Right: `{ "a" : 2 , "z" : 1 }`},
		{Name: "cluster-b", Type: "CLUSTER", Status: differ.Added,
			Left:  "",
			Right: `{"name":"cluster-b"}`},
		{Name: "route-c ", Type: "RouteConfiguration", Status: differ.Removed,
			Left:  `not-json`,
			Right: ""},
	}
}

func TestApply_TrimsName(t *testing.T) {
	results := makeResults()
	out := normalizer.Apply(results, normalizer.Options{TrimName: true})
	if out[0].Name != "listener-a" {
		t.Errorf("expected trimmed name, got %q", out[0].Name)
	}
	if out[2].Name != "route-c" {
		t.Errorf("expected trimmed name, got %q", out[2].Name)
	}
}

func TestApply_LowercasesType(t *testing.T) {
	results := makeResults()
	out := normalizer.Apply(results, normalizer.Options{LowercaseType: true})
	for _, r := range out {
		if r.Type != string([]rune(r.Type)) { // just ensure no panic
		}
	}
	if out[0].Type != "listener" {
		t.Errorf("expected lowercase type, got %q", out[0].Type)
	}
	if out[1].Type != "cluster" {
		t.Errorf("expected lowercase type, got %q", out[1].Type)
	}
}

func TestApply_CanonicalJSON_SortsKeys(t *testing.T) {
	results := makeResults()
	out := normalizer.Apply(results, normalizer.Options{CanonicalJSON: true})
	// Both sides should produce identical canonical form
	if out[0].Left != out[0].Right {
		t.Errorf("expected canonical JSON to match: left=%q right=%q", out[0].Left, out[0].Right)
	}
}

func TestApply_CanonicalJSON_InvalidJSON_Passthrough(t *testing.T) {
	results := makeResults()
	out := normalizer.Apply(results, normalizer.Options{CanonicalJSON: true})
	if out[2].Left != "not-json" {
		t.Errorf("expected invalid JSON to pass through unchanged, got %q", out[2].Left)
	}
}

func TestApply_DefaultOptions_DoesNotMutateOriginal(t *testing.T) {
	results := makeResults()
	origName := results[0].Name
	normalizer.Apply(results, normalizer.DefaultOptions())
	if results[0].Name != origName {
		t.Error("Apply must not mutate the original slice")
	}
}

func TestApply_EmptySlice(t *testing.T) {
	out := normalizer.Apply([]differ.Result{}, normalizer.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d results", len(out))
	}
}
