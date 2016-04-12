package subscriptions

import (
	"fmt"

	"github.com/RichardKnop/pinglist-api/config"
	stripe "github.com/stripe/stripe-go"
	stripeCard "github.com/stripe/stripe-go/card"
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

// CreateCustomer creates a new customer
func (a *StripeAdapter) CreateCustomer(email, token string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: email,
		Desc:  fmt.Sprintf("Customer for %s", email),
	}
	// Optionally add a payment source to the customer
	if token != "" {
		params.SetSource(token)
	}
	return stripeCustomer.New(params)
}

// GetCustomer retrieves a customer
func (a *StripeAdapter) GetCustomer(customerID string) (*stripe.Customer, error) {
	return stripeCustomer.Get(customerID, &stripe.CustomerParams{})
}

// GetOrCreateCustomer tries to retrieve a customer first, otherwise creates a new one
func (a *StripeAdapter) GetOrCreateCustomer(customerID, email, token string) (*stripe.Customer, bool, error) {
	var (
		c       *stripe.Customer
		created bool
		err     error
	)
	c, err = a.GetCustomer(customerID)
	if err != nil {
		c, err = a.CreateCustomer(email, token)
		created = true
	}
	return c, created, err
}

// CreateCard creates a new card
func (a *StripeAdapter) CreateCard(customerID, token string) (*stripe.Card, error) {
	params := &stripe.CardParams{
		Customer: customerID,
		Token:    token,
	}
	return stripeCard.New(params)
}

// DeleteCard deletes a card
func (a *StripeAdapter) DeleteCard(cardID, customerID string) (*stripe.Card, error) {
	params := &stripe.CardParams{
		Customer: customerID,
	}
	return stripeCard.Del(cardID, params)
}

// CreateSubscription creates a new subscription
func (a *StripeAdapter) CreateSubscription(customerID, planID string) (*stripe.Sub, error) {
	params := &stripe.SubParams{
		Customer: customerID,
		Plan:     planID,
	}
	return stripeSubscription.New(params)
}

// GetSubscription retrieves a subscription
func (a *StripeAdapter) GetSubscription(subscriptionID, customerID string) (*stripe.Sub, error) {
	params := &stripe.SubParams{
		Customer: customerID,
	}
	return stripeSubscription.Get(subscriptionID, params)
}

// UpdateSubscription upgrades or downgrades a subscription plan and/or
// changes the source associated with the subscription
func (a *StripeAdapter) UpdateSubscription(subscriptionID, customerID, planID string) (*stripe.Sub, error) {
	s, err := a.GetSubscription(subscriptionID, customerID)
	if err != nil {
		return nil, err
	}
	params := &stripe.SubParams{
		Customer: customerID,
		Plan:     planID,
	}
	return stripeSubscription.Update(s.ID, params)
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
