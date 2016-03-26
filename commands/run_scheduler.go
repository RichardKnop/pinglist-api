package commands

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/scheduler"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
)

// RunScheduler runs the scheduler
func RunScheduler() error {
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialise the oauth service
	oauthService := oauth.NewService(cnf, db)

	// Initialise the email service
	emailService := email.NewService(cnf)

	// Initialise the accounts service
	accountsService := accounts.NewService(
		cnf,
		db,
		oauthService,
		emailService,
		nil, // accounts.EmailFactory
	)

	// Initialise the subscriptions service
	subscriptionsService := subscriptions.NewService(
		cnf,
		db,
		accountsService,
		nil, // subscriptions.StripeAdapter
	)

	// Initialise the teams service
	teamsService := teams.NewService(
		cnf,
		db,
		accountsService,
		subscriptionsService,
	)

	// Initialise the metrics service
	metricsService := metrics.NewService(
		cnf,
		db,
		accountsService,
	)

	// Initialise the alarms service
	alarmsService := alarms.NewService(
		cnf,
		db,
		accountsService,
		subscriptionsService,
		teamsService,
		metricsService,
		emailService,
		nil, // alarms.EmailFactory
		nil, // HTTP client
	)

	// Initialise the scheduler
	theScheduler := scheduler.New(metricsService, alarmsService)

	// Run the scheduler
	theScheduler.Run(
		time.Duration(10),  // alarms interval = 10s
		time.Duration(600), // partition / rotate interval  = 10m
	)

	return nil
}
