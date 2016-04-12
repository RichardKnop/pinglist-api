package subscriptions

import (
	stripe "github.com/stripe/stripe-go"
)

// customer.subscription.updated
// Occurs whenever a subscription changes. Examples would include switching
// from one plan to another, or switching status from trial to active.
func (s *Service) stripeEventCustomerSubscriptionUpdated(e *stripe.Event) error {
	// Fetch the subscription record from our database
	subscription, err := s.FindSubscriptionBySubscriptionID(e.GetObjValue("id"))
	if err != nil {
		return err
	}

	// Unmarshal into Stripe subscription
	stripeSubscription := new(stripe.Sub)
	if err := stripeSubscription.UnmarshalJSON(e.Data.Raw); err != nil {
		return err
	}

	// Fetch the plan
	plan, err := s.FindPlanByPlanID(stripeSubscription.Plan.ID)
	if err != nil {
		return err
	}

	// Update the subscription
	return s.updateSusbcriptionCommon(
		s.db,
		subscription,
		plan,
		stripeSubscription,
	)
}

// customer.subscription.deleted
// Occurs whenever a customer ends their subscription.
func (s *Service) stripeEventCustomerSubscriptionDeleted(e *stripe.Event) error {
	// Fetch the subscription record from our database
	subscription, err := s.FindSubscriptionBySubscriptionID(e.GetObjValue("id"))
	if err != nil {
		return err
	}

	// Unmarshal into Stripe subscription
	stripeSubscription := new(stripe.Sub)
	if err := stripeSubscription.UnmarshalJSON(e.Data.Raw); err != nil {
		return err
	}

	// Update the subscription
	return s.updateSusbcriptionCommon(
		s.db,
		subscription,
		subscription.Plan,
		stripeSubscription,
	)
}

// customer.subscription.trial_will_end
// Sent 3 days before the trial period ends.
func (s *Service) stripeEventCustomerSubscriptionTrialWillEnd(e *stripe.Event) error {
	// TODO - probably send a remainder email here
	return nil
}
