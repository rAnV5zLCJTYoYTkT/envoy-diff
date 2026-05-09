package pipeline

import (
	"io"

	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/profiler"
)

// ProfiledPipeline wraps a Pipeline and records per-stage timing via a
// profiler.Profile. Use WriteProfile to print results after Run.
type ProfiledPipeline struct {
	pipeline *Pipeline
	profile  *profiler.Profile
}

// NewProfiled creates a ProfiledPipeline that instruments every stage of p.
func NewProfiled(p *Pipeline) *ProfiledPipeline {
	return &ProfiledPipeline{
		pipeline: p,
		profile:  profiler.New(),
	}
}

// Run executes all stages sequentially, recording timing for each.
func (pp *ProfiledPipeline) Run(results []differ.Result) []differ.Result {
	current := results
	for _, s := range pp.pipeline.stages {
		tmr := pp.profile.Start(s.Name())
		current = s.Apply(current)
		tmr.Stop(len(current))
	}
	return current
}

// Profile returns the underlying profiler.Profile for inspection or export.
func (pp *ProfiledPipeline) Profile() *profiler.Profile {
	return pp.profile
}

// WriteProfile writes the collected timing summary to w.
func (pp *ProfiledPipeline) WriteProfile(w io.Writer) {
	pp.profile.Write(w)
}
