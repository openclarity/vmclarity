package local

import (
	"fmt"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func NewStore() (types.Store, error) {
	// Create database
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
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
