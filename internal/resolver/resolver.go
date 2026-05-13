// Package resolver maps resource names to their canonical xDS type URLs,
// enabling downstream components to normalise type strings that may arrive
// in short-form (e.g. "listeners") or as full type URLs
// (e.g. "type.googleapis.com/envoy.config.listener.v3.Listener").
package resolver

import "strings"

// knownTypes is the canonical mapping from short alias → full type URL.
var knownTypes = map[string]string{
	"listener":  "type.googleapis.com/envoy.config.listener.v3.Listener",
	"listeners": "type.googleapis.com/envoy.config.listener.v3.Listener",
	"cluster":   "type.googleapis.com/envoy.config.cluster.v3.Cluster",
	"clusters":  "type.googleapis.com/envoy.config.cluster.v3.Cluster",
	"route":     "type.googleapis.com/envoy.config.route.v3.RouteConfiguration",
	"routes":    "type.googleapis.com/envoy.config.route.v3.RouteConfiguration",
	"endpoint":  "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment",
	"endpoints": "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment",
	"secret":    "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",
	"secrets":   "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",
}

// Resolve returns the canonical full type URL for the given input string.
// If the input is already a full URL (contains "/") it is returned as-is.
// If the input matches a known alias (case-insensitive) the canonical URL
// is returned. Otherwise the original input is returned unchanged.
func Resolve(input string) string {
	if strings.Contains(input, "/") {
		return input
	}
	if canonical, ok := knownTypes[strings.ToLower(input)]; ok {
		return canonical
	}
	return input
}

// ShortName returns a human-friendly short name for a type URL.
// For known types the short alias is returned; otherwise the last
// path segment of the URL is used.
func ShortName(typeURL string) string {
	for alias, canonical := range knownTypes {
		// Prefer the singular form (no trailing 's') as the short name.
		if canonical == typeURL && !strings.HasSuffix(alias, "s") {
			return alias
		}
	}
	// Fall back to the last segment of the URL.
	parts := strings.Split(typeURL, "/")
	return parts[len(parts)-1]
}

// KnownTypeURLs returns all unique canonical type URLs registered in the
// resolver.
func KnownTypeURLs() []string {
	seen := make(map[string]struct{})
	for _, v := range knownTypes {
		seen[v] = struct{}{}
	}
	urls := make([]string, 0, len(seen))
	for u := range seen {
		urls = append(urls, u)
	}
	return urls
}
