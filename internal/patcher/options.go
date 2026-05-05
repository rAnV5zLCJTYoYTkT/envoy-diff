package patcher

// Option is a functional option for configuring a Patcher.
type Option func(*Patcher)

// WithRules appends additional rules to the patcher.
func WithRules(rules ...Rule) Option {
	return func(p *Patcher) {
		p.rules = append(p.rules, rules...)
	}
}

// NewWithOptions constructs a Patcher using the provided options.
func NewWithOptions(opts ...Option) *Patcher {
	p := &Patcher{}
	for _, o := range opts {
		o(p)
	}
	return p
}

// DefaultSuppressRules returns a set of commonly used suppression rules
// for well-known transient or expected diff patterns.
func DefaultSuppressRules() []Rule {
	return []Rule{
		{
			NameContains:   "healthcheck",
			SuppressStatus: true,
		},
		{
			NameContains:   "stats",
			SuppressStatus: true,
		},
	}
}
