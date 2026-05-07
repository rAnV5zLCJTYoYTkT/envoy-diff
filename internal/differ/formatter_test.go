package differ

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildResult(status Status) Result {
	return Result{
		Type:     "cluster",
		Name:     "my-cluster",
		Status:   status,
		OldValue: `{"name":"my-cluster"}`,
		NewValue: `{"name":"my-cluster","extra":true}`,
	}
}

func TestRenderText_Added(t *testing.T) {
	var buf bytes.Buffer
	err := Render(&buf, []Result{buildResult(Added)}, "text")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "+ [cluster]") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderText_Removed(t *testing.T) {
	var buf bytes.Buffer
	_ = Render(&buf, []Result{buildResult(Removed)}, "text")
	if !strings.Contains(buf.String(), "- [cluster]") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	_ = Render(&buf, nil, "text")
	if !strings.Contains(buf.String(), "No diffs") {
		t.Errorf("expected no-diffs message, got: %s", buf.String())
	}
}

func TestRenderJSON_Valid(t *testing.T) {
	var buf bytes.Buffer
	_ = Render(&buf, []Result{buildResult(Modified)}, "json")
	var out []Result
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 result, got %d", len(out))
	}
}

func TestRenderJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	_ = Render(&buf, nil, "json")
	var out []Result
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
