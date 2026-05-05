// Package patcher provides rule-based patching of envoy-diff results.
//
// Patch rules can suppress known-good differences (e.g. healthcheck routes
// that change frequently) or override labels for downstream consumers.
//
// Example usage:
//
//	p := patcher.New(patcher.DefaultSuppressRules())
//	patched := p.Apply(results)
package patcher
