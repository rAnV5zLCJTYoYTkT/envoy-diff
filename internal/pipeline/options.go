package pipeline

import (
	"github.com/yourorg/envoy-diff/internal/annotator"
	"github.com/yourorg/envoy-diff/internal/filter"
	"github.com/yourorg/envoy-diff/internal/patcher"
	"github.com/yourorg/envoy-diff/internal/tagger"
)

// Options configures a Pipeline created via NewWithOptions.
type Options struct {
	Filter    filter.Options
	Patcher   *patcher.Patcher
	Tagger    *tagger.Tagger
	Annotator *annotator.Annotator
	// SkipFilter disables the filter stage even when Filter is non-zero.
	SkipFilter bool
}

// NewWithOptions builds a Pipeline from Options, falling back to default
// component instances when fields are nil.
func NewWithOptions(opts Options) *Pipeline {
	var stages []Stage

	if !opts.SkipFilter {
		stages = append(stages, FilterStage(opts.Filter))
	}

	pt := opts.Patcher
	if pt == nil {
		pt = patcher.New()
	}
	stages = append(stages, PatchStage(pt))

	tg := opts.Tagger
	if tg == nil {
		tg = tagger.New()
	}
	stages = append(stages, TagStage(tg))

	an := opts.Annotator
	if an == nil {
		an = annotator.New()
	}
	stages = append(stages, AnnotateStage(an))

	return New(stages...)
}
