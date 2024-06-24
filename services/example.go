package services

import "github.com/shunwuse/go-hris/lib"

type ExampleService struct {
	logger lib.Logger
}

func NewExampleService() ExampleService {
	logger := lib.GetLogger()

	return ExampleService{
		logger: logger,
	}
}

func (s ExampleService) Ping() string {
	s.logger.Info("Ping service")

	return "pong"
}
