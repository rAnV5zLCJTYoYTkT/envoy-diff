// Package highlighter provides field-level diff highlighting for Envoy xDS
// resource changes. It operates on differ.Result slices and produces
// Highlight values that identify which JSON fields changed between the
// before and after snapshots of a modified resource.
//
// Usage:
//
//	highlights := highlighter.Apply(results)
//	for _, h := range highlights {
//		for _, f := range h.Fields {
//			fmt.Printf("%s.%s: %q -> %q\n", h.Name, f.Field, f.Before, f.After)
//		}
//	}
package highlighter
