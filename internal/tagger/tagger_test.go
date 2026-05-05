package tagger_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/differ"
	"github.com/your-org/envoy-diff/internal/tagger"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Type: "type.googleapis.com/envoy.config.listener.v3.Listener", Name: "ingress", Status: differ.Added},
		{Type: "type.googleapis.com/envoy.config.cluster.v3.Cluster", Name: "backend", Status: differ.Removed},
		{Type: "type.googleapis.com/envoy.config.route.v3.RouteConfiguration", Name: "default", Status: differ.Modified},
		{Type: "type.googleapis.com/envoy.config.listener.v3.Listener", Name: "egress", Status: differ.Unchanged},
	}
}

func findTag(tags []tagger.Tag, key string) (string, bool) {
	for _, t := range tags {
		if t.Key == key {
			return t.Value, true
		}
	}
	return "", false
}

func TestApply_DefaultRules_ChangeTag(t *testing.T) {
	tgr := tagger.New(tagger.DefaultRules())
	tagged := tgr.Apply(makeResults())

	expected := map[string]string{
		"ingress": "added",
		"backend": "removed",
		"default": "modified",
	}
	for _, tr := range tagged {
		if want, ok := expected[tr.Name]; ok {
			v, found := findTag(tr.Tags, "change")
			if !found {
				t.Errorf("result %q missing 'change' tag", tr.Name)
				continue
			}
			if v != want {
				t.Errorf("result %q: want change=%q, got %q", tr.Name, want, v)
			}
		}
	}
}

func TestApply_DefaultRules_ResourceFamily(t *testing.T) {
	tgr := tagger.New(tagger.DefaultRules())
	tagged := tgr.Apply(makeResults())

	for _, tr := range tagged {
		switch tr.Name {
		case "ingress", "egress":
			v, found := findTag(tr.Tags, "resource-family")
			if !found || v != "listener" {
				t.Errorf("%q: expected resource-family=listener, got %q", tr.Name, v)
			}
		case "backend":
			v, found := findTag(tr.Tags, "resource-family")
			if !found || v != "cluster" {
				t.Errorf("%q: expected resource-family=cluster, got %q", tr.Name, v)
			}
		}
	}
}

func TestApply_CustomRule(t *testing.T) {
	rules := []tagger.Rule{
		{
			Match: func(r differ.Result) bool { return r.Name == "ingress" },
			Tags:  []tagger.Tag{{Key: "env", Value: "prod"}},
		},
	}
	tgr := tagger.New(rules)
	tagged := tgr.Apply(makeResults())

	for _, tr := range tagged {
		v, found := findTag(tr.Tags, "env")
		if tr.Name == "ingress" {
			if !found || v != "prod" {
				t.Errorf("ingress: expected env=prod")
			}
		} else if found {
			t.Errorf("%q: unexpected env tag", tr.Name)
		}
	}
}

func TestApply_EmptyResults(t *testing.T) {
	tgr := tagger.New(tagger.DefaultRules())
	tagged := tgr.Apply(nil)
	if len(tagged) != 0 {
		t.Errorf("expected empty slice, got %d", len(tagged))
	}
}
