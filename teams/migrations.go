package teams

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

// Migrate0001 creates accounts schema
func migrate0001(db *gorm.DB) error {
	migrationName := "teams_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.Infof("Skipping %s migration", migrationName)
		return nil
	}

	logger.Infof("Running %s migration", migrationName)

	var err error
	// Create team_teams table
	if err := db.CreateTable(new(Team)).Error; err != nil {
		return fmt.Errorf("Error creating team_teams table: %s", err)
	}

	// Add foreign key on team_teams.owner_id
	err = db.Model(new(Team)).AddForeignKey(
		"owner_id",
		"account_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_teams.owner_id for account_users(id): %s", err)
	}

	// Add foreign key on team_team_members.team_id
	err = db.Table("team_team_members").AddForeignKey(
		"team_id",
		"team_teams(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"team_team_members.team_id for team_teams(id): %s",
			err,
		)
	}

	// Add foreign key on team_team_members.user_id
	err = db.Table("team_team_members").AddForeignKey(
		"user_id",
		"account_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"team_team_members.user_id for account_users(id): %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
