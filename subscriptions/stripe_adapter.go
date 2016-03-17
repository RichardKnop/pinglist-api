package subscriptions

import (
	"github.com/RichardKnop/pinglist-api/config"
	stripe "github.com/stripe/stripe-go"
	stripeCustomer "github.com/stripe/stripe-go/customer"
	stripeEvent "github.com/stripe/stripe-go/event"
	stripeSubscription "github.com/stripe/stripe-go/sub"
)

// StripeAdapter ...
type StripeAdapter struct {
}

// NewStripeAdapter starts a new StripeAdapter instance
func NewStripeAdapter(cnf *config.Config) *StripeAdapter {
	// Assign secret key from configuration to Stripe
	stripe.Key = cnf.Stripe.SecretKey

	return &StripeAdapter{}
}

// CreateSubscription creates a new Stripe customer and subscribes him/her to a plan
func (a *StripeAdapter) CreateSubscription(planID, stripeEmail, stripeToken string) (*stripe.Customer, error) {
	// Create a new Stripe customer and subscribe him/her to a plan
	params := &stripe.CustomerParams{
		Plan:  planID,
		Email: stripeEmail,
	}
	params.SetSource(stripeToken)
	return stripeCustomer.New(params)
}

// GetSubscription retrieves a subscription
func (a *StripeAdapter) GetSubscription(subscriptionID string) (*stripe.Sub, error) {
	return stripeSubscription.Get(subscriptionID, &stripe.SubParams{})
}

// CancelSubscription cancels a subscription
func (a *StripeAdapter) CancelSubscription(subscriptionID, customerID string) (*stripe.Sub, error) {
	return stripeSubscription.Cancel(
		subscriptionID,
		&stripe.SubParams{Customer: customerID},
	)
}

// GetEvent retrieves an event
func (a *StripeAdapter) GetEvent(eventID string) (*stripe.Event, error) {
	return stripeEvent.Get(eventID, &stripe.Params{})
}
