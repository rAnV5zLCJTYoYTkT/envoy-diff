package comparator_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envoy-diff/internal/comparator"
	"github.com/your-org/envoy-diff/internal/snapshot"
)

func makeSnap(resources map[string]string) *snapshot.Snapshot {
	s := snapshot.New()
	for name, body := range resources {
		s.Add("listener", name, []byte(body))
	}
	return s
}

func TestCompare_SequentialPairs(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"dev":     makeSnap(map[string]string{"a": `{"x":1}`}),
		"staging": makeSnap(map[string]string{"a": `{"x":2}`}),
		"prod":    makeSnap(map[string]string{"a": `{"x":2}`}),
	}
	r, err := comparator.Compare(envs, comparator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(r.Pairs))
	}
	if r.Pairs[0].Left != "dev" || r.Pairs[0].Right != "staging" {
		t.Errorf("unexpected first pair: %s→%s", r.Pairs[0].Left, r.Pairs[0].Right)
	}
}

func TestCompare_AllPairs(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"a": makeSnap(nil),
		"b": makeSnap(nil),
		"c": makeSnap(nil),
	}
	opts := comparator.Options{Sequential: false}
	r, err := comparator.Compare(envs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Pairs) != 3 {
		t.Fatalf("expected 3 pairs for 3 envs, got %d", len(r.Pairs))
	}
}

func TestCompare_TooFewEnvironments(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"only": makeSnap(nil),
	}
	_, err := comparator.Compare(envs, comparator.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for single environment, got nil")
	}
}

func TestWriteText_ContainsPairHeader(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"dev":  makeSnap(map[string]string{"svc": `{"port":80}`}),
		"prod": makeSnap(map[string]string{"svc": `{"port":443}`}),
	}
	r, _ := comparator.Compare(envs, comparator.DefaultOptions())
	var buf bytes.Buffer
	comparator.WriteText(&buf, r)
	if !strings.Contains(buf.String(), "dev → prod") {
		t.Errorf("expected header in output, got:\n%s", buf.String())
	}
}

func TestWriteJSON_Valid(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"x": makeSnap(nil),
		"y": makeSnap(nil),
	}
	r, _ := comparator.Compare(envs, comparator.DefaultOptions())
	var buf bytes.Buffer
	if err := comparator.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "Pairs") {
		t.Errorf("expected JSON output to contain 'Pairs', got:\n%s", buf.String())
	}
}
