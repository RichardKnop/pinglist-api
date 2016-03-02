package commands

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/health"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/web"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
)

// RunServer runs the app
func RunServer() error {
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialise the health service
	healthService := health.NewService(db)

	// Initialise the oauth service
	oauthService := oauth.NewService(cnf, db)

	// Initialise the accounts service
	accountsService := accounts.NewService(cnf, db, oauthService)

	// Initialise the subscriptions service
	subscriptionsService := subscriptions.NewService(cnf, db, accountsService)

	// Initialise the alarms service
	alarmsService := alarms.NewService(
		cnf,
		db,
		accountsService,
		subscriptionsService,
		nil, // HTTP client
	)

	// Initialise the web service
	webService := web.NewService(cnf, accountsService)

	// Start a negroni app
	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(negroni.NewStatic(http.Dir("public")))

	// Create a router instance
	router := mux.NewRouter()

	// Add routes for the health service (healthcheck endpoint)
	health.RegisterRoutes(router, healthService)

	// Add routes for the oauth service (tokens endpoint)
	oauth.RegisterRoutes(router, oauthService)

	// Register routes for the accounts service
	accounts.RegisterRoutes(router, accountsService)

	// Register routes for the alarms service
	alarms.RegisterRoutes(router, alarmsService)

	// Register routes for the subscriptions service
	subscriptions.RegisterRoutes(router, subscriptionsService)

	// Register routes for the web service
	web.RegisterRoutes(router, webService)

	// Set the router
	app.UseHandler(router)

	// Run the server on port 8080
	app.Run(":8080")

	return nil
}
