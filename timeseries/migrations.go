package timeseries

import (
	"fmt"

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

// Migrate0001 creates files schema
func migrate0001(db *gorm.DB) error {
	migrationName := "timeseries_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.Infof("Skipping %s migration", migrationName)
		return nil
	}

	logger.Infof("Running %s migration", migrationName)

	// Create timeseries_sub_tables table
	if err := db.CreateTable(new(SubTable)).Error; err != nil {
		return fmt.Errorf("Error creating timeseries_sub_tables table: %s", err)
	}

	// Create timeseries_request_times table
	if err := db.CreateTable(new(RequestTime)).Error; err != nil {
		return fmt.Errorf("Error creating timeseries_request_times table: %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
