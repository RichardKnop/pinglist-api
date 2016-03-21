package subscriptions

import (
	stripe "github.com/stripe/stripe-go"
)

// StripeAdapterInterface defines exported methods
type StripeAdapterInterface interface {
	// Exported methods
	CreateCustomer(stripeEmail, stripeToken string) (*stripe.Customer, error)
	GetCustomer(customerID string) (*stripe.Customer, error)
	GetOrCreateCustomer(customerID, stripeEmail, stripeToken string) (*stripe.Customer, error, bool)
	CreateSubscription(customerID, planID string) (*stripe.Sub, error)
	GetSubscription(subscriptionID string) (*stripe.Sub, error)
	CancelSubscription(subscriptionID, customerID string) (*stripe.Sub, error)
	GetEvent(eventID string) (*stripe.Event, error)
}
