// Package redactor masks sensitive field values in diff results before
// rendering or exporting, preventing secrets from appearing in audit output.
package redactor

import (
	"encoding/json"
	"strings"

	"github.com/example/envoy-diff/internal/differ"
)

// DefaultSensitiveKeys is the set of JSON field names redacted by default.
var DefaultSensitiveKeys = []string{
	"secret", "password", "token", "private_key", "tls_certificate",
	"inline_string", "inline_bytes",
}

const redactedPlaceholder = "[REDACTED]"

// Options controls redactor behaviour.
type Options struct {
	SensitiveKeys []string
}

// DefaultOptions returns an Options with the default sensitive key list.
func DefaultOptions() Options {
	return Options{SensitiveKeys: DefaultSensitiveKeys}
}

// Apply redacts sensitive fields from the Before and After JSON bodies of each
// result. Results are modified in place; a copy is returned for convenience.
func Apply(results []differ.Result, opts Options) []differ.Result {
	keySet := make(map[string]struct{}, len(opts.SensitiveKeys))
	for _, k := range opts.SensitiveKeys {
		keySet[strings.ToLower(k)] = struct{}{}
	}

	out := make([]differ.Result, len(results))
	for i, r := range results {
		c := r.Clone()
		c.Before = redactJSON(c.Before, keySet)
		c.After = redactJSON(c.After, keySet)
		out[i] = c
	}
	return out
}

// redactJSON walks a JSON string and replaces values whose keys are sensitive.
func redactJSON(raw string, keys map[string]struct{}) string {
	if raw == "" {
		return raw
	}
	var node interface{}
	if err := json.Unmarshal([]byte(raw), &node); err != nil {
		return raw // not valid JSON – return as-is
	}
	redacted := walkNode(node, keys)
	b, err := json.Marshal(redacted)
	if err != nil {
		return raw
	}
	return string(b)
}

func walkNode(node interface{}, keys map[string]struct{}) interface{} {
	switch v := node.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if _, sensitive := keys[strings.ToLower(k)]; sensitive {
				v[k] = redactedPlaceholder
			} else {
				v[k] = walkNode(val, keys)
			}
		}
		return v
	case []interface{}:
		for i, elem := range v {
			v[i] = walkNode(elem, keys)
		}
		return v
	}
	return node
}
