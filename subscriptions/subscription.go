package subscriptions

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	stripe "github.com/stripe/stripe-go"
)

var (
	// ErrSubscriptionNotFound ...
	ErrSubscriptionNotFound = errors.New("Subscription not found")
	// ErrUserCanOnlyHaveOneActiveSubscription ...
	ErrUserCanOnlyHaveOneActiveSubscription = errors.New("User can only have one active subscription")
)

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

// FindSubscriptionByCardID looks up a subscription by a card ID and returns it
func (s *Service) FindSubscriptionByCardID(cardID uint) (*Subscription, error) {
	// Fetch the subscription from the database
	subscription := new(Subscription)
	notFound := s.db.Preload("Customer.User").Preload("Plan").
		Where("card_id = ?", cardID).
		First(subscription).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrSubscriptionNotFound
	}

	return subscription, nil
}

// createSubscription creates a new Stripe user and subscribes him/her to a plan
func (s *Service) createSubscription(user *accounts.User, subscriptionRequest *SubscriptionRequest) (*Subscription, error) {
	// Fetch the active user subscription
	_, err := s.FindActiveSubscriptionByUserID(user.ID)

	// User should only have one active subscription at any time
	if err != ErrUserHasNoActiveSubscription {
		return nil, ErrUserCanOnlyHaveOneActiveSubscription
	}

	// Fetch the customer
	customer, err := s.FindCustomerByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Fetch the plan
	plan, err := s.FindPlanByID(subscriptionRequest.PlanID)
	if err != nil {
		return nil, err
	}

	// Calculate trial end
	trialPeriodDuration, err := s.calculateTrialPeriodDuration(customer, plan)
	if err != nil {
		return nil, err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Create a new Stripe subscription
	stripeSubscription, err := s.stripeAdapter.CreateSubscription(
		customer.CustomerID,
		plan.PlanID,
		time.Now().Add(trialPeriodDuration).Unix(),
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	logger.INFO.Printf("Created subscription: %s", stripeSubscription.ID)

	// Parse subscription times
	startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd := getStripeSubscriptionTimes(stripeSubscription)

	// Create a new subscription object
	subscription := NewSubscription(
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
		string(stripeSubscription.Status),
	)

	// Save the subscription to the database
	if err := tx.Create(subscription).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Assign related objects
	subscription.Customer = customer
	subscription.Plan = plan

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return subscription, nil
}

// changeSubscriptionPlan upgrades or downgrades a subscription plan
func (s *Service) updateSubscription(subscription *Subscription, subscriptionRequest *SubscriptionRequest) error {
	// Fetch the plan
	plan, err := s.FindPlanByID(subscriptionRequest.PlanID)
	if err != nil {
		return err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Change the subscription plan and card
	stripeSubscription, err := s.stripeAdapter.UpdateSubscription(
		subscription.SubscriptionID,
		subscription.Customer.CustomerID,
		plan.PlanID,
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Update the subscription
	err = s.updateSusbcriptionCommon(tx, subscription, plan, stripeSubscription)
	if err != nil {
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

// cancelSubscription cancels a subscription immediatelly
func (s *Service) cancelSubscription(subscription *Subscription) error {
	// Begin a transaction
	tx := s.db.Begin()

	logger.INFO.Print(subscription.SubscriptionID)

	// Cancel the subscription
	stripeSubscription, err := s.stripeAdapter.CancelSubscription(
		subscription.SubscriptionID,
		subscription.Customer.CustomerID,
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	logger.INFO.Printf("Cancelled subscription: %s", subscription.SubscriptionID)

	// Update the subscription
	err = s.updateSusbcriptionCommon(tx, subscription, subscription.Plan, stripeSubscription)
	if err != nil {
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

// updateSusbcriptionCommon updates a subscription
func (s *Service) updateSusbcriptionCommon(tx *gorm.DB, subscription *Subscription, plan *Plan, stripeSubscription *stripe.Sub) error {
	// Parse subscription times
	startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd := getStripeSubscriptionTimes(stripeSubscription)

	// Update the subscription plan
	if err := tx.Model(subscription).UpdateColumn(Subscription{
		PlanID:      util.PositiveIntOrNull(int64(plan.ID)),
		StartedAt:   util.TimeOrNull(startedAt),
		CancelledAt: util.TimeOrNull(cancelledAt),
		EndedAt:     util.TimeOrNull(endedAt),
		PeriodStart: util.TimeOrNull(periodStart),
		PeriodEnd:   util.TimeOrNull(periodEnd),
		TrialStart:  util.TimeOrNull(trialStart),
		TrialEnd:    util.TimeOrNull(trialEnd),
		Status:      string(stripeSubscription.Status),
		Model:       gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		return err
	}
	subscription.Plan = plan

	return nil
}

// subscriptionsCount returns a total count of subscriptions
// Can be optionally filtered by user
func (s *Service) subscriptionsCount(user *accounts.User) (int, error) {
	var count int
	if err := s.subscriptionsQuery(user).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedSubscriptions returns paginated subscription records
// Results can optionally be filtered by user
func (s *Service) findPaginatedSubscriptions(offset, limit int, orderBy string, user *accounts.User) ([]*Subscription, error) {
	var subscriptions []*Subscription

	// Get the pagination query
	subscriptionsQuery := s.subscriptionsQuery(user)

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

// subscriptionsQuery returns a generic db query for fetching subscriptions
func (s *Service) subscriptionsQuery(user *accounts.User) *gorm.DB {
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
