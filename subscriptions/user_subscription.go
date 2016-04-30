package subscriptions

import (
	"errors"
	"time"
)

var (
	// ErrUserHasNoActiveSubscription ...
	ErrUserHasNoActiveSubscription = errors.New("User has no active subscription")
)

// FindActiveSubscriptionByUserID returns the currently active user subscription
func (s *Service) FindActiveSubscriptionByUserID(userID uint) (*Subscription, error) {
	// Fetch the subscription from the database
	subscription := new(Subscription)
	where := "subscription_customers.user_id = ? AND ended_at IS NULL"
	notFound := s.db.Preload("Customer.User").Preload("Plan").
		Joins("inner join subscription_customers on subscription_customers.id = subscription_subscriptions.customer_id").
		Joins("inner join account_users on account_users.id = subscription_customers.user_id").
		Where(where, userID).Last(subscription).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserHasNoActiveSubscription
	}

	return subscription, nil
}

// calculateTrialPeriodDuration calculates duration of a trial period
// If the user has been subscribed to the same plan before, the trial period
// gets decreased by the time already spent in the trialing mode
func (s *Service) calculateTrialPeriodDuration(customer *Customer, plan *Plan) (time.Duration, error) {
	trialPeriodDurarion := time.Duration(plan.TrialPeriod) * time.Hour * 24
	var prevSubscriptions []*Subscription
	if err := s.db.Where("cancelled_at IS NOT NULL").Where(map[string]interface{}{
		"plan_id":     plan.ID,
		"customer_id": customer.ID,
	}).Find(&prevSubscriptions).Error; err != nil {
		return 0, err
	}
	for _, prevSubscription := range prevSubscriptions {
		delta := prevSubscription.CancelledAt.Time.Sub(prevSubscription.StartedAt.Time)
		trialPeriodDurarion -= delta
	}
	if trialPeriodDurarion < 0 {
		trialPeriodDurarion = 0
	}
	return trialPeriodDurarion, nil
}
