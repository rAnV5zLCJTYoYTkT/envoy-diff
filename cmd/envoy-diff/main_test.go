package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const baseSnap = `{
  "cluster": [
    {"name": "cluster-a", "connect_timeout": "1s"},
    {"name": "cluster-b", "connect_timeout": "2s"}
  ]
}`

const headSnap = `{
  "cluster": [
    {"name": "cluster-a", "connect_timeout": "5s"},
    {"name": "cluster-c", "connect_timeout": "1s"}
  ]
}`

func writeTempSnap(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func TestMain_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	base := writeTempSnap(t, "base.json", baseSnap)
	head := writeTempSnap(t, "head.json", headSnap)

	cmd := exec.Command("go", "run", ".",
		"--base", base,
		"--head", head,
		"--format", "text",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\noutput: %s", err, out)
	}
	outStr := string(out)
	for _, want := range []string{"cluster-a", "cluster-b", "cluster-c"} {
		if !containsStr(outStr, want) {
			t.Errorf("output missing %q; got:\n%s", want, outStr)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
