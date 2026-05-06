package sorter

// Option is a functional option for configuring sorting.
type Option func(*Options)

// WithField sets the field to sort by.
func WithField(f SortBy) Option {
	return func(o *Options) {
		o.Field = f
	}
}

// WithDescending reverses the sort order.
func WithDescending() Option {
	return func(o *Options) {
		o.Ascending = false
	}
}

// ApplyWithOptions builds an Options from functional opts and sorts results.
func ApplyWithOptions(results []differ.DiffResult, opts ...Option) []differ.DiffResult {
	cfg := DefaultOptions()
	for _, o := range opts {
		o(&cfg)
	}
	return Apply(results, cfg)
}
