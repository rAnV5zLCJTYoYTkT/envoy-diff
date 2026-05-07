package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/exporter"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "cluster", Name: "svc-a", Status: differ.StatusAdded},
		{Type: "listener", Name: "http", Status: differ.StatusRemoved},
		{Type: "cluster", Name: "svc-b", Status: differ.StatusUnchanged},
	}
}

func TestExport_JSON_Stdout(t *testing.T) {
	results := makeResults()
	opts := exporter.Options{Format: exporter.FormatJSON}
	if err := exporter.Export(results, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExport_JSON_File(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.json")
	results := makeResults()
	opts := exporter.Options{Format: exporter.FormatJSON, FilePath: tmp}
	if err := exporter.Export(results, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var got []differ.Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(got))
	}
}

func TestExport_CSV_File(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.csv")
	opts := exporter.Options{Format: exporter.FormatCSV, FilePath: tmp}
	if err := exporter.Export(makeResults(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(data), "type,name,status") {
		t.Errorf("missing csv header in output")
	}
	if !strings.Contains(string(data), "svc-a") {
		t.Errorf("missing expected row in csv output")
	}
}

func TestExport_Text_File(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.txt")
	opts := exporter.Options{Format: exporter.FormatText, FilePath: tmp}
	if err := exporter.Export(makeResults(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "[added]") {
		t.Errorf("expected status in text output")
	}
}

func TestExport_InvalidFormat(t *testing.T) {
	opts := exporter.Options{Format: "xml"}
	if err := exporter.Export(makeResults(), opts); err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExport_BadFilePath(t *testing.T) {
	opts := exporter.Options{Format: exporter.FormatJSON, FilePath: "/nonexistent/dir/out.json"}
	if err := exporter.Export(makeResults(), opts); err == nil {
		t.Error("expected error for bad file path")
	}
}

func TestExport_EmptyResults(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "empty.json")
	opts := exporter.Options{Format: exporter.FormatJSON, FilePath: tmp}
	if err := exporter.Export([]differ.Result{}, opts); err != nil {
		t.Fatalf("unexpected error for empty results: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var got []differ.Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json for empty results: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}
