package infra

import (
	"context"
	"errors"
	"testing"
)

func TestLifecycle_Start(t *testing.T) {
	lc := NewLifecycle()
	startOrder := []int{}

	lc.Append(Hook{
		OnStart: func(ctx context.Context) error {
			startOrder = append(startOrder, 1)
			return nil
		},
	})
	lc.Append(Hook{
		OnStart: func(ctx context.Context) error {
			startOrder = append(startOrder, 2)
			return nil
		},
	})

	err := lc.Start(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(startOrder) != 2 || startOrder[0] != 1 || startOrder[1] != 2 {
		t.Errorf("incorrect start order, expected [1, 2], got %v", startOrder)
	}
}

func TestLifecycle_StopLIFO(t *testing.T) {
	lc := NewLifecycle()
	stopOrder := []int{}

	lc.Append(Hook{
		OnStop: func(ctx context.Context) error {
			stopOrder = append(stopOrder, 1)
			return nil
		},
	})
	lc.Append(Hook{
		OnStop: func(ctx context.Context) error {
			stopOrder = append(stopOrder, 2)
			return nil
		},
	})

	err := lc.Stop(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stopOrder) != 2 || stopOrder[0] != 2 || stopOrder[1] != 1 {
		t.Errorf("incorrect stop order (LIFO), expected [2, 1], got %v", stopOrder)
	}
}

func TestLifecycle_StartError(t *testing.T) {
	lc := NewLifecycle()
	lc.Append(Hook{
		OnStart: func(ctx context.Context) error {
			return errors.New("start failed")
		},
	})

	err := lc.Start(context.Background())
	if err == nil || err.Error() != "start failed" {
		t.Errorf("expected 'start failed' error, got %v", err)
	}
}

func TestLifecycle_StopMultiError(t *testing.T) {
	lc := NewLifecycle()
	lc.Append(Hook{
		OnStop: func(ctx context.Context) error {
			return errors.New("stop 1 failed")
		},
	})
	lc.Append(Hook{
		OnStop: func(ctx context.Context) error {
			return errors.New("stop 2 failed")
		},
	})

	err := lc.Stop(context.Background())
	if err == nil {
		t.Fatal("expected multi-error, got nil")
	}

	errMsg := err.Error()
	if !hasError(errMsg, "stop 1 failed") || !hasError(errMsg, "stop 2 failed") {
		t.Errorf("combined error should contains both failures, got: %s", errMsg)
	}
}

func hasError(msg, target string) bool {
	return errors.New(msg).Error() == target || (errors.Unwrap(errors.New(msg)) != nil && hasError(errors.Unwrap(errors.New(msg)).Error(), target)) || (len(msg) >= len(target) && (msg[:len(target)] == target || msg[len(msg)-len(target):] == target))
}

// Simplified error check for Join
func TestLifecycle_ErrorJoin(t *testing.T) {
	lc := NewLifecycle()
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	lc.Append(Hook{OnStop: func(ctx context.Context) error { return err1 }})
	lc.Append(Hook{OnStop: func(ctx context.Context) error { return err2 }})

	err := lc.Stop(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	// errors.Join preserves the individual errors
	if !errors.Is(err, err1) {
		t.Error("error 1 not found in joined error")
	}
	if !errors.Is(err, err2) {
		t.Error("error 2 not found in joined error")
	}
}
