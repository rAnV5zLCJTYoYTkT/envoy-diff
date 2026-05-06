// Package pipeline chains envoy-diff processing stages into a single
// composable execution unit.
//
// A Pipeline is constructed from one or more Stage functions, each of which
// accepts and returns a []differ.DiffResult slice. Stages are executed in
// declaration order, allowing callers to mix built-in helpers (FilterStage,
// PatchStage, TagStage, AnnotateStage) with custom logic.
//
// Example:
//
//	p := pipeline.New(
//		pipeline.FilterStage(filter.Options{Types: []string{"cluster"}}),
//		pipeline.TagStage(tagger.New()),
//		pipeline.AnnotateStage(annotator.New()),
//	)
//	results = p.Run(results)
package pipeline
