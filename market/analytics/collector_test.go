package analytics

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCollector_StartIdempotent(t *testing.T) {
	c := NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ensure not running initially
	c.mu.RLock()
	running := c.running
	c.mu.RUnlock()
	if running {
		t.Fatal("expected collector to not be running initially")
	}

	// First start
	c.Start(ctx)

	c.mu.RLock()
	running = c.running
	c.mu.RUnlock()
	if !running {
		t.Fatal("expected collector to be running after first Start")
	}

	// Second start
	c.Start(ctx)

	c.mu.RLock()
	running = c.running
	c.mu.RUnlock()
	if !running {
		t.Fatal("expected collector to still be running after second Start")
	}

	// Wait for the goroutine to enter the select loop
	time.Sleep(50 * time.Millisecond)

	// Stop the collector
	c.Stop()

	// Wait for goroutine to exit and update state
	// Use a small loop to avoid flaky tests
	stopped := false
	for i := 0; i < 50; i++ {
		c.mu.RLock()
		r := c.running
		c.mu.RUnlock()
		if !r {
			stopped = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !stopped {
		t.Fatal("expected collector to stop after Stop() is called")
	}

	// Context cancellation should also stop it if started again
	ctx2, cancel2 := context.WithCancel(context.Background())
	c.Start(ctx2)

	c.mu.RLock()
	if !c.running {
		c.mu.RUnlock()
		t.Fatal("expected collector to be running again")
	}
	c.mu.RUnlock()

	cancel2()

	stopped = false
	for i := 0; i < 50; i++ {
		c.mu.RLock()
		r := c.running
		c.mu.RUnlock()
		if !r {
			stopped = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !stopped {
		t.Fatal("expected collector to stop after context cancellation")
	}
}

func TestCollector_StartIdempotent_Concurrent(t *testing.T) {
	c := NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Start(ctx)
		}()
	}
	wg.Wait()

	c.mu.RLock()
	running := c.running
	c.mu.RUnlock()
	if !running {
		t.Fatal("expected collector to be running")
	}

	time.Sleep(50 * time.Millisecond)
	c.Stop()
	
	stopped := false
	for i := 0; i < 50; i++ {
		c.mu.RLock()
		r := c.running
		c.mu.RUnlock()
		if !r {
			stopped = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !stopped {
		t.Fatal("expected collector to stop after concurrent start")
	}
}
