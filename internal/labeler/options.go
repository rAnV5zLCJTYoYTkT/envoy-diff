package labeler

import "github.com/yourorg/envoy-diff/internal/differ"

// Option is a functional option for configuring a Labeler.
type Option func(*Labeler)

// WithRules replaces the default rule set with the provided rules.
func WithRules(rules []Rule) Option {
	return func(l *Labeler) {
		l.rules = rules
	}
}

// WithExtraRules appends rules to the existing rule set without replacing them.
func WithExtraRules(rules ...Rule) Option {
	return func(l *Labeler) {
		l.rules = append(l.rules, rules...)
	}
}

// NewWithOptions constructs a Labeler starting from DefaultRules and applying
// each option in order.
func NewWithOptions(opts ...Option) *Labeler {
	l := &Labeler{rules: DefaultRules()}
	for _, o := range opts {
		o(l)
	}
	return l
}

// MatchType returns a predicate that matches results whose Type equals t
// (case-insensitive).
func MatchType(t string) func(differ.Result) bool {
	return func(r differ.Result) bool {
		return len(r.Type) > 0 && len(t) > 0 &&
			len(r.Type) == len(t) &&
			stringsEqualFold(r.Type, t)
	}
}

// MatchStatus returns a predicate that matches results with the given status.
func MatchStatus(s differ.DiffStatus) func(differ.Result) bool {
	return func(r differ.Result) bool {
		return r.Status == s
	}
}

func stringsEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
