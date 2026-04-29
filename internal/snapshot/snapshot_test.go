package snapshot_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/snapshot"
)

func TestNew(t *testing.T) {
	s := snapshot.New("production")
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if s.Environment != "production" {
		t.Errorf("expected environment %q, got %q", "production", s.Environment)
	}
	if s.Resources == nil {
		t.Error("expected resources map to be initialized")
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestAddResource(t *testing.T) {
	s := snapshot.New("staging")
	res := snapshot.Resource{
		Name:    "my-cluster",
		Version: "v1",
		Body:    map[string]interface{}{"connect_timeout": "5s"},
	}
	s.AddResource(snapshot.ResourceTypeClusters, res)

	clusters := s.Resources[snapshot.ResourceTypeClusters]
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	if clusters[0].Name != "my-cluster" {
		t.Errorf("expected name %q, got %q", "my-cluster", clusters[0].Name)
	}
}

func TestResourceNames(t *testing.T) {
	s := snapshot.New("dev")
	s.AddResource(snapshot.ResourceTypeListeners, snapshot.Resource{Name: "listener-a"})
	s.AddResource(snapshot.ResourceTypeListeners, snapshot.Resource{Name: "listener-b"})
	s.AddResource(snapshot.ResourceTypeClusters, snapshot.Resource{Name: "cluster-x"})

	names := s.ResourceNames(snapshot.ResourceTypeListeners)
	if len(names) != 2 {
		t.Errorf("expected 2 listener names, got %d", len(names))
	}
	if _, ok := names["listener-a"]; !ok {
		t.Error("expected listener-a in names")
	}
	if _, ok := names["listener-b"]; !ok {
		t.Error("expected listener-b in names")
	}

	clusterNames := s.ResourceNames(snapshot.ResourceTypeClusters)
	if len(clusterNames) != 1 {
		t.Errorf("expected 1 cluster name, got %d", len(clusterNames))
	}
}

func TestResourceNamesEmptyType(t *testing.T) {
	s := snapshot.New("dev")
	names := s.ResourceNames(snapshot.ResourceTypeRoutes)
	if len(names) != 0 {
		t.Errorf("expected 0 names for empty type, got %d", len(names))
	}
}
