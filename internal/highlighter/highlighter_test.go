package highlighter_test

import (
	"testing"

	"github.com/envoy-diff/internal/differ"
	"github.com/envoy-diff/internal/highlighter"
)

func makeResults(status, before, after string) []differ.Result {
	return []differ.Result{
		{Type: "listener", Name: "ingress", Status: status, Before: before, After: after},
	}
}

func findField(fields []highlighter.FieldDiff, name string) (highlighter.FieldDiff, bool) {
	for _, f := range fields {
		if f.Field == name {
			return f, true
		}
	}
	return highlighter.FieldDiff{}, false
}

func TestApply_Modified_DetectsChangedField(t *testing.T) {
	before := `{"port":8080,"name":"old"}`
	after := `{"port":9090,"name":"old"}`
	results := makeResults("modified", before, after)

	highlights := highlighter.Apply(results)
	if len(highlights) != 1 {
		t.Fatalf("expected 1 highlight, got %d", len(highlights))
	}
	h := highlights[0]
	if len(h.Fields) != 1 {
		t.Fatalf("expected 1 field diff, got %d", len(h.Fields))
	}
	f, ok := findField(h.Fields, "port")
	if !ok {
		t.Fatal("expected field 'port' in diff")
	}
	if f.Before != "8080" || f.After != "9090" {
		t.Errorf("unexpected values: before=%q after=%q", f.Before, f.After)
	}
}

func TestApply_Modified_MultipleFields(t *testing.T) {
	before := `{"port":8080,"tls":false}`
	after := `{"port":9090,"tls":true}`
	results := makeResults("modified", before, after)

	highlights := highlighter.Apply(results)
	if len(highlights[0].Fields) != 2 {
		t.Errorf("expected 2 field diffs, got %d", len(highlights[0].Fields))
	}
}

func TestApply_Unchanged_NoFields(t *testing.T) {
	body := `{"port":8080}`
	results := makeResults("unchanged", body, body)

	highlights := highlighter.Apply(results)
	if len(highlights[0].Fields) != 0 {
		t.Errorf("expected no fields for unchanged, got %d", len(highlights[0].Fields))
	}
}

func TestApply_Added_NoFields(t *testing.T) {
	results := makeResults("added", "", `{"port":8080}`)
	highlights := highlighter.Apply(results)
	if len(highlights[0].Fields) != 0 {
		t.Errorf("expected no fields for added, got %d", len(highlights[0].Fields))
	}
}

func TestApply_Removed_NoFields(t *testing.T) {
	results := makeResults("removed", `{"port":8080}`, "")
	highlights := highlighter.Apply(results)
	if len(highlights[0].Fields) != 0 {
		t.Errorf("expected no fields for removed, got %d", len(highlights[0].Fields))
	}
}

func TestApply_Modified_FieldAdded(t *testing.T) {
	before := `{"port":8080}`
	after := `{"port":8080,"tls":true}`
	results := makeResults("modified", before, after)

	highlights := highlighter.Apply(results)
	_, ok := findField(highlights[0].Fields, "tls")
	if !ok {
		t.Error("expected 'tls' to appear as a new field diff")
	}
}

func TestApply_PreservesMetadata(t *testing.T) {
	results := makeResults("modified", `{"a":1}`, `{"a":2}`)
	h := highlighter.Apply(results)[0]
	if h.Type != "listener" || h.Name != "ingress" || h.Status != "modified" {
		t.Errorf("metadata mismatch: type=%q name=%q status=%q", h.Type, h.Name, h.Status)
	}
}
