package profiler

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew_InitialisesStartTime(t *testing.T) {
	before := time.Now()
	p := New()
	after := time.Now()
	if p.start.Before(before) || p.start.After(after) {
		t.Errorf("start time %v not in expected range [%v, %v]", p.start, before, after)
	}
}

func TestRecord_AppendedAndTotal(t *testing.T) {
	p := New()
	p.Record("stage-a", 10*time.Millisecond, 5)
	p.Record("stage-b", 20*time.Millisecond, 3)

	if got := p.Total(); got != 30*time.Millisecond {
		t.Errorf("Total() = %v, want 30ms", got)
	}
	if len(p.entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(p.entries))
	}
}

func TestEntries_SortedByDurationDescending(t *testing.T) {
	p := New()
	p.Record("fast", 1*time.Millisecond, 1)
	p.Record("slow", 50*time.Millisecond, 10)
	p.Record("medium", 10*time.Millisecond, 4)

	entries := p.Entries()
	if entries[0].Stage != "slow" {
		t.Errorf("expected slow first, got %s", entries[0].Stage)
	}
	if entries[2].Stage != "fast" {
		t.Errorf("expected fast last, got %s", entries[2].Stage)
	}
}

func TestWrite_ContainsStageNames(t *testing.T) {
	p := New()
	p.Record("filter", 5*time.Millisecond, 8)
	p.Record("patch", 2*time.Millisecond, 8)

	var buf bytes.Buffer
	p.Write(&buf)
	out := buf.String()

	for _, want := range []string{"filter", "patch", "TOTAL", "Profile"} {
		if !strings.Contains(out, want) {
			t.Errorf("Write() output missing %q", want)
		}
	}
}

func TestTimer_StopRecordsEntry(t *testing.T) {
	p := New()
	tmr := p.Start("my-stage")
	time.Sleep(1 * time.Millisecond)
	d := tmr.Stop(42)

	if d < time.Millisecond {
		t.Errorf("expected duration >= 1ms, got %v", d)
	}
	if len(p.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(p.entries))
	}
	if p.entries[0].Stage != "my-stage" {
		t.Errorf("stage = %q, want my-stage", p.entries[0].Stage)
	}
	if p.entries[0].Count != 42 {
		t.Errorf("count = %d, want 42", p.entries[0].Count)
	}
}

func TestTimer_ElapsedDoesNotRecord(t *testing.T) {
	p := New()
	tmr := p.Start("probe")
	_ = tmr.Elapsed()
	if len(p.entries) != 0 {
		t.Errorf("Elapsed() should not record an entry")
	}
}

func TestTotal_Empty(t *testing.T) {
	p := New()
	if p.Total() != 0 {
		t.Errorf("Total() on empty profile = %v, want 0", p.Total())
	}
}
