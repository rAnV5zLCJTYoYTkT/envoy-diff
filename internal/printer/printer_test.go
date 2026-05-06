package printer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/printer"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "cluster", Name: "api-cluster", Status: "added", Tags: map[string]string{"label": "infra"}},
		{Type: "listener", Name: "public-listener", Status: "removed", Tags: map[string]string{"breaking": "true"}},
		{Type: "route", Name: "default-route", Status: "modified", Tags: map[string]string{}},
		{Type: "cluster", Name: "old-cluster", Status: "unchanged", Tags: map[string]string{}},
	}
}

func TestWriteTable_HidesUnchangedByDefault(t *testing.T) {
	var buf bytes.Buffer
	opts := printer.DefaultOptions()
	err := printer.WriteTable(&buf, makeResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "unchanged") {
		t.Error("expected unchanged rows to be hidden by default")
	}
	if !strings.Contains(out, "added") {
		t.Error("expected added row to be present")
	}
}

func TestWriteTable_ShowUnchanged(t *testing.T) {
	var buf bytes.Buffer
	opts := printer.DefaultOptions()
	opts.ShowUnchanged = true
	err := printer.WriteTable(&buf, makeResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "unchanged") {
		t.Error("expected unchanged row when ShowUnchanged=true")
	}
}

func TestWriteTable_TruncatesLongNames(t *testing.T) {
	results := []differ.Result{
		{Type: "cluster", Name: strings.Repeat("x", 80), Status: "added", Tags: map[string]string{}},
	}
	var buf bytes.Buffer
	opts := printer.DefaultOptions()
	opts.MaxNameWidth = 20
	if err := printer.WriteTable(&buf, results, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), strings.Repeat("x", 80)) {
		t.Error("expected long name to be truncated")
	}
	if !strings.Contains(buf.String(), "...") {
		t.Error("expected ellipsis in truncated name")
	}
}

func TestWriteTable_NoDiffs(t *testing.T) {
	results := []differ.Result{
		{Type: "cluster", Name: "c", Status: "unchanged", Tags: map[string]string{}},
	}
	var buf bytes.Buffer
	if err := printer.WriteTable(&buf, results, printer.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no differences found") {
		t.Error("expected no-differences message when all results are unchanged")
	}
}

func TestWriteTable_BreakingNotes(t *testing.T) {
	results := []differ.Result{
		{Type: "listener", Name: "l", Status: "removed", Tags: map[string]string{"breaking": "true", "label": "edge"}},
	}
	var buf bytes.Buffer
	if err := printer.WriteTable(&buf, results, printer.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "BREAKING") {
		t.Error("expected BREAKING annotation in notes")
	}
	if !strings.Contains(out, "edge") {
		t.Error("expected label annotation in notes")
	}
}
