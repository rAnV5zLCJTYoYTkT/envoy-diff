package streamer

import "io"

// Options configures a Streamer.
type Options struct {
	Format Format
	Writer io.Writer
}

// DefaultOptions returns options with text format and a nil writer.
// Callers must supply a non-nil Writer before use.
func DefaultOptions() Options {
	return Options{
		Format: FormatText,
	}
}

// WithFormat returns a copy of opts with the format set to f.
func WithFormat(opts Options, f Format) Options {
	opts.Format = f
	return opts
}

// WithWriter returns a copy of opts with the writer set to w.
func WithWriter(opts Options, w io.Writer) Options {
	opts.Writer = w
	return opts
}

// NewWithOptions constructs a Streamer from the provided Options.
// It panics if Options.Writer is nil.
func NewWithOptions(opts Options) *Streamer {
	if opts.Writer == nil {
		panic("streamer: Options.Writer must not be nil")
	}
	return New(opts.Writer, opts.Format)
}
