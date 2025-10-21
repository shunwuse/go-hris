package infra

import (
	"log/slog"
	"os"

	ent "github.com/shunwuse/go-hris/ent/entgen"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Client *ent.Client
}

var globalDatabase *Database

func GetDatabase() Database {
	if globalDatabase == nil {
		db := newDatabase(NewEnv())
		globalDatabase = &db
	}

	return *globalDatabase
}

func newDatabase(env Env) Database {
	client, err := ent.Open("sqlite3", env.Sqlite.Database)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database connected")

	return Database{
		Client: client,
	}
}
