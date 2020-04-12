package pq

import (
	"database/sql"
	"fmt"

	"github.com/houstonj1/go-postgres/config"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// Pq postgres with lib/pq
func Pq(logger *zap.SugaredLogger) {
	config := config.NewConfig()
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", config.DBUsername, config.DBUsername)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error opening db connection: %w", err))
	}
	if err = db.Ping(); err != nil {
		logger.Fatalf("%s", fmt.Errorf("error connecting to postgres: %w", err))
	}
	logger.Info("connected to postgres")
	create(db, logger)
}

func create(db *sql.DB, logger *zap.SugaredLogger) {
	_, err := db.Query(`
		CREATE TABLE public.item (
			id text NOT NULL PRIMARY KEY,
			name text NOT NULL,
			description text
		);
	`)
	if err != nil {
		if err.(*pq.Error).Code.Name() == "duplicate_table" {
			logger.Info("item table already exists")
			return
		}
		logger.Fatalf("%s", fmt.Errorf("error creating item table: %w", err))
	}
	logger.Info("item table created")
}
