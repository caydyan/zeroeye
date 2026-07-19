package analytics_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/yourproject/market/analytics"
)

func TestCollectorStartIdempotency(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collector := &analytics.Collector{}

	// First start
	collector.Start(ctx)
	if!collector.started {
		t.Errorf("Collector should be started after the first call to Start")
	}

	// Second start
	collector.Start(ctx)
	if!collector.started {
		t.Errorf("Collector should remain started after the second call to Start")
	}

	// Third start
	collector.Start(ctx)
	if!collector.started {
		t.Errorf("Collector should remain started after the third call to Start")
	}
}

func TestCollectorStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collector := &analytics.Collector{}
	collector.Start(ctx)

	// Stop the collector
	collector.Stop()
	if collector.started {
		t.Errorf("Collector should not be started after calling Stop")
	}

	// Try to start again
	collector.Start(ctx)
	if!collector.started {
		t.Errorf("Collector should be started after calling Start again")
	}
}

func TestCollectorContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	collector := &analytics.Collector{}
	collector.Start(ctx)

	// Cancel the context
	cancel()

	// Wait for the goroutine to stop
	time.Sleep(100 * time.Millisecond)

	if collector.started {
		t.Errorf("Collector should not be started after context cancellation")
	}
}
