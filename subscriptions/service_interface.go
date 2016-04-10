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
	FindCustomerByUserID(userID uint) (*Customer, error)
	FindCustomerByCustomerID(customerID string) (*Customer, error)
	FindCardByID(cardID uint) (*Card, error)
	FindCardByCardID(cardID string) (*Card, error)
	FindSubscriptionByID(subscriptionID uint) (*Subscription, error)
	FindSubscriptionBySubscriptionID(subscriptionID string) (*Subscription, error)
	FindActiveSubscriptionByUserID(userID uint) (*Subscription, error)

	// Needed for the newRoutes to be able to register handlers
	listPlansHandler(w http.ResponseWriter, r *http.Request)
	createCardHandler(w http.ResponseWriter, r *http.Request)
	getCardHandler(w http.ResponseWriter, r *http.Request)
	listCardsHandler(w http.ResponseWriter, r *http.Request)
	deleteCardHandler(w http.ResponseWriter, r *http.Request)
	createSubscriptionHandler(w http.ResponseWriter, r *http.Request)
	getSubscriptionHandler(w http.ResponseWriter, r *http.Request)
	listSubscriptionsHandler(w http.ResponseWriter, r *http.Request)
	updateSubscriptionHandler(w http.ResponseWriter, r *http.Request)
	cancelSubscriptionHandler(w http.ResponseWriter, r *http.Request)
	stripeWebhookHandler(w http.ResponseWriter, r *http.Request)
}
