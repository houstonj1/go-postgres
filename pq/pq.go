package pq

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/houstonj1/go-postgres/config"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// Item shape
type Item struct {
	ID          string
	Name        string
	Description string
}

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
	insert(db, logger)
	items := selectAll(db, logger)
	for _, item := range items {
		selectByID(db, logger, item.ID)
		if err != nil {
			logger.Errorf("error selecting item %s: %s", fmt.Errorf("%w", err))
		}
	}
}

func create(db *sql.DB, logger *zap.SugaredLogger) {
	_, err := db.Query(`
		CREATE TABLE public.item (
			id text NOT NULL PRIMARY KEY,
			name text NOT NULL UNIQUE,
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

func insert(db *sql.DB, logger *zap.SugaredLogger) {
	items := []Item{
		{
			ID:          uuid.New().String(),
			Name:        "keyboard",
			Description: "A computer keyboard",
		},
		{
			ID:          uuid.New().String(),
			Name:        "monitor",
			Description: "A computer monitor",
		},
		{
			ID:          uuid.New().String(),
			Name:        "mouse",
			Description: "A computer mouse",
		},
	}
	for _, item := range items {
		query := fmt.Sprintf("INSERT INTO item VALUES ('%s', '%s', '%s')", item.ID, item.Name, item.Description)
		logger.Infof("executing query: %s", query)
		_, err := db.Query(query)
		if err != nil {
			if err.(*pq.Error).Code.Name() == "unique_violation" {
				logger.Info("item already exists")
			} else {
				logger.Fatalf("%s", fmt.Errorf("error executing query: %w", err))
			}
		} else {
			logger.Info("created item")
		}
	}
}

func selectAll(db *sql.DB, logger *zap.SugaredLogger) []Item {
	query := fmt.Sprintf("SELECT id, name, description FROM item")
	logger.Info("selecting all items")
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error selecting all items: %w", err))
	}
	defer rows.Close()
	var items []Item
	var item Item
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Name, &item.Description)
		if err != nil {
			logger.Fatalf("%s", fmt.Errorf("error scanning next row: %w", err))
		}
		logger.Infof("%v", item)
		items = append(items, item)
	}
	return items
}

func selectByID(db *sql.DB, logger *zap.SugaredLogger, id string) Item {
	query := fmt.Sprintf("SELECT id, name, description FROM item WHERE id = '%s'", id)
	logger.Infof("selecting item by ID: %s", id)
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error selecting item with id %s: %w", id, err))
	}
	var item Item
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Name, &item.Description)
		if err != nil {
			logger.Fatalf("%s", fmt.Errorf("error scanning next row: %w", err))
		}
		logger.Infof("%v", item)
	}
	return item
}
