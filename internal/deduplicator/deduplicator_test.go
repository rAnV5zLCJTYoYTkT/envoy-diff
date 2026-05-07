package deduplicator_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/deduplicator"
	"github.com/your-org/envoy-diff/internal/differ"
)

func makeResults(pairs ...string) []differ.Result {
	if len(pairs)%2 != 0 {
		panic("makeResults: need even number of args (type, name)")
	}
	out := make([]differ.Result, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, differ.Result{Type: pairs[i], Name: pairs[i+1], Status: differ.StatusUnchanged})
	}
	return out
}

func TestApply_NoDuplicates(t *testing.T) {
	results := makeResults("CDS", "cluster-a", "LDS", "listener-b")
	got := deduplicator.Apply(results, deduplicator.DefaultOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_RemovesDuplicates_KeepFirst(t *testing.T) {
	results := makeResults("CDS", "cluster-a", "CDS", "cluster-a", "LDS", "listener-b")
	results[1].Status = differ.StatusModified // second occurrence differs

	opts := deduplicator.DefaultOptions() // KeepFirst = true
	got := deduplicator.Apply(results, opts)

	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].Status != differ.StatusUnchanged {
		t.Errorf("expected first occurrence (Unchanged) to be kept")
	}
}

func TestApply_RemovesDuplicates_KeepLast(t *testing.T) {
	results := makeResults("CDS", "cluster-a", "CDS", "cluster-a")
	results[1].Status = differ.StatusModified

	opts := deduplicator.Options{KeyFn: nil, KeepFirst: false}
	got := deduplicator.Apply(results, opts)

	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Status != differ.StatusModified {
		t.Errorf("expected last occurrence (Modified) to be kept")
	}
}

func TestApply_CustomKeyFn(t *testing.T) {
	results := makeResults("CDS", "cluster-a", "LDS", "cluster-a")
	// deduplicate by name only
	opts := deduplicator.Options{
		KeyFn:     func(r differ.Result) string { return r.Name },
		KeepFirst: true,
	}
	got := deduplicator.Apply(results, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 result with name-only key, got %d", len(got))
	}
}

func TestCount_ReturnsDuplicateCount(t *testing.T) {
	results := makeResults("CDS", "c1", "CDS", "c1", "CDS", "c1", "LDS", "l1")
	n := deduplicator.Count(results, deduplicator.DefaultOptions())
	if n != 2 {
		t.Errorf("expected 2 duplicates, got %d", n)
	}
}

func TestApply_Empty(t *testing.T) {
	got := deduplicator.Apply(nil, deduplicator.DefaultOptions())
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d", len(got))
	}
}
