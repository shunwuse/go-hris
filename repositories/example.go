package repositories

import "github.com/shunwuse/go-hris/lib"

type ExampleRepository struct {
	logger lib.Logger
	lib.Database
}

func NewExampleRepository(
	logger lib.Logger,
	db lib.Database,
) ExampleRepository {
	return ExampleRepository{
		logger:   logger,
		Database: db,
	}
}

func (r ExampleRepository) Ping() string {
	r.logger.Info("Ping repository")

	return "pong"
}
