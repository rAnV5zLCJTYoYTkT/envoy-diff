package exporter_test

import (
	"testing"

	"github.com/example/envoy-diff/internal/exporter"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected exporter.Format
	}{
		{"json", exporter.FormatJSON},
		{"csv", exporter.FormatCSV},
		{"text", exporter.FormatText},
	}
	for _, tc := range cases {
		f, err := exporter.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if f != tc.expected {
			t.Errorf("ParseFormat(%q) = %v, want %v", tc.input, f, tc.expected)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := exporter.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestFormat_String(t *testing.T) {
	cases := []struct {
		format   exporter.Format
		expected string
	}{
		{exporter.FormatJSON, "json"},
		{exporter.FormatCSV, "csv"},
		{exporter.FormatText, "text"},
	}
	for _, tc := range cases {
		if got := tc.format.String(); got != tc.expected {
			t.Errorf("Format.String() = %q, want %q", got, tc.expected)
		}
	}
}

func TestFormats_Count(t *testing.T) {
	if got := exporter.Formats(); len(got) != 3 {
		t.Errorf("expected 3 formats, got %d", len(got))
	}
}
