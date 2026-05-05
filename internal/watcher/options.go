package watcher

import "time"

// Option is a functional option for configuring a Watcher.
type Option func(*Watcher)

// WithInterval sets the polling interval.
func WithInterval(d time.Duration) Option {
	return func(w *Watcher) {
		if d > 0 {
			w.interval = d
		}
	}
}

// NewWithOptions creates a Watcher with functional options applied.
// The default polling interval is 10 seconds if none is provided.
func NewWithOptions(baseFile, targetFile string, opts ...Option) *Watcher {
	w := &Watcher{
		baseFile:   baseFile,
		targetFile: targetFile,
		interval:   10 * time.Second,
	}
	for _, o := range opts {
		o(w)
	}
	return w
}
