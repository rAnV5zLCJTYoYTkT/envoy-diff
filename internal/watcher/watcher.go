// Package watcher polls snapshot files for changes and emits diffs automatically.
package watcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/envoy-diff/internal/differ"
	"github.com/envoy-diff/internal/loader"
)

// Event holds the diff results produced when a change is detected.
type Event struct {
	Results []differ.Result
	Err     error
}

// Watcher watches two snapshot files and emits Events when either changes.
type Watcher struct {
	baseFile   string
	targetFile string
	interval   time.Duration
	prevHash   [2][32]byte
}

// New creates a Watcher that polls the given files every interval.
func New(baseFile, targetFile string, interval time.Duration) *Watcher {
	return &Watcher{
		baseFile:   baseFile,
		targetFile: targetFile,
		interval:   interval,
	}
}

// Watch starts polling and sends Events on the returned channel until ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan Event {
	ch := make(chan Event, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.poll(ctx, ch)
			}
		}
	}()
	return ch
}

func (w *Watcher) poll(ctx context.Context, ch chan<- Event) {
	h0, err := hashFile(w.baseFile)
	if err != nil {
		select {
		case ch <- Event{Err: fmt.Errorf("hashing base: %w", err)}:
		case <-ctx.Done():
		}
		return
	}
	h1, err := hashFile(w.targetFile)
	if err != nil {
		select {
		case ch <- Event{Err: fmt.Errorf("hashing target: %w", err)}:
		case <-ctx.Done():
		}
		return
	}
	if h0 == w.prevHash[0] && h1 == w.prevHash[1] {
		return
	}
	w.prevHash[0] = h0
	w.prevHash[1] = h1
	base, err := loader.FromFile(w.baseFile)
	if err != nil {
		select {
		case ch <- Event{Err: fmt.Errorf("loading base: %w", err)}:
		case <-ctx.Done():
		}
		return
	}
	target, err := loader.FromFile(w.targetFile)
	if err != nil {
		select {
		case ch <- Event{Err: fmt.Errorf("loading target: %w", err)}:
		case <-ctx.Done():
		}
		return
	}
	results := differ.Compare(base, target)
	select {
	case ch <- Event{Results: results}:
	case <-ctx.Done():
	}
}

func hashFile(path string) ([32]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return [32]byte{}, err
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out, nil
}
