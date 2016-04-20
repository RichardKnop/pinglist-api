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
			Name:        "get_team",
			Method:      "GET",
			Pattern:     "/teams/{id:[0-9]+}",
			HandlerFunc: service.getTeamHandler,
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
		routes.Route{
			Name:        "delete_team",
			Method:      "DELETE",
			Pattern:     "/teams/{id:[0-9]+}",
			HandlerFunc: service.deleteTeamHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "list_teams",
			Method:      "GET",
			Pattern:     "/teams",
			HandlerFunc: service.listTeamsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "invite_user",
			Method:      "POST",
			Pattern:     "/teams/{id:[0-9]+}/invitations",
			HandlerFunc: service.inviteUserHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
	}
}
