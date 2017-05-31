package subscriptions

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

// Migrate0001 creates files schema
func migrate0001(db *gorm.DB) error {
	migrationName := "subscriptions_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.INFO.Printf("Skipping %s migration", migrationName)
		return nil
	}

	logger.INFO.Printf("Running %s migration", migrationName)

	var err error

	// Create subscription_stripe_event_logs
	if err := db.CreateTable(new(StripeEventLog)).Error; err != nil {
		return fmt.Errorf("Error creating subscription_stripe_event_logs table: %s", err)
	}

	// Create subscription_plans table
	if err := db.CreateTable(new(Plan)).Error; err != nil {
		return fmt.Errorf("Error creating subscription_plans table: %s", err)
	}

	// Create subscription_customers table
	if err := db.CreateTable(new(Customer)).Error; err != nil {
		return fmt.Errorf("Error creating subscription_customers table: %s", err)
	}

	// Create subscription_cards table
	if err := db.CreateTable(new(Card)).Error; err != nil {
		return fmt.Errorf("Error creating subscription_cards table: %s", err)
	}

	// Create subscription_subscriptions table
	if err := db.CreateTable(new(Subscription)).Error; err != nil {
		return fmt.Errorf("Error creating subscription_subscriptions table: %s", err)
	}

	// Add foreign key on subscription_customers.user_id
	err = db.Model(new(Customer)).AddForeignKey(
		"user_id",
		"account_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"subscription_customers.user_id for account_users(id): %s", err)
	}

	// Add foreign key on subscription_cards.customer_id
	err = db.Model(new(Card)).AddForeignKey(
		"customer_id",
		"subscription_customers(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"subscription_cards.customer_id for subscription_customers(id): %s", err)
	}

	// Add foreign key on subscription_subscriptions.customer_id
	err = db.Model(new(Subscription)).AddForeignKey(
		"customer_id",
		"subscription_customers(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"subscription_subscriptions.customer_id for subscription_customers(id): %s", err)
	}

	// Add foreign key on subscription_subscriptions.plan_id
	err = db.Model(new(Subscription)).AddForeignKey(
		"plan_id",
		"subscription_plans(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"subscription_subscriptions.plan_id for subscription_plans(id): %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
