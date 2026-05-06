// Package pipeline provides a composable processing pipeline for diff results.
// It chains multiple transformation steps (filter, patch, tag, annotate, group)
// into a single reusable execution unit.
package pipeline

import (
	"github.com/yourorg/envoy-diff/internal/annotator"
	"github.com/yourorg/envoy-diff/internal/differ"
	"github.com/yourorg/envoy-diff/internal/filter"
	"github.com/yourorg/envoy-diff/internal/patcher"
	"github.com/yourorg/envoy-diff/internal/tagger"
)

// Stage is a function that transforms a slice of DiffResults.
type Stage func([]differ.DiffResult) []differ.DiffResult

// Pipeline holds an ordered list of processing stages.
type Pipeline struct {
	stages []Stage
}

// New creates a Pipeline with the provided stages applied in order.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run executes all stages sequentially and returns the final results.
func (p *Pipeline) Run(results []differ.DiffResult) []differ.DiffResult {
	out := results
	for _, s := range p.stages {
		out = s(out)
	}
	return out
}

// FilterStage wraps filter.Apply as a pipeline Stage.
func FilterStage(opts filter.Options) Stage {
	return func(results []differ.DiffResult) []differ.DiffResult {
		return filter.Apply(results, opts)
	}
}

// PatchStage wraps a patcher.Patcher as a pipeline Stage.
func PatchStage(pt *patcher.Patcher) Stage {
	return func(results []differ.DiffResult) []differ.DiffResult {
		return pt.Apply(results)
	}
}

// TagStage wraps a tagger.Tagger as a pipeline Stage.
func TagStage(tg *tagger.Tagger) Stage {
	return func(results []differ.DiffResult) []differ.DiffResult {
		return tg.Apply(results)
	}
}

// AnnotateStage wraps an annotator.Annotator as a pipeline Stage.
func AnnotateStage(an *annotator.Annotator) Stage {
	return func(results []differ.DiffResult) []differ.DiffResult {
		return an.Apply(results)
	}
}

// DefaultPipeline builds a pipeline using default rules for patching, tagging,
// and annotation with no active filters.
func DefaultPipeline() *Pipeline {
	return New(
		PatchStage(patcher.New()),
		TagStage(tagger.New()),
		AnnotateStage(annotator.New()),
	)
}
