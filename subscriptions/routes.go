package subscriptions

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
			Name:        "list_plans",
			Method:      "GET",
			Pattern:     "/plans",
			HandlerFunc: service.listPlansHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "create_card",
			Method:      "POST",
			Pattern:     "/cards",
			HandlerFunc: service.createCardHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "get_card",
			Method:      "GET",
			Pattern:     "/cards/{id:[0-9]+}",
			HandlerFunc: service.getCardHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "delete_card",
			Method:      "DELETE",
			Pattern:     "/cards/{id:[0-9]+}",
			HandlerFunc: service.deleteCardHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "list_cards",
			Method:      "GET",
			Pattern:     "/cards",
			HandlerFunc: service.listCardsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "create_subscription",
			Method:      "POST",
			Pattern:     "/subscriptions",
			HandlerFunc: service.createSubscriptionHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "update_subscription",
			Method:      "PUT",
			Pattern:     "/subscriptions/{id:[0-9]+}",
			HandlerFunc: service.updateSubscriptionHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "cancel_subscription",
			Method:      "DELETE",
			Pattern:     "/subscriptions/{id:[0-9]+}",
			HandlerFunc: service.cancelSubscriptionHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "list_subscriptions",
			Method:      "GET",
			Pattern:     "/subscriptions",
			HandlerFunc: service.listSubscriptionsHandler,
			Middlewares: []negroni.Handler{
				accounts.NewUserAuthMiddleware(service.GetAccountsService()),
			},
		},
		routes.Route{
			Name:        "stripe_webhook",
			Method:      "POST",
			Pattern:     "/stripe-webhook",
			HandlerFunc: service.stripeWebhookHandler,
		},
	}
}
