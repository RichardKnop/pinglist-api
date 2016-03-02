package subscriptions

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface
	FindPlanByID(planID uint) (*Plan, error)
	FindSubscriptionByID(subscriptionID uint) (*Subscription, error)

	// Needed for the newRoutes to be able to register handlers
	listPlansHandler(w http.ResponseWriter, r *http.Request)
	subscribeUserHandler(w http.ResponseWriter, r *http.Request)
	stripeWebhookHandler(w http.ResponseWriter, r *http.Request)
}
