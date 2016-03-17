package alarms

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
	migrationName := "alarms_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.Infof("Skipping %s migration", migrationName)
		return nil
	}

	logger.Infof("Running %s migration", migrationName)

	var err error

	// Create alarm_states table
	if err := db.CreateTable(new(AlarmState)).Error; err != nil {
		return fmt.Errorf("Error creating alarm_states table: %s", err)
	}

	// Create alarm_alarms table
	if err := db.CreateTable(new(Alarm)).Error; err != nil {
		return fmt.Errorf("Error creating alarm_alarms table: %s", err)
	}

	// Create alarm_incident_types table
	if err := db.CreateTable(new(IncidentType)).Error; err != nil {
		return fmt.Errorf("Error creating alarm_incident_types table: %s", err)
	}

	// Create alarm_incidents table
	if err := db.CreateTable(new(Incident)).Error; err != nil {
		return fmt.Errorf("Error creating alarm_incidents table: %s", err)
	}

	// Create alarm_result_sub_tables table
	if err := db.CreateTable(new(ResultSubTable)).Error; err != nil {
		return fmt.Errorf("Error creating alarm_result_sub_tables table: %s", err)
	}

	// Create alarm_results table
	if err := db.CreateTable(new(Result)).Error; err != nil {
		return fmt.Errorf("Error creating alarm_results table: %s", err)
	}

	// Add foreign key on alarm_alarms.user_id
	err = db.Model(new(Alarm)).AddForeignKey(
		"user_id",
		"account_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"alarm_alarms.user_id for account_users(id): %s", err)
	}

	// Add foreign key on alarm_alarms.alarm_state_id
	err = db.Model(new(Alarm)).AddForeignKey(
		"alarm_state_id",
		"alarm_states(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"alarm_alarms.alarm_state_id for alarm_states(id): %s", err)
	}

	// Add foreign key on alarm_incidents.alarm_id
	err = db.Model(new(Incident)).AddForeignKey(
		"alarm_id",
		"alarm_alarms(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"alarm_incidents.alarm_id for alarm_alarms(id): %s", err)
	}

	// Add foreign key on alarm_incidents.incident_type_id
	err = db.Model(new(Incident)).AddForeignKey(
		"incident_type_id",
		"alarm_incident_types(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"alarm_incidents.incident_type_id for alarm_incident_types(id): %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
