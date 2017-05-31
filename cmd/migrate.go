package cmd

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/migrations"
	"github.com/RichardKnop/pinglist-api/notifications"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/jinzhu/gorm"
)

var migrationFunctions = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
	metrics.MigrateAll,
	subscriptions.MigrateAll,
	alarms.MigrateAll,
	teams.MigrateAll,
	notifications.MigrateAll,
}

// Migrate runs database migrations
func Migrate() error {
	_, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()

	migrations.MigrateAll(db, migrationFunctions)

	return nil
}
