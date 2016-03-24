package metrics

import (
	"github.com/RichardKnop/pinglist-api/routes"
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
    // TODO
	}
}
