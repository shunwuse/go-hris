package services

import (
	"context"

	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/ports/service"
	"github.com/shunwuse/go-hris/repositories"
)

type exampleService struct {
	logger            lib.Logger
	exampleRepository repositories.ExampleRepository
}

func NewExampleService(
	logger lib.Logger,
	exampleRepository repositories.ExampleRepository,
) service.ExampleService {
	return exampleService{
		logger:            logger,
		exampleRepository: exampleRepository,
	}
}

func (s exampleService) Ping(ctx context.Context) string {
	s.logger.Info("Ping service")

	pong := s.exampleRepository.Ping(ctx)

	return pong
}
