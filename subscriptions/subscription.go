package subscriptions

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	stripe "github.com/stripe/stripe-go"
)

var (
	// ErrSubscriptionNotFound ...
	ErrSubscriptionNotFound = errors.New("Subscription not found")
	// ErrUserHasNoActiveSubscription ...
	ErrUserHasNoActiveSubscription = errors.New("User has no active subscription")
	// ErrUserCanOnlyHaveOneActiveSubscription ...
	ErrUserCanOnlyHaveOneActiveSubscription = errors.New("User can only have one active subscriptions")
)

// IsActive returns true if the subscription has not ended yet
func (s *Subscription) IsActive() bool {
	return !s.EndedAt.Valid || s.EndedAt.Time.After(time.Now())
}

// FindSubscriptionByID looks up a subscription by an ID and returns it
func (s *Service) FindSubscriptionByID(subscriptionID uint) (*Subscription, error) {
	// Fetch the subscription from the database
	subscription := new(Subscription)
	notFound := s.db.Preload("Customer.User").Preload("Plan").
		First(subscription, subscriptionID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrSubscriptionNotFound
	}

	return subscription, nil
}

// FindSubscriptionBySubscriptionID looks up a subscription by a subscription ID and returns it
func (s *Service) FindSubscriptionBySubscriptionID(subscriptionID string) (*Subscription, error) {
	// Fetch the subscription from the database
	subscription := new(Subscription)
	notFound := s.db.Preload("Customer.User").Preload("Plan").
		Where("subscription_id = ?", subscriptionID).
		First(subscription).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrSubscriptionNotFound
	}

	return subscription, nil
}

// FindActiveUserSubscription returns the currently active user subscription
func (s *Service) FindActiveUserSubscription(userID uint) (*Subscription, error) {
	// Fetch all active user subscriptions
	activeUserSubscriptions, err := s.findActiveUserSubscriptions(userID)
	if err != nil {
		return nil, err
	}

	// User has no active subscription
	if len(activeUserSubscriptions) == 0 {
		return nil, ErrUserHasNoActiveSubscription
	}

	return activeUserSubscriptions[0], nil
}

// createSubscription creates a new Stripe user and subscribes him/her to a plan
func (s *Service) createSubscription(user *accounts.User, plan *Plan, stripeToken, stripeEmail string) (*Subscription, error) {
	// Fetch all active user subscriptions
	activeUserSubscriptions, err := s.findActiveUserSubscriptions(user.ID)
	if err != nil {
		return nil, err
	}

	// User should only have one active subscription at any time
	if len(activeUserSubscriptions) > 0 {
		return nil, ErrUserCanOnlyHaveOneActiveSubscription
	}

	// Create a new Stripe customer and subscribe him/her to a plan
	stripeCustomer, err := s.stripeAdapter.CreateSubscription(
		plan.PlanID,
		stripeEmail,
		stripeToken,
	)
	if err != nil {
		return nil, err
	}

	// Create a new customer object
	customer := newCustomer(user, stripeCustomer.ID)

	// Begin a transaction
	tx := s.db.Begin()

	// Save the customer to the database
	if err := tx.Create(customer).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Assign Stripe subscription to a less confusing variable
	stripeSubscription := stripeCustomer.Subs.Values[0]

	// Parse subscription times
	startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd := getStripeSubscriptionTimes(stripeSubscription)

	// Create a new subscription object
	subscription := newSubscription(
		customer,
		plan,
		stripeSubscription.ID,
		startedAt,
		cancelledAt,
		endedAt,
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

// getStripeSubscriptionTimes parses UNIX timestamps from a subscription
func getStripeSubscriptionTimes(stripeSubscription *stripe.Sub) (startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd *time.Time) {
	if stripeSubscription.Start > 0 {
		t := time.Unix(stripeSubscription.Start, 0)
		startedAt = &t
	}
	if stripeSubscription.Canceled > 0 {
		t := time.Unix(stripeSubscription.Canceled, 0)
		cancelledAt = &t
	}
	if stripeSubscription.Ended > 0 {
		t := time.Unix(stripeSubscription.Ended, 0)
		endedAt = &t
	}
	if stripeSubscription.PeriodStart > 0 {
		t := time.Unix(stripeSubscription.PeriodStart, 0)
		periodStart = &t
	}
	if stripeSubscription.PeriodEnd > 0 {
		t := time.Unix(stripeSubscription.PeriodEnd, 0)
		periodEnd = &t
	}
	if stripeSubscription.TrialStart > 0 {
		t := time.Unix(stripeSubscription.TrialStart, 0)
		trialStart = &t
	}
	if stripeSubscription.TrialEnd > 0 {
		t := time.Unix(stripeSubscription.TrialEnd, 0)
		trialEnd = &t
	}
	return
}

// cancelSubscription cancells a subscription immediatelly
func (s *Service) cancelSubscription(subscription *Subscription) error {
	// Cancel the subscription
	stripeSubscription, err := s.stripeAdapter.CancelSubscription(
		subscription.SubscriptionID,
		subscription.Customer.CustomerID,
	)
	if err != nil {
		return err
	}

	// Update the subscription's cancelled_at field
	cancelledAt := time.Unix(stripeSubscription.Canceled, 0)
	if err := s.db.Model(subscription).UpdateColumn(Subscription{
		CancelledAt: util.TimeOrNull(&cancelledAt),
	}).Error; err != nil {
		return err
	}

	return nil
}

// findActiveUserSubscriptions returns only active subscriptions belonging to a user
func (s *Service) findActiveUserSubscriptions(userID uint) ([]*Subscription, error) {
	var activeUserSubscriptions []*Subscription

	// Fetch all user subscriptions first
	userSubscriptions, err := s.findUserSubscriptions(userID)
	if err != nil {
		return activeUserSubscriptions, err
	}

	// Filter out active subscriptions only
	for _, userSubscription := range userSubscriptions {
		if userSubscription.IsActive() {
			activeUserSubscriptions = append(activeUserSubscriptions, userSubscription)
		}
	}

	return activeUserSubscriptions, err
}

// findUserSubscriptions returns subscriptions belonging to a user
func (s *Service) findUserSubscriptions(userID uint) ([]*Subscription, error) {
	var userSubscriptions []*Subscription
	userObj := &accounts.User{Model: gorm.Model{ID: userID}}
	return userSubscriptions, s.paginatedSubscriptionsQuery(userObj).
		Preload("Customer.User").Preload("Plan").
		Order("id").Find(&userSubscriptions).Error
}

// paginatedSubscriptionsCount returns a total count of subscriptions
// Can be optionally filtered by user
func (s *Service) paginatedSubscriptionsCount(user *accounts.User) (int, error) {
	var count int
	if err := s.paginatedSubscriptionsQuery(user).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedSubscriptions returns paginated subscription records
// Results can optionally be filtered by user
func (s *Service) findPaginatedSubscriptions(offset, limit int, orderBy string, user *accounts.User) ([]*Subscription, error) {
	var subscriptions []*Subscription

	// Get the pagination query
	subscriptionsQuery := s.paginatedSubscriptionsQuery(user)

	// Default ordering
	if orderBy == "" {
		orderBy = "id"
	}

	// Retrieve paginated results from the database
	err := subscriptionsQuery.Offset(offset).Limit(limit).Order(orderBy).
		Preload("Customer.User").Preload("Plan").Find(&subscriptions).Error
	if err != nil {
		return subscriptions, err
	}

	return subscriptions, nil
}

// paginatedSubscriptionsQuery returns a db query for paginated subscriptions
func (s *Service) paginatedSubscriptionsQuery(user *accounts.User) *gorm.DB {
	// Basic query
	subscriptionsQuery := s.db.Model(new(Subscription))

	// Optionally filter by user
	if user != nil {
		subscriptionsQuery = subscriptionsQuery.
			Joins("inner join subscription_customers on subscription_customers.id = subscription_subscriptions.customer_id").
			Joins("inner join account_users on account_users.id = subscription_customers.user_id").
			Where("subscription_customers.user_id = ?", user.ID)
	}

	return subscriptionsQuery
}
