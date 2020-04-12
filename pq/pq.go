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
	logger.Info("------------ CREATE TABLE -----------")
	create(db, logger)
	logger.Info("-------------------------------------")
	logger.Info("-------------- INSERT ---------------")
	insert(db, logger)
	logger.Info("-------------------------------------")
	logger.Info("------------ SELECT ALL -------------")
	items := selectAll(db, logger)
	logger.Info("-------------------------------------")
	logger.Info("----------- SELECT BY ID ------------")
	for _, item := range items {
		selectByID(db, logger, item.ID)
		if err != nil {
			logger.Errorf("error selecting item %s: %s", item.ID, fmt.Errorf("%w", err))
		}
	}
	logger.Info("-------------------------------------")
	logger.Info("-------------- UPSERT ---------------")
	items = append(items, Item{
		ID:          uuid.New().String(),
		Name:        "wireless mouse",
		Description: "A wireless computer mouse",
	})
	items[2].Description = "A wired computer mouse"
	upsert(db, logger, items)
	logger.Info("-------------------------------------")
	logger.Info("----------- DELETE BY ID ------------")
	items = selectAll(db, logger)
	for _, item := range items {
		deleteByID(db, logger, item.ID)
		if err != nil {
			logger.Errorf("error deleting item %s: %s", item.ID, fmt.Errorf("%w", err))
		}
	}
	logger.Info("-------------------------------------")
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

func deleteByID(db *sql.DB, logger *zap.SugaredLogger, id string) {
	query := fmt.Sprintf("DELETE FROM item WHERE id = '%s'", id)
	logger.Infof("deleting item by ID: %s", id)
	_, err := db.Query(query)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error deleting item with id %s: %w", id, err))
	}
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

func upsert(db *sql.DB, logger *zap.SugaredLogger, items []Item) {
	length := len(items) - 1
	values := ""
	for index, item := range items {
		var value string
		if index != length {
			value = fmt.Sprintf("('%s', '%s', '%s'),", item.ID, item.Name, item.Description)
		} else {
			value = fmt.Sprintf("('%s', '%s', '%s')", item.ID, item.Name, item.Description)
		}
		logger.Infof("upserting on id %s: %s", item.ID, value)
		values += value
	}
	query := fmt.Sprintf(`
		INSERT INTO item(id, name, description)
		VALUES %s
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name,
		description = EXCLUDED.description`,
		values,
	)
	_, err := db.Query(query)
	if err != nil {
		logger.Fatalf("%s", fmt.Errorf("error upserting values: %w", err))
	}
}
