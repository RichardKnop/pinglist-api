package commands

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/migrations"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
)

// Migrate runs database migrations
func Migrate() error {
	_, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()

	// Bootstrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := oauth.MigrateAll(db); err != nil {
		return err
	}

	// Run migrations for the accounts service
	if err := accounts.MigrateAll(db); err != nil {
		return err
	}

	// Run migrations for the alarms service
	if err := alarms.MigrateAll(db); err != nil {
		return err
	}

	// Run migrations for the subscriptions service
	if err := subscriptions.MigrateAll(db); err != nil {
		return err
	}

	return nil
}
