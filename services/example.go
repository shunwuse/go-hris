package services

import (
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/repositories"
)

type ExampleService struct {
	logger            lib.Logger
	exampleRepository repositories.ExampleRepository
}

func NewExampleService() ExampleService {
	logger := lib.GetLogger()

	// Initialize repositories
	exampleRepository := repositories.NewExampleRepository()

	return ExampleService{
		logger:            logger,
		exampleRepository: exampleRepository,
	}
}

func (s ExampleService) Ping() string {
	s.logger.Info("Ping service")

	pong := s.exampleRepository.Ping()

	return pong
}
