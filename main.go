package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/houstonj1/go-postgres/config"
	"github.com/houstonj1/go-postgres/pq"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic("error creating logger")
	}
	logger := zapLogger.Sugar()

	config := config.NewConfig()
	logger.Debugf("%s", config.Print())
	migrateUp(config, logger)

	logger.Info("-------------------------------------")
	logger.Info("-------------  lib/pq  --------------")
	logger.Info("-------------------------------------")
	pq.Pq(logger)

	migrateDown(config, logger)
}

func migrateUp(config *config.Config, logger *zap.SugaredLogger) {
	connStr := fmt.Sprintf("postgres://%s@localhost:5432/%s?sslmode=disable", config.DBUsername, config.DBPassword)
	migrations, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error in migration setup: %w", err))
	}
	err = migrations.Up()
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error migrating database: %w", err))
	}
	logger.Info("migrateUp completed")
}

func migrateDown(config *config.Config, logger *zap.SugaredLogger) {
	connStr := fmt.Sprintf("postgres://%s@localhost:5432/%s?sslmode=disable", config.DBUsername, config.DBPassword)
	migrations, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error in migration setup: %w", err))
	}
	err = migrations.Down()
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error migrating database: %w", err))
	}
	logger.Info("migrateDown completed")
}
