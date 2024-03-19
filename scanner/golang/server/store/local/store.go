package local

import (
	"fmt"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type handler struct {
	db *gorm.DB
}

func NewStore() (store.Store, error) {
	// Create database
	db, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	// Create and initialize db handler
	handler := &handler{db: db}
	if err := handler.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	return handler, nil
}

func (h *handler) init() error {
	// Add extensions
	//if err := h.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
	//	return fmt.Errorf("failed to create uuid extension: %w", err)
	//}

	// Migrate models
	if err := h.db.AutoMigrate(
		scanModel{},
		findingModel{},
	); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	return nil
}

func getMetaKVSelectors(column string, metaSelectors map[string]string) *datatypes.JSONQueryExpression {
	if len(metaSelectors) == 0 {
		return nil
	}

	// this extracts selectors such as key=value specified for metaSelector query param
	jqe := datatypes.JSONQuery(column)
	for key, value := range metaSelectors {
		key, value := key, value
		jqe = jqe.Equals(value, key)
	}
	return jqe
}
