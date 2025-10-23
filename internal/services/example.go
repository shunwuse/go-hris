package services

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
)

type exampleService struct {
	logger            *infra.Logger
	exampleRepository *repositories.ExampleRepository
}

func NewExampleService(
	logger *infra.Logger,
	exampleRepository *repositories.ExampleRepository,
) service.ExampleService {
	return &exampleService{
		logger:            logger,
		exampleRepository: exampleRepository,
	}
}

func (s *exampleService) Ping(ctx context.Context) string {
	s.logger.WithContext(ctx).Info("ping service invoked")

	pong := s.exampleRepository.Ping(ctx)

	return pong
}
