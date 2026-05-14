package streamer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/envoy-diff/internal/differ"
	"github.com/example/envoy-diff/internal/streamer"
)

func makeResult(status differ.DiffStatus, typ, name string) differ.Result {
	return differ.Result{
		Type:   typ,
		Name:   name,
		Status: status,
		Labels: map[string]string{},
		Tags:   []string{},
	}
}

func TestWriteText_SingleResult(t *testing.T) {
	var buf bytes.Buffer
	s := streamer.New(&buf, streamer.FormatText)
	_ = s.Open()
	_ = s.Write(makeResult(differ.Added, "listener", "ingress"))
	_ = s.Close()

	out := buf.String()
	if !strings.Contains(out, "[added]") {
		t.Errorf("expected status in output, got: %s", out)
	}
	if !strings.Contains(out, "listener/ingress") {
		t.Errorf("expected type/name in output, got: %s", out)
	}
}

func TestWriteText_MultipleResults(t *testing.T) {
	var buf bytes.Buffer
	s := streamer.New(&buf, streamer.FormatText)
	_ = s.Open()
	_ = s.Write(makeResult(differ.Added, "cluster", "a"))
	_ = s.Write(makeResult(differ.Removed, "cluster", "b"))
	_ = s.Close()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %s", len(lines), buf.String())
	}
}

func TestWriteJSON_ProducesValidArray(t *testing.T) {
	var buf bytes.Buffer
	s := streamer.New(&buf, streamer.FormatJSON)
	_ = s.Open()
	_ = s.Write(makeResult(differ.Modified, "route", "r1"))
	_ = s.Write(makeResult(differ.Unchanged, "route", "r2"))
	_ = s.Close()

	var results []differ.Result
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, buf.String())
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestNewWithOptions_NilWriterPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil writer")
		}
	}()
	streamer.NewWithOptions(streamer.DefaultOptions())
}

func TestNewWithOptions_ValidOptions(t *testing.T) {
	var buf bytes.Buffer
	opts := streamer.WithWriter(streamer.WithFormat(streamer.DefaultOptions(), streamer.FormatText), &buf)
	s := streamer.NewWithOptions(opts)
	_ = s.Open()
	_ = s.Write(makeResult(differ.Added, "endpoint", "ep0"))
	_ = s.Close()
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestOpen_Idempotent(t *testing.T) {
	var buf bytes.Buffer
	s := streamer.New(&buf, streamer.FormatJSON)
	_ = s.Open()
	_ = s.Open() // second call should be a no-op
	_ = s.Close()

	// Should still be a valid (empty) JSON array.
	var results []differ.Result
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON after double Open: %v\n%s", err, buf.String())
	}
}
