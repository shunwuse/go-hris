package infra

import (
	ent "github.com/shunwuse/go-hris/ent/entgen"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Client *ent.Client
}

var globalDatabase *Database

func GetDatabase() *Database {
	if globalDatabase == nil {
		db := newDatabase(GetConfig(), GetLogger())
		globalDatabase = &db
	}

	return globalDatabase
}

func newDatabase(config Config, logger *Logger) Database {
	client, err := ent.Open("sqlite3", config.SqliteDBPath)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	logger.Info("database connected successfully")

	return Database{
		Client: client,
	}
}
