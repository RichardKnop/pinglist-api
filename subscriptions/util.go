package subscriptions

import (
	"time"

	stripe "github.com/stripe/stripe-go"
)

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
