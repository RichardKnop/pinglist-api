package commands

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/scheduler"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/jinzhu/gorm"
)

// RunScheduler runs the scheduler
func RunScheduler() error {
	// Init config and database
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()

	// Init the scheduler
	theScheduler, err := initScheduler(cnf, db)
	if err != nil {
		return err
	}

	// Run the scheduling goroutines
	wg := theScheduler.Run(
		time.Duration(10),  // alarms check interval = 10s
		time.Duration(600), // partition / rotate interval = 10m
	)

	// The Run method returns sync.WaitGroup, use it to block the return
	wg.Wait()

	return nil
}

// initScheduler starts a scheduler instance
func initScheduler(cnf *config.Config, db *gorm.DB) (*scheduler.Scheduler, error) {
	// Initialise services
	oauthService := oauth.NewService(cnf, db)
	emailService := email.NewService(cnf)
	accountsService := accounts.NewService(
		cnf,
		db,
		oauthService,
		emailService,
		nil, // accounts.EmailFactory
	)
	subscriptionsService := subscriptions.NewService(
		cnf,
		db,
		accountsService,
		nil, // subscriptions.StripeAdapter
	)
	teamsService := teams.NewService(
		cnf,
		db,
		accountsService,
		subscriptionsService,
	)
	metricsService := metrics.NewService(
		cnf,
		db,
		accountsService,
	)
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

	return scheduler.New(metricsService, alarmsService), nil
}
