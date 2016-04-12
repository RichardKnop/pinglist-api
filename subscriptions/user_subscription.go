package subscriptions

import (
	"errors"

	"github.com/RichardKnop/pinglist-api/subscriptions/subscriptionstatuses"
)

var (
	// ErrUserHasNoActiveSubscription ...
	ErrUserHasNoActiveSubscription = errors.New("User has no active subscription")
)

// FindActiveSubscriptionByUserID returns the currently active user subscription
func (s *Service) FindActiveSubscriptionByUserID(userID uint) (*Subscription, error) {
	// Fetch the subscription from the database
	subscription := new(Subscription)
	where := "subscription_customers.user_id = ? AND status != ? " +
		"AND cancelled_at IS NULL AND ended_at IS NULL"
	notFound := s.db.Preload("Customer.User").Preload("Plan").
		Joins("inner join subscription_customers on subscription_customers.id = subscription_subscriptions.customer_id").
		Joins("inner join account_users on account_users.id = subscription_customers.user_id").
		Where(where, userID, subscriptionstatuses.Cancelled).First(subscription).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserHasNoActiveSubscription
	}

	return subscription, nil
}
