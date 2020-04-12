package pq

import (
	"database/sql"
	"fmt"

	"github.com/houstonj1/go-postgres/config"
	_ "github.com/lib/pq" // postgres driver
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
	logger.Debug("connected to postgres")
}
