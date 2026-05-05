// Package watcher provides file-based snapshot polling for envoy-diff.
//
// It watches two snapshot JSON files and emits a diff Event each time
// either file's content changes. Polling is driven by a configurable
// interval and respects context cancellation for clean shutdown.
//
// Basic usage:
//
//	w := watcher.New("base.json", "target.json", 5*time.Second)
//	for ev := range w.Watch(ctx) {
//		if ev.Err != nil {
//			log.Println("watch error:", ev.Err)
//			continue
//		}
//		// process ev.Results
//	}
package watcher
