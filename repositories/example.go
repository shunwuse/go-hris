package repositories

import "github.com/shunwuse/go-hris/lib"

type ExampleRepository struct {
	logger lib.Logger
}

func NewExampleRepository() ExampleRepository {
	logger := lib.GetLogger()

	return ExampleRepository{
		logger: logger,
	}
}

func (r ExampleRepository) Ping() string {
	r.logger.Info("Ping repository")

	return "pong"
}
