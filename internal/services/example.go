package services

import (
	"context"
	"log/slog"

	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
)

type exampleService struct {
	exampleRepository repositories.ExampleRepository
}

func NewExampleService(
	exampleRepository repositories.ExampleRepository,
) service.ExampleService {
	return exampleService{
		exampleRepository: exampleRepository,
	}
}

func (s exampleService) Ping(ctx context.Context) string {
	slog.Info("Ping service")

	pong := s.exampleRepository.Ping(ctx)

	return pong
}
