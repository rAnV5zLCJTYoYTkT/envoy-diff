// Package streamer provides a streaming diff writer that emits
// differ.Result values incrementally to an io.Writer as they are processed.
package streamer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/example/envoy-diff/internal/differ"
)

// Format controls how streamed results are encoded.
type Format int

const (
	FormatText Format = iota
	FormatJSON
)

// Streamer writes differ.Result values to a destination as they arrive.
type Streamer struct {
	w      io.Writer
	fmt    Format
	first  bool // used to manage JSON array commas
	opened bool
}

// New returns a Streamer that writes to w using the given format.
func New(w io.Writer, f Format) *Streamer {
	return &Streamer{w: w, fmt: f, first: true}
}

// Open writes any necessary preamble (e.g. opening JSON array bracket).
func (s *Streamer) Open() error {
	if s.opened {
		return nil
	}
	s.opened = true
	if s.fmt == FormatJSON {
		_, err := fmt.Fprintln(s.w, "[")
		return err
	}
	return nil
}

// Write encodes a single result and writes it to the underlying writer.
func (s *Streamer) Write(r differ.Result) error {
	if !s.opened {
		if err := s.Open(); err != nil {
			return err
		}
	}
	switch s.fmt {
	case FormatJSON:
		return s.writeJSON(r)
	default:
		return s.writeText(r)
	}
}

func (s *Streamer) writeJSON(r differ.Result) error {
	if !s.first {
		if _, err := fmt.Fprintln(s.w, ","); err != nil {
			return err
		}
	}
	s.first = false
	enc := json.NewEncoder(s.w)
	enc.SetIndent("  ", "  ")
	return enc.Encode(r)
}

func (s *Streamer) writeText(r differ.Result) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s/%s", r.Status, r.Type, r.Name))
	if len(r.Labels) > 0 {
		sb.WriteString(" labels:")
		for k, v := range r.Labels {
			sb.WriteString(fmt.Sprintf(" %s=%s", k, v))
		}
	}
	sb.WriteString("\n")
	_, err := fmt.Fprint(s.w, sb.String())
	return err
}

// Close writes any necessary postamble (e.g. closing JSON array bracket).
func (s *Streamer) Close() error {
	if s.fmt == FormatJSON && s.opened {
		_, err := fmt.Fprintln(s.w, "]")
		return err
	}
	return nil
}
