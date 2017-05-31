package cmd

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/facebook"
	"github.com/RichardKnop/pinglist-api/health"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/notifications"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/jinzhu/gorm"
)

var (
	healthService        health.ServiceInterface
	oauthService         oauth.ServiceInterface
	accountsService      accounts.ServiceInterface
	facebookService      facebook.ServiceInterface
	subscriptionsService subscriptions.ServiceInterface
	teamsService         teams.ServiceInterface
	metricsService       metrics.ServiceInterface
	notificationsService notifications.ServiceInterface
	alarmsService        alarms.ServiceInterface
)

// initConfigDB loads the configuration and connects to the database
func initConfigDB(mustLoadOnce, keepReloading bool) (*config.Config, *gorm.DB, error) {
	// Config
	cnf := config.NewConfig(mustLoadOnce, keepReloading)

	// Database
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return nil, nil, err
	}

	return cnf, db, nil
}

// initServices starts up all services and sets above defined variables
func initServices(cnf *config.Config, db *gorm.DB) error {
	// Initialise services
	healthService = health.NewService(db)
	oauthService = oauth.NewService(cnf, db)
	accountsService = accounts.NewService(
		cnf,
		db,
		oauthService,
		nil, // email.Service
		nil, // accounts.EmailFactory
	)
	facebookService = facebook.NewService(
		cnf,
		db,
		accountsService,
		nil, // facebook.Adapter
	)
	subscriptionsService = subscriptions.NewService(
		cnf,
		db,
		accountsService,
		nil, // subscriptions.StripeAdapter
	)
	teamsService = teams.NewService(cnf, db, accountsService, subscriptionsService)
	metricsService = metrics.NewService(cnf, db, accountsService)
	notificationsService = notifications.NewService(
		cnf,
		db,
		accountsService,
		nil, // notifications.SNSAdapter
	)
	alarmsService = alarms.NewService(
		cnf,
		db,
		accountsService,
		subscriptionsService,
		teamsService,
		metricsService,
		notificationsService,
		nil, // email.Service
		nil, // alarms.EmailFactory
		nil, // alarms.SlackFactory
		nil, // HTTP client
	)

	return nil
}
