package profiler

import "time"

// Timer is a convenience helper that measures elapsed time for a single stage
// and records the result into a Profile when stopped.
type Timer struct {
	profile *Profile
	stage   string
	start   time.Time
}

// Start begins timing a named stage and returns a Timer.
// Call Stop (or defer it) to record the measurement.
func (p *Profile) Start(stage string) *Timer {
	return &Timer{profile: p, stage: stage, start: time.Now()}
}

// Stop records the elapsed time since Start was called.
// count is the number of items processed during the stage.
func (t *Timer) Stop(count int) time.Duration {
	d := time.Since(t.start)
	t.profile.Record(t.stage, d, count)
	return d
}

// Elapsed returns the duration since the timer was started without stopping it.
func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}
