package subscriptions

import (
	stripe "github.com/stripe/stripe-go"
)

// StripeAdapterInterface defines exported methods
type StripeAdapterInterface interface {
	// Exported methods
	CreateCustomer(stripeEmail, stripeToken string) (*stripe.Customer, error)
	GetCustomer(customerID string) (*stripe.Customer, error)
	GetOrCreateCustomer(customerID, stripeEmail, stripeToken string) (*stripe.Customer, bool, error)
	CreateCard(customerID, stripeToken string) (*stripe.Card, error)
	DeleteCard(customerID, cardID string) (*stripe.Card, error)
	CreateSubscription(customerID, planID string) (*stripe.Sub, error)
	GetSubscription(subscriptionID, customerID string) (*stripe.Sub, error)
	ChangeSubscriptionPlan(subscriptionID, customerID, planID string) (*stripe.Sub, error)
	CancelSubscription(subscriptionID, customerID string) (*stripe.Sub, error)
	GetEvent(eventID string) (*stripe.Event, error)
}
