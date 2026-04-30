// Package filter provides post-processing utilities for envoy-diff comparison
// results. It allows callers to narrow down a slice of diff results by xDS
// resource type, change status, or a substring match on the resource name.
//
// Typical usage:
//
//	opts := filter.Options{
//	    Types:         []string{"listeners", "clusters"},
//	    StatusInclude: []string{"added", "modified"},
//	    NameContains:  "prod",
//	}
//	filtered := filter.Apply(allResults, opts)
//
// All criteria are ANDed together; an empty value for any field disables that
// particular criterion so that all entries pass it.
package filter
