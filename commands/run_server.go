package commands

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/facebook"
	"github.com/RichardKnop/pinglist-api/health"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/notifications"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"
)

// RunServer runs the app
func RunServer() error {
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := initServices(cnf, db); err != nil {
		return err
	}

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
	notifications.RegisterRoutes(router, notificationsService)
	alarms.RegisterRoutes(router, alarmsService)

	// Set the router
	app.UseHandler(router)

	return app, nil
}
