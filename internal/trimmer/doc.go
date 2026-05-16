// Package trimmer provides a lightweight filter that discards diff results
// whose body content falls below a configurable size threshold.
//
// This is useful when downstream rendering or auditing pipelines generate
// noise from resources that technically differ but whose actual serialised
// diff is empty or trivially small (e.g. a single whitespace character).
//
// Usage:
//
//	opts := trimmer.DefaultOptions()
//	opts.MinBodyLen = 10
//	filtered := trimmer.Apply(results, opts)
package trimmer
