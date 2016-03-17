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
	FindPlanByPlanID(planID string) (*Plan, error)
	FindCustomerByID(customerID uint) (*Customer, error)
	FindCustomerByCustomerID(customerID string) (*Customer, error)
	FindSubscriptionByID(subscriptionID uint) (*Subscription, error)
	FindSubscriptionBySubscriptionID(subscriptionID string) (*Subscription, error)
	FindActiveUserSubscription(userID uint) (*Subscription, error)

	// Needed for the newRoutes to be able to register handlers
	listPlansHandler(w http.ResponseWriter, r *http.Request)
	subscribeUserHandler(w http.ResponseWriter, r *http.Request)
	listSubscriptionsHandler(w http.ResponseWriter, r *http.Request)
	stripeWebhookHandler(w http.ResponseWriter, r *http.Request)
}
