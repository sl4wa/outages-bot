package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunNotifierLoop_FirstRunImmediate(t *testing.T) {
	started := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())

	runFn := func(ctx context.Context) error {
		select {
		case started <- struct{}{}:
		default:
		}
		return nil
	}

	done := make(chan struct{})
	go func() {
		runNotifierLoop(ctx, runFn, time.Hour, log.Default())
		close(done)
	}()

	select {
	case <-started:
		// First run happened immediately.
	case <-time.After(time.Second):
		t.Fatal("first run did not start immediately")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("loop did not exit after context cancellation")
	}
}

func TestRunNotifierLoop_ContinuesOnError(t *testing.T) {
	var count atomic.Int32
	ready := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runFn := func(ctx context.Context) error {
		n := count.Add(1)
		if n >= 2 {
			select {
			case ready <- struct{}{}:
			default:
			}
		}
		return fmt.Errorf("test error %d", n)
	}

	done := make(chan struct{})
	go func() {
		runNotifierLoop(ctx, runFn, time.Millisecond, log.Default())
		close(done)
	}()

	select {
	case <-ready:
		// Loop continued past the first error.
	case <-time.After(time.Second):
		t.Fatal("loop did not continue after error")
	}

	if c := count.Load(); c < 2 {
		t.Fatalf("expected at least 2 runs, got %d", c)
	}

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("loop did not exit after context cancellation")
	}
}

func TestRunNotifierLoop_ExitsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	runFn := func(ctx context.Context) error {
		cancel()
		return nil
	}

	done := make(chan struct{})
	go func() {
		runNotifierLoop(ctx, runFn, time.Hour, log.Default())
		close(done)
	}()

	select {
	case <-done:
		// Loop exited.
	case <-time.After(time.Second):
		t.Fatal("loop did not exit after context cancellation")
	}
}

func TestRunNotifierLoop_RejectsNonPositiveInterval(t *testing.T) {
	ctx := context.Background()
	runFn := func(ctx context.Context) error { return nil }

	err := runNotifierLoop(ctx, runFn, 0, log.Default())
	if err == nil {
		t.Fatal("expected error for zero interval")
	}

	err = runNotifierLoop(ctx, runFn, -time.Second, log.Default())
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}
