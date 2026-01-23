package lifecycle

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLifecycle_Start(t *testing.T) {
	lc := New()
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
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2}, startOrder)
}

func TestLifecycle_StopLIFO(t *testing.T) {
	lc := New()
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
	require.NoError(t, err)

	assert.Equal(t, []int{2, 1}, stopOrder)
}

func TestLifecycle_StartError(t *testing.T) {
	lc := New()
	lc.Append(Hook{
		OnStart: func(ctx context.Context) error {
			return errors.New("start failed")
		},
	})

	err := lc.Start(context.Background())
	assert.EqualError(t, err, "start failed")
}

func TestLifecycle_StopMultiError(t *testing.T) {
	lc := New()
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
	require.Error(t, err)

	errMsg := err.Error()
	assert.Contains(t, errMsg, "stop 1 failed")
	assert.Contains(t, errMsg, "stop 2 failed")
}

// Simple test for ErrorJoin consistency
func TestLifecycle_ErrorJoin(t *testing.T) {
	lc := New()
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	lc.Append(Hook{OnStop: func(ctx context.Context) error { return err1 }})
	lc.Append(Hook{OnStop: func(ctx context.Context) error { return err2 }})

	err := lc.Stop(context.Background())
	require.Error(t, err)

	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
}
