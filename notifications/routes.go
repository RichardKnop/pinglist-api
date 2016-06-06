package notifications

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/routes"
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the agencies service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the agencies service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "register_device",
			Method:      "POST",
			Pattern:     "/devices",
			HandlerFunc: service.registerDeviceHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
	}
}
