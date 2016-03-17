package subscriptions

import (
	stripe "github.com/stripe/stripe-go"
)

// StripeAdapterInterface defines exported methods
type StripeAdapterInterface interface {
	// Exported methods
	CreateSubscription(planID, stripeEmail, stripeToken string) (*stripe.Customer, error)
	GetSubscription(subscriptionID string) (*stripe.Sub, error)
	CancelSubscription(subscriptionID, customerID string) (*stripe.Sub, error)
	GetEvent(eventID string) (*stripe.Event, error)
}
