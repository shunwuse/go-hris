package lib

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var globalDatabase *Database

func GetDatabase() Database {
	if globalDatabase == nil {
		db := newDatabase(NewEnv(), GetLogger())
		globalDatabase = &db
	}

	return *globalDatabase
}

func newDatabase(env Env, logger Logger) Database {
	db, err := gorm.Open(sqlite.Open(env.Sqlite.Database), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database, %v", err)
	}

	logger.Info("Database connected")

	return Database{
		DB: db,
	}
}
