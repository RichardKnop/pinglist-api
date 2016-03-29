package commands

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/RichardKnop/pinglist-api/facebook"
	"github.com/RichardKnop/pinglist-api/health"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/notifications"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/RichardKnop/pinglist-api/web"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/phyber/negroni-gzip/gzip"
)

// RunServer runs the app
func RunServer() error {
	// Init config and database
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()

	// Init the app
	app, err := initApp(cnf, db)
	if err != nil {
		return err
	}

	// Run the server on port 8080
	app.Run(":8080")

	return nil
}

// initApp starts all services, creates a negroni app, registers all routes
func initApp(cnf *config.Config, db *gorm.DB) (*negroni.Negroni, error) {
	// Initialise services
	healthService := health.NewService(db)
	oauthService := oauth.NewService(cnf, db)
	emailService := email.NewService(cnf)
	accountsService := accounts.NewService(
		cnf,
		db,
		oauthService,
		emailService,
		nil, // accounts.EmailFactory
	)
	facebookService := facebook.NewService(
		cnf,
		db,
		accountsService,
		nil, // facebook.Adapter
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
	notificationsService := notifications.NewService(
		cnf,
		db,
		accountsService,
		nil, // notifications.SNSAdapter
	)
	webService := web.NewService(cnf, accountsService)

	// Start a negroni app
	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(negroni.NewStatic(http.Dir("public")))

	// Create a router instance
	router := mux.NewRouter()

	// Register routes
	health.RegisterRoutes(router, healthService)
	oauth.RegisterRoutes(router, oauthService)
	accounts.RegisterRoutes(router, accountsService)
	facebook.RegisterRoutes(router, facebookService)
	subscriptions.RegisterRoutes(router, subscriptionsService)
	teams.RegisterRoutes(router, teamsService)
	metrics.RegisterRoutes(router, metricsService)
	alarms.RegisterRoutes(router, alarmsService)
	notifications.RegisterRoutes(router, notificationsService)
	web.RegisterRoutes(router, webService)

	// Set the router
	app.UseHandler(router)

	return app, nil
}
