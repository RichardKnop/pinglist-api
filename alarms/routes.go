package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/routes"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the alarms service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the alarms service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "list_regions",
			Method:      "GET",
			Pattern:     "/regions",
			HandlerFunc: service.listRegionsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "create_alarm",
			Method:      "POST",
			Pattern:     "/alarms",
			HandlerFunc: service.createAlarmHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "update_alarm",
			Method:      "PUT",
			Pattern:     "/alarms/{id:[0-9]+}",
			HandlerFunc: service.updateAlarmHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "delete_alarm",
			Method:      "DELETE",
			Pattern:     "/alarms/{id:[0-9]+}",
			HandlerFunc: service.deleteAlarmHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "list_alarms",
			Method:      "GET",
			Pattern:     "/alarms",
			HandlerFunc: service.listAlarmsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "list_alarm_incidents",
			Method:      "GET",
			Pattern:     "/alarms/{id:[0-9]+}/incidents",
			HandlerFunc: service.listAlarmIncidentsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "list_alarm_results",
			Method:      "GET",
			Pattern:     "/alarms/{id:[0-9]+}/results",
			HandlerFunc: service.listAlarmResultsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
	}
}
