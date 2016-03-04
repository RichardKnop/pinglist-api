package migrations

import (
	"github.com/jinzhu/gorm"
)

// MigrateAllTest runs bootstrap, then all migration functions listed against
// the specified database and logs any errors
func MigrateAllTest(db *gorm.DB, migrationFunctions []func(*gorm.DB) error) {
	if err := Bootstrap(db); err != nil {
		logger.Error(err)
	}

	for _, m := range migrationFunctions {
		if err := m(db); err != nil {
			logger.Error(err)
		}
	}
}
