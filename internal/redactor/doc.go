// Package redactor provides field-level redaction for envoy-diff results.
//
// Sensitive field names (e.g. "tls_certificate", "token") are matched
// case-insensitively against keys in the Before/After JSON bodies of each
// differ.Result. Matched values are replaced with the placeholder "[REDACTED]"
// so that audit output can be shared without exposing secrets.
//
// Usage:
//
//	results = redactor.Apply(results, redactor.DefaultOptions())
package redactor
