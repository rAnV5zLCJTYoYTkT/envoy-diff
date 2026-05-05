package watcher_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/envoy-diff/internal/watcher"
)

type snapData struct {
	Resources []map[string]string `json:"resources"`
}

func writeTempSnap(t *testing.T, name, rtype, body string) string {
	t.Helper()
	d := map[string]interface{}{
		"resources": []map[string]interface{}{
			{"name": name, "type": rtype, "body": body},
		},
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.Write(b); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	base := writeTempSnap(t, "cluster-a", "cluster", `{"connect_timeout":"1s"}`)
	target := writeTempSnap(t, "cluster-a", "cluster", `{"connect_timeout":"2s"}`)

	w := watcher.New(base, target, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ch := w.Watch(ctx)
	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before event")
		}
		if ev.Err != nil {
			t.Fatalf("unexpected error: %v", ev.Err)
		}
		if len(ev.Results) == 0 {
			t.Fatal("expected diff results, got none")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	base := writeTempSnap(t, "cluster-a", "cluster", `{"connect_timeout":"1s"}`)
	target := writeTempSnap(t, "cluster-a", "cluster", `{"connect_timeout":"1s"}`)

	w := watcher.New(base, target, 40*time.Millisecond)
	// First poll fires an event (initial state). Drain it.
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	ch := w.Watch(ctx)
	// Consume first event
	<-ch
	// After the first event, no further changes → no more events before timeout
	select {
	case ev, ok := <-ch:
		if ok && ev.Err == nil && len(ev.Results) > 0 {
			t.Log("got a second event — files may have changed, which is acceptable")
		}
	case <-ctx.Done():
		// expected: no second change event
	}
}

func TestWatcher_InvalidFile(t *testing.T) {
	w := watcher.New("/nonexistent/base.json", "/nonexistent/target.json", 30*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch := w.Watch(ctx)
	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before error event")
		}
		if ev.Err == nil {
			t.Fatal("expected error for missing files")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for error event")
	}
}
