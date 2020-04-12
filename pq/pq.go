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
	connStr := fmt.Sprintf("user=%s dbname=%s", config.DBUsername, config.DBUsername)
	_, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatalf("error creating db connection: %w", err)
	}
	logger.Debug("connected to postgres")
}
