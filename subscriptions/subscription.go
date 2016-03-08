package subscriptions

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	stripe "github.com/stripe/stripe-go"
	stripeCustomer "github.com/stripe/stripe-go/customer"
	stripeSubscription "github.com/stripe/stripe-go/sub"
)

var (
	// ErrSubscriptionNotFound ...
	ErrSubscriptionNotFound = errors.New("Subscription not found")
)

// FindSubscriptionByID looks up a subscription by an ID and returns it
func (s *Service) FindSubscriptionByID(subscriptionID uint) (*Subscription, error) {
	// Fetch the subscription from the database
	subscription := new(Subscription)
	notFound := s.db.Preload("Customer").Preload("Plan").
		First(subscription, subscriptionID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrSubscriptionNotFound
	}

	return subscription, nil
}

// createSubscription creates a new Stripe user and subscribes him/her to a plan
func (s *Service) createSubscription(user *accounts.User, plan *Plan, stripeToken, stripeEmail string) (*Subscription, error) {
	// Begin a transaction
	tx := s.db.Begin()

	// Create a new Stripe customer and subscribe him/her to a plan
	params := &stripe.CustomerParams{
		Plan:  plan.PlanID,
		Email: stripeEmail,
		// TrialEnd: time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	params.SetSource(stripeToken)
	cus, err := stripeCustomer.New(params)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new incident object
	customer := newCustomer(user, cus.ID)

	// Save the customer to the database
	if err := tx.Create(customer).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Prepare subscription data
	sub := cus.Subs.Values[0]
	var (
		startedAt   *time.Time
		periodStart *time.Time
		periodEnd   *time.Time
		trialStart  *time.Time
		trialEnd    *time.Time
	)
	if sub.Start > 0 {
		t := time.Unix(sub.Start, 0)
		startedAt = &t
	}
	if sub.PeriodStart > 0 {
		t := time.Unix(sub.Start, 0)
		periodStart = &t
	}
	if sub.PeriodEnd > 0 {
		t := time.Unix(sub.Start, 0)
		periodEnd = &t
	}
	if sub.TrialStart > 0 {
		t := time.Unix(sub.Start, 0)
		trialStart = &t
	}
	if sub.TrialEnd > 0 {
		t := time.Unix(sub.Start, 0)
		trialEnd = &t
	}

	// Create a new subscription object
	subscription := newSubscription(
		customer,
		plan,
		sub.ID,
		startedAt,
		periodStart,
		periodEnd,
		trialStart,
		trialEnd,
	)

	// Save the subscription to the database
	if err := tx.Create(subscription).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return subscription, nil
}

// cancelSubscription cancells a subscription immediatelly
func (s *Service) cancelSubscription(subscription *Subscription) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Cancel the Stripe subscription
	sub, err := stripeSubscription.Cancel(
		subscription.SubscriptionID,
		&stripe.SubParams{Customer: subscription.Customer.CustomerID},
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Update the subscription's cancelled_at field
	cancelledAt := time.Unix(sub.Canceled, 0)
	if err := tx.Model(subscription).UpdateColumn(Subscription{
		CancelledAt: util.TimeOrNull(&cancelledAt),
	}).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}
