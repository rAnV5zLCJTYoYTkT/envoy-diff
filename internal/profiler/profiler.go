// Package profiler measures timing and resource usage across diff pipeline stages.
package profiler

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Entry records a single timed span.
type Entry struct {
	Stage    string
	Duration time.Duration
	Count    int // number of results processed
}

// Profile holds all collected entries.
type Profile struct {
	entries []Entry
	start   time.Time
}

// New creates a new Profile with the wall-clock start time set to now.
func New() *Profile {
	return &Profile{start: time.Now()}
}

// Record appends a timing entry for the given stage.
func (p *Profile) Record(stage string, d time.Duration, count int) {
	p.entries = append(p.entries, Entry{Stage: stage, Duration: d, Count: count})
}

// Total returns the sum of all recorded durations.
func (p *Profile) Total() time.Duration {
	var total time.Duration
	for _, e := range p.entries {
		total += e.Duration
	}
	return total
}

// Entries returns a copy of all recorded entries sorted by duration descending.
func (p *Profile) Entries() []Entry {
	out := make([]Entry, len(p.entries))
	copy(out, p.entries)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Duration > out[j].Duration
	})
	return out
}

// Write prints a human-readable profile summary to w.
func (p *Profile) Write(w io.Writer) {
	fmt.Fprintf(w, "Profile (wall: %s)\n", time.Since(p.start).Round(time.Millisecond))
	fmt.Fprintf(w, "%-30s %12s %8s\n", "Stage", "Duration", "Count")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 54))
	for _, e := range p.Entries() {
		fmt.Fprintf(w, "%-30s %12s %8d\n", e.Stage, e.Duration.Round(time.Microsecond), e.Count)
	}
	fmt.Fprintf(w, "%s\n", repeatChar('-', 54))
	fmt.Fprintf(w, "%-30s %12s\n", "TOTAL", p.Total().Round(time.Microsecond))
}

func repeatChar(c rune, n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = c
	}
	return string(out)
}
