package accounts

import (
	"github.com/RichardKnop/pinglist-api/routes"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the accounts service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1/accounts").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the accounts service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "create_user",
			Method:      "POST",
			Pattern:     "/users",
			HandlerFunc: service.createUserHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "get_my_user",
			Method:      "GET",
			Pattern:     "/me",
			HandlerFunc: service.getMyUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "get_user",
			Method:      "GET",
			Pattern:     "/users/{id:[0-9]+}",
			HandlerFunc: service.getUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "update_user",
			Method:      "PUT",
			Pattern:     "/users/{id:[0-9]+}",
			HandlerFunc: service.updateUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "invite_user",
			Method:      "POST",
			Pattern:     "/invitations",
			HandlerFunc: service.inviteUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "create_password_reset",
			Method:      "POST",
			Pattern:     "/password-reset",
			HandlerFunc: service.createPasswordResetHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
	}
}
