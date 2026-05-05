package annotator_test

import (
	"strings"
	"testing"

	"github.com/your-org/envoy-diff/internal/annotator"
	"github.com/your-org/envoy-diff/internal/differ"
)

func makeResults() []differ.DiffResult {
	return []differ.DiffResult{
		{Type: "envoy.config.listener.v3.Listener", Name: "listener-a", Status: differ.StatusAdded},
		{Type: "envoy.config.cluster.v3.Cluster", Name: "cluster-b", Status: differ.StatusRemoved},
		{Type: "envoy.config.route.v3.RouteConfiguration", Name: "route-c", Status: differ.StatusModified},
		{Type: "envoy.config.endpoint.v3.ClusterLoadAssignment", Name: "ep-d", Status: differ.StatusUnchanged},
	}
}

func TestAnnotate_DefaultRules(t *testing.T) {
	a := annotator.New(annotator.DefaultRules()...)
	annotated := a.Annotate(makeResults())

	if len(annotated) != 4 {
		t.Fatalf("expected 4 annotated results, got %d", len(annotated))
	}
}

func TestAnnotate_BreakingChange(t *testing.T) {
	a := annotator.New(annotator.BreakingChangeRule())
	annotated := a.Annotate(makeResults())

	var found *annotator.AnnotatedResult
	for i := range annotated {
		if annotated[i].Result.Status == differ.StatusRemoved {
			found = &annotated[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected a removed result")
	}
	if found.Annotation.Label != "breaking" {
		t.Errorf("expected label 'breaking', got %q", found.Annotation.Label)
	}
	if !strings.Contains(found.Annotation.Summary, "cluster-b") {
		t.Errorf("expected summary to mention resource name, got %q", found.Annotation.Summary)
	}
}

func TestAnnotate_ResourceTypeTag(t *testing.T) {
	a := annotator.New(annotator.ResourceTypeRule())
	annotated := a.Annotate(makeResults())

	for _, ar := range annotated {
		found := false
		for _, tag := range ar.Annotation.Tags {
			if strings.HasPrefix(tag, "type:") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected type tag for %q, got tags %v", ar.Result.Name, ar.Annotation.Tags)
		}
	}
}

func TestAnnotate_DefaultLabelFallback(t *testing.T) {
	// No rules — label should fall back to status string.
	a := annotator.New()
	annotated := a.Annotate(makeResults())

	for _, ar := range annotated {
		if ar.Annotation.Label == "" {
			t.Errorf("expected non-empty label for %q", ar.Result.Name)
		}
	}
}

func TestAnnotate_CustomRule(t *testing.T) {
	customRule := func(r differ.DiffResult) (annotator.Annotation, bool) {
		if r.Status == differ.StatusAdded {
			return annotator.Annotation{Tags: []string{"new-resource"}}, true
		}
		return annotator.Annotation{}, false
	}

	a := annotator.New(customRule)
	annotated := a.Annotate(makeResults())

	for _, ar := range annotated {
		if ar.Result.Status == differ.StatusAdded {
			found := false
			for _, tag := range ar.Annotation.Tags {
				if tag == "new-resource" {
					found = true
				}
			}
			if !found {
				t.Errorf("expected 'new-resource' tag on added result %q", ar.Result.Name)
			}
		}
	}
}
