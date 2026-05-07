package splitter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/splitter"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "Listener", Name: "l1", Status: differ.Added},
		{Type: "Listener", Name: "l2", Status: differ.Removed},
		{Type: "Cluster", Name: "c1", Status: differ.Modified},
		{Type: "Cluster", Name: "c2", Status: differ.Unchanged},
		{Type: "Route", Name: "r1", Status: differ.Added, Annotations: map[string]string{"env": "prod"}},
	}
}

func TestApply_ByType(t *testing.T) {
	buckets := splitter.Apply(makeResults(), splitter.ByType())
	if len(buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(buckets))
	}
	b := splitter.Find(buckets, "Listener")
	if b == nil || len(b.Results) != 2 {
		t.Errorf("expected 2 Listener results, got %v", b)
	}
}

func TestApply_ByStatus(t *testing.T) {
	buckets := splitter.Apply(makeResults(), splitter.ByStatus())
	added := splitter.Find(buckets, string(differ.Added))
	if added == nil || len(added.Results) != 2 {
		t.Errorf("expected 2 added results, got %v", added)
	}
	unchanged := splitter.Find(buckets, string(differ.Unchanged))
	if unchanged == nil || len(unchanged.Results) != 1 {
		t.Errorf("expected 1 unchanged result, got %v", unchanged)
	}
}

func TestApply_ByLabel_Labeled(t *testing.T) {
	buckets := splitter.Apply(makeResults(), splitter.ByLabel("env"))
	prod := splitter.Find(buckets, "prod")
	if prod == nil || len(prod.Results) != 1 {
		t.Errorf("expected 1 prod result, got %v", prod)
	}
	unlabeled := splitter.Find(buckets, "_unlabeled")
	if unlabeled == nil || len(unlabeled.Results) != 4 {
		t.Errorf("expected 4 unlabeled results, got %v", unlabeled)
	}
}

func TestApply_Empty(t *testing.T) {
	buckets := splitter.Apply(nil, splitter.ByType())
	if len(buckets) != 0 {
		t.Errorf("expected empty buckets, got %d", len(buckets))
	}
}

func TestKeys_Order(t *testing.T) {
	buckets := splitter.Apply(makeResults(), splitter.ByType())
	keys := splitter.Keys(buckets)
	if keys[0] != "Listener" {
		t.Errorf("expected first key Listener, got %s", keys[0])
	}
	if keys[1] != "Cluster" {
		t.Errorf("expected second key Cluster, got %s", keys[1])
	}
}

func TestFind_Missing(t *testing.T) {
	buckets := splitter.Apply(makeResults(), splitter.ByType())
	if splitter.Find(buckets, "Endpoint") != nil {
		t.Error("expected nil for missing key")
	}
}
