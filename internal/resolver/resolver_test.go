package resolver_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/resolver"
)

func TestResolve_ShortAlias(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"listener", "type.googleapis.com/envoy.config.listener.v3.Listener"},
		{"listeners", "type.googleapis.com/envoy.config.listener.v3.Listener"},
		{"cluster", "type.googleapis.com/envoy.config.cluster.v3.Cluster"},
		{"routes", "type.googleapis.com/envoy.config.route.v3.RouteConfiguration"},
		{"endpoint", "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment"},
		{"secret", "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := resolver.Resolve(tc.input)
			if got != tc.want {
				t.Errorf("Resolve(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	got := resolver.Resolve("LISTENER")
	want := "type.googleapis.com/envoy.config.listener.v3.Listener"
	if got != want {
		t.Errorf("Resolve(LISTENER) = %q; want %q", got, want)
	}
}

func TestResolve_FullURL_PassThrough(t *testing.T) {
	full := "type.googleapis.com/envoy.config.listener.v3.Listener"
	if got := resolver.Resolve(full); got != full {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestResolve_Unknown_PassThrough(t *testing.T) {
	input := "unknown-resource"
	if got := resolver.Resolve(input); got != input {
		t.Errorf("expected passthrough for unknown, got %q", got)
	}
}

func TestShortName_KnownType(t *testing.T) {
	got := resolver.ShortName("type.googleapis.com/envoy.config.listener.v3.Listener")
	if got != "listener" {
		t.Errorf("ShortName = %q; want \"listener\"", got)
	}
}

func TestShortName_UnknownURL(t *testing.T) {
	url := "type.googleapis.com/some.custom.v1.Widget"
	got := resolver.ShortName(url)
	if got != "Widget" {
		t.Errorf("ShortName = %q; want \"Widget\"", got)
	}
}

func TestKnownTypeURLs_NoDuplicates(t *testing.T) {
	urls := resolver.KnownTypeURLs()
	seen := make(map[string]int)
	for _, u := range urls {
		seen[u]++
	}
	for u, count := range seen {
		if count > 1 {
			t.Errorf("duplicate URL %q appeared %d times", u, count)
		}
	}
}

func TestKnownTypeURLs_AllAreFullURLs(t *testing.T) {
	for _, u := range resolver.KnownTypeURLs() {
		if !strings.Contains(u, "/") {
			t.Errorf("expected full URL, got %q", u)
		}
	}
}
