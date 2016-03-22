package subscriptions

import (
	"time"

	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	stripe "github.com/stripe/stripe-go"
)

// customer.created
func (s *Service) stripeEventCustomerCreated(e *stripe.Event) error {
	return nil
}

// customer.subscription.created
func (s *Service) stripeEventCustomerSubscriptionCreated(e *stripe.Event) error {
	return nil
}

// customer.subscription.trial_will_end
// // Sent 3 days before the trial period ends
func (s *Service) stripeEventCustomerSubscriptionTrialWillEnd(e *stripe.Event) error {
	return nil
}

// invoice.created
func (s *Service) stripeEventInvoiceCreated(e *stripe.Event) error {
	// Once the trial period is up, Stripe will generate an invoice and send out
	// an invoice.created event notification. Approximately an hour later,
	// Stripe will attempt to charge that invoice. Assuming that the payment
	// attempt succeeded, you’ll receive notifications of the following events:
	// - charge.succeeded
	// - invoice.payment_succeeded
	// - customer.subscription.updated (reflecting an update from a trial to an active subscription)
	// Customer ID:
	// e.GetObjValue("invoice", "customer")
	// Subscription ID:
	// e.GetObjValue("invoice", "subscription")

	return nil
}

// charge.succeeded
func (s *Service) stripeEventChargeSucceeded(e *stripe.Event) error {
	return nil
}

// invoice.payment_succeeded
func (s *Service) stripeEventPaymentSucceeded(e *stripe.Event) error {
	return nil
}

// customer.subscription.updated
// Also received when subscription plan is upgraded or downgraded
func (s *Service) stripeEventCustomerSubscriptionUpdated(e *stripe.Event) error {
	// Fetch the subscription record from our database
	subscriptionID := e.GetObjValue("subscription", "id")
	subscription, err := s.FindSubscriptionBySubscriptionID(subscriptionID)
	if err != nil {
		return err
	}

	// Verify the subscription by fetching it from Stripe
	stripeSubscription, err := s.stripeAdapter.GetSubscription(
		subscription.SubscriptionID,
		subscription.Customer.CustomerID,
	)
	if err != nil {
		return err
	}

	// Fetch the plan
	plan, err := s.FindPlanByPlanID(stripeSubscription.Plan.ID)
	if err != nil {
		return err
	}

	// Parse subscription times
	startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd := getStripeSubscriptionTimes(stripeSubscription)

	if plan.ID != subscription.Plan.ID {
		// Plan changed (upgraded or downgraded)
	}

	if cancelledAt != nil && !subscription.CancelledAt.Valid {
		// Subscription cancelled
	}

	if endedAt != nil && !subscription.EndedAt.Valid {
		// Subscription ended
	}

	// Update the subscription
	if err := s.db.Model(subscription).UpdateColumn(Subscription{
		PlanID:      util.PositiveIntOrNull(int64(plan.ID)),
		StartedAt:   util.TimeOrNull(startedAt),
		CancelledAt: util.TimeOrNull(cancelledAt),
		EndedAt:     util.TimeOrNull(endedAt),
		PeriodStart: util.TimeOrNull(periodStart),
		PeriodEnd:   util.TimeOrNull(periodEnd),
		TrialStart:  util.TimeOrNull(trialStart),
		TrialEnd:    util.TimeOrNull(trialEnd),
		Model:       gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		return nil
	}

	return nil
}

// customer.subscription.deleted
// When customer subscription is cancelled immediately
func (s *Service) stripeEventCustomerSubscriptionDeleted(e *stripe.Event) error {
	// Fetch the subscription record from our database
	subscriptionID := e.GetObjValue("subscription", "id")
	subscription, err := s.FindSubscriptionBySubscriptionID(subscriptionID)
	if err != nil {
		return err
	}

	// Verify the subscription by fetching it from Stripe
	stripeSubscription, err := s.stripeAdapter.GetSubscription(
		subscription.SubscriptionID,
		subscription.Customer.CustomerID,
	)
	if err != nil {
		return err
	}

	// Update the subscription's cancelled_at field
	cancelledAt := time.Unix(stripeSubscription.Canceled, 0)
	if err := s.db.Model(subscription).UpdateColumns(Subscription{
		CancelledAt: util.TimeOrNull(&cancelledAt),
		Model:       gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		return err
	}

	return nil
}
