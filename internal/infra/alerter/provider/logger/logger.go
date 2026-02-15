package logalerter

import (
	"context"

	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/infra"
	"go.uber.org/zap"
)

type sender struct {
	logger *logger.Logger
}

func New(log *logger.Logger) infra.Alerter {
	return &sender{
		logger: log,
	}
}

func (s *sender) Send(ctx context.Context, msg infra.Message) error {
	log := s.logger.WithContext(ctx)

	fields := []zap.Field{
		zap.String("content", msg.Content),
	}

	if msg.TraceID != "" {
		fields = append(fields, zap.String("trace_id", msg.TraceID))
	}

	switch msg.Level {
	case infra.LevelCritical:
		if msg.StackTrace != "" {
			fields = append(fields, zap.String("stack_trace", msg.StackTrace))
		}
		log.Error("CRITICAL ALERT: "+msg.Title, fields...)
	default:
		log.Error("ALERT: "+msg.Title, fields...)
	}

	return nil
}

// func init() {
// alert.Register(alert.ProviderLog, &sender{})
// }
