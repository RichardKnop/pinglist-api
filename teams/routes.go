package teams

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/routes"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the accounts service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the accounts service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "create_team",
			Method:      "POST",
			Pattern:     "/teams",
			HandlerFunc: service.createTeamHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "get_own_team",
			Method:      "GET",
			Pattern:     "/ownteam",
			HandlerFunc: service.getOwnTeamHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "update_team",
			Method:      "PUT",
			Pattern:     "/teams/{id:[0-9]+}",
			HandlerFunc: service.updateTeamHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
	}
}
