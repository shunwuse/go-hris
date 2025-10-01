package services

import (
	"context"

	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/repositories"
)

type ExampleService struct {
	logger            lib.Logger
	exampleRepository repositories.ExampleRepository
}

func NewExampleService(
	logger lib.Logger,
	exampleRepository repositories.ExampleRepository,
) ExampleService {
	return ExampleService{
		logger:            logger,
		exampleRepository: exampleRepository,
	}
}

func (s ExampleService) Ping(ctx context.Context) string {
	s.logger.Info("Ping service")

	pong := s.exampleRepository.Ping(ctx)

	return pong
}
