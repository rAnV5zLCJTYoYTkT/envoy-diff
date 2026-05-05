package baseline

import "github.com/example/envoy-diff/internal/differ"

// SaveOption configures a Save call.
type SaveOption func(*saveConfig)

type saveConfig struct {
	label string
}

// WithLabel sets a human-readable label stored alongside the baseline.
func WithLabel(label string) SaveOption {
	return func(c *saveConfig) {
		c.label = label
	}
}

// SaveWithOptions is a variadic wrapper around Save that accepts SaveOptions.
func SaveWithOptions(path string, results []differ.DiffResult, opts ...SaveOption) error {
	cfg := &saveConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return Save(path, cfg.label, results)
}
