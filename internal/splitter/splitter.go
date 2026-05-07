// Package splitter partitions a slice of differ.Result into multiple named
// buckets using a configurable key function. This is useful for splitting
// diff results by environment label, cluster, or any arbitrary dimension
// before further processing.
package splitter

import "github.com/yourorg/envoy-diff/internal/differ"

// Bucket holds a named collection of diff results.
type Bucket struct {
	Key     string
	Results []differ.Result
}

// KeyFn extracts a bucket key from a single Result.
type KeyFn func(r differ.Result) string

// ByType returns a KeyFn that splits results by resource type.
func ByType() KeyFn {
	return func(r differ.Result) string {
		return r.Type
	}
}

// ByStatus returns a KeyFn that splits results by diff status.
func ByStatus() KeyFn {
	return func(r differ.Result) string {
		return string(r.Status)
	}
}

// ByLabel returns a KeyFn that splits results by a specific annotation label.
// Results missing the label are placed in the "_unlabeled" bucket.
func ByLabel(label string) KeyFn {
	return func(r differ.Result) string {
		if r.Annotations == nil {
			return "_unlabeled"
		}
		if v, ok := r.Annotations[label]; ok {
			return v
		}
		return "_unlabeled"
	}
}

// Apply partitions results into buckets using the provided key function.
// The order of buckets follows first-seen key order.
func Apply(results []differ.Result, fn KeyFn) []Bucket {
	index := make(map[string]int)
	var buckets []Bucket

	for _, r := range results {
		key := fn(r)
		if i, exists := index[key]; exists {
			buckets[i].Results = append(buckets[i].Results, r)
		} else {
			index[key] = len(buckets)
			buckets = append(buckets, Bucket{Key: key, Results: []differ.Result{r}})
		}
	}

	return buckets
}

// Keys returns the ordered list of bucket keys from a split result.
func Keys(buckets []Bucket) []string {
	keys := make([]string, len(buckets))
	for i, b := range buckets {
		keys[i] = b.Key
	}
	return keys
}

// Find returns the bucket with the given key, or nil if not found.
func Find(buckets []Bucket, key string) *Bucket {
	for i := range buckets {
		if buckets[i].Key == key {
			return &buckets[i]
		}
	}
	return nil
}
