package alerter

import (
	"context"
	"sync"
	"time"

	logalerter "github.com/shunwuse/go-hris/internal/infra/alerter/provider/logger"
	termalerter "github.com/shunwuse/go-hris/internal/infra/alerter/provider/terminal"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/infra"
)

type multiAlerter struct {
	alerters []infra.Alerter
	timeout  time.Duration
}

// NewMultiAlerter creates an alerter that dispatches to multiple alerters in parallel.
func NewMultiAlerter(log *logger.Logger) infra.Alerter {
	return &multiAlerter{
		alerters: []infra.Alerter{
			termalerter.New(),
			logalerter.New(log),
		},
		timeout: 10 * time.Second,
	}
}

func (m *multiAlerter) Send(ctx context.Context, msg infra.Message) error {
	if len(m.alerters) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, alerter := range m.alerters {
		wg.Go(func() {
			_ = alerter.Send(ctx, msg)
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
