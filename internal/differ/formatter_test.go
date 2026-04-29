package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envoy-diff/internal/differ"
)

func buildResult(diffs []differ.ResourceDiff) *differ.Result {
	return &differ.Result{Diffs: diffs}
}

func TestRenderText_Added(t *testing.T) {
	result := buildResult([]differ.ResourceDiff{
		{Type: "cluster", Name: "my-cluster", Status: differ.StatusAdded, Right: `{"name":"my-cluster"}`},
	})
	var buf bytes.Buffer
	if err := differ.Render(&buf, result, differ.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "+ [cluster] my-cluster") {
		t.Errorf("unexpected output: %s", output)
	}
}

func TestRenderText_Removed(t *testing.T) {
	result := buildResult([]differ.ResourceDiff{
		{Type: "listener", Name: "l1", Status: differ.StatusRemoved, Left: `{"name":"l1"}`},
	})
	var buf bytes.Buffer
	if err := differ.Render(&buf, result, differ.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "- [listener] l1") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderText_NoDiffs(t *testing.T) {
	result := buildResult(nil)
	var buf bytes.Buffer
	if err := differ.Render(&buf, result, differ.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderJSON_Valid(t *testing.T) {
	result := buildResult([]differ.ResourceDiff{
		{Type: "cluster", Name: "c1", Status: differ.StatusModified, Left: `{}`, Right: `{"x":1}`},
	})
	var buf bytes.Buffer
	if err := differ.Render(&buf, result, differ.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"Status"`) {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestRenderUnsupportedFormat(t *testing.T) {
	result := buildResult(nil)
	var buf bytes.Buffer
	err := differ.Render(&buf, result, differ.Format("yaml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
