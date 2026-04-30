package loader_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-diff/internal/loader"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestFromBytes_Valid(t *testing.T) {
	data := `{
		"cluster": [{"name": "cluster-a", "connect_timeout": "1s"}],
		"listener": [{"name": "listener-1"}]
	}`
	snap, err := loader.FromBytes([]byte(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := snap.ResourceNames("cluster")
	if len(names) != 1 || names[0] != "cluster-a" {
		t.Errorf("expected [cluster-a], got %v", names)
	}
}

func TestFromBytes_MissingName(t *testing.T) {
	data := `{"cluster": [{"connect_timeout": "1s"}]}`
	_, err := loader.FromBytes([]byte(data))
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestFromBytes_InvalidJSON(t *testing.T) {
	_, err := loader.FromBytes([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFromFile_Valid(t *testing.T) {
	content := `{"endpoint": [{"name": "ep-1", "address": "127.0.0.1"}]}`
	path := writeTemp(t, content)
	snap, err := loader.FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := snap.Resource("endpoint", "ep-1")
	if raw == nil {
		t.Fatal("expected resource ep-1, got nil")
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal resource: %v", err)
	}
}

func TestFromFile_NotFound(t *testing.T) {
	_, err := loader.FromFile("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
