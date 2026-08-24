package analytics

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestStartIdempotent(t *testing.T) {
	before := runtime.NumGoroutine()
	c := NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// First Start launches the flush loop.
	c.Start(ctx)

	// Allow the goroutine to start and perform its immediate flush.
	time.Sleep(50 * time.Millisecond)

	// Repeated Starts must not spawn additional loops. We can't observe
	// goroutine identity directly, so compare the count before/after.
	for i := 0; i < 10; i++ {
		c.Start(ctx)
	}
	time.Sleep(50 * time.Millisecond)

	runtime.GC()
	// Give any leaked goroutines time to show up in the count.
	time.Sleep(100 * time.Millisecond)
	after := runtime.NumGoroutine()

	// The first Start creates one goroutine. If idempotency is broken,
	// 11 additional goroutines would exist here.
	if after-before > 1 {
		t.Fatalf("expected at most 1 extra goroutine after repeated Start calls, got %d", after-before)
	}

	cancel()
	c.Stop()
}
