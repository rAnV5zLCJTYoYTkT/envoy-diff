package differ

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/snapshot"
)

func makeSnapshot(t *testing.T, resources map[string]map[string]string) *snapshot.Snapshot {
	t.Helper()
	s := snapshot.New()
	for typ, items := range resources {
		for name, val := range items {
			s.Add(typ, name, val)
		}
	}
	return s
}

func TestCompare_Added(t *testing.T) {
	a := makeSnapshot(t, nil)
	b := makeSnapshot(t, map[string]map[string]string{
		"cluster": {"c1": `{"name":"c1"}`},
	})
	res := Compare(a, b)
	if len(res) != 1 || res[0].Status != Added {
		t.Fatalf("expected 1 added, got %+v", res)
	}
}

func TestCompare_Removed(t *testing.T) {
	a := makeSnapshot(t, map[string]map[string]string{
		"cluster": {"c1": `{"name":"c1"}`},
	})
	b := makeSnapshot(t, nil)
	res := Compare(a, b)
	if len(res) != 1 || res[0].Status != Removed {
		t.Fatalf("expected 1 removed, got %+v", res)
	}
}

func TestCompare_Modified(t *testing.T) {
	a := makeSnapshot(t, map[string]map[string]string{
		"listener": {"l1": `{"port":80}`},
	})
	b := makeSnapshot(t, map[string]map[string]string{
		"listener": {"l1": `{"port":443}`},
	})
	res := Compare(a, b)
	if len(res) != 1 || res[0].Status != Modified {
		t.Fatalf("expected 1 modified, got %+v", res)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	a := makeSnapshot(t, map[string]map[string]string{
		"route": {"r1": `{"prefix":"/"}`},
	})
	b := makeSnapshot(t, map[string]map[string]string{
		"route": {"r1": `{"prefix":"/"}`},
	})
	res := Compare(a, b)
	if len(res) != 1 || res[0].Status != Unchanged {
		t.Fatalf("expected 1 unchanged, got %+v", res)
	}
}

func TestCompare_MultiType(t *testing.T) {
	a := makeSnapshot(t, map[string]map[string]string{
		"cluster": {"c1": `{}`},
	})
	b := makeSnapshot(t, map[string]map[string]string{
		"cluster":  {"c1": `{}`},
		"listener": {"l1": `{}`},
	})
	res := Compare(a, b)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}
