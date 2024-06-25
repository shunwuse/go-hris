package lib

import (
	"github.com/shunwuse/go-hris/models"
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

	logger.Info("Migrating database")
	if err = migrate(db); err != nil {
		logger.Fatalf("Error migrating database, %v", err)
	}

	return Database{
		DB: db,
	}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
