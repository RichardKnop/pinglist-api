package notifications

import (
	"fmt"

	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/migrations"
	"github.com/jinzhu/gorm"
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	if err := migrate0001(db); err != nil {
		return err
	}

	return nil
}

// Migrate0001 creates notifications schema
func migrate0001(db *gorm.DB) error {
	migrationName := "notifications_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.INFO.Printf("Skipping %s migration", migrationName)
		return nil
	}

	logger.INFO.Printf("Running %s migration", migrationName)

	// Create agency_states table
	if err := db.CreateTable(new(Endpoint)).Error; err != nil {
		return fmt.Errorf("Error creating agency_states table: %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
