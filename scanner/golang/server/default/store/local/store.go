package local

import (
	"fmt"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

type baseModel struct {
	ID string `sql:"type:uuid;primary_key"`
}

type handler struct {
	db *gorm.DB
}

func NewStore() (types.Store, error) {
	// Create database
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	// Migrate models
	if err := db.AutoMigrate(
		scanModel{},
		scanResultModel{},
	); err != nil {
		return nil, fmt.Errorf("failed to run auto migration: %w", err)
	}

	return &handler{
		db: db,
	}, nil
}
