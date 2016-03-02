package subscriptions

import (
	stripe "github.com/stripe/stripe-go"
)

// customer.created
func (s *Service) stripeEventCustomerCreated(e *stripe.Event) error {
	// access event data via e.GetObjValue("resource_name_based_on_type", "resource_property_name")
	// alternatively you can access values via e.Data.Obj["resource_name_based_on_type"].(map[string]interface{})["resource_property_name"]

	// access previous attributes via e.GetPrevValue("resource_name_based_on_type", "resource_property_name")
	// alternatively you can access values via e.Data.Prev["resource_name_based_on_type"].(map[string]interface{})["resource_property_name"]

	return nil
}

// customer.subscription.created
func (s *Service) stripeEventCustomerSubscriptionCreated(e *stripe.Event) error {
	// e.GetObjValue("subscription", "id")
	// Customer ID:
	// e.GetObjValue("subscription", "customer")

	return nil
}

// customer.subscription.trial_will_end
func (s *Service) stripeEventCustomerSubscriptionTrialWillEnd(e *stripe.Event) error {
	// Sent 3 days before the trial period ends
	// e.GetObjValue("subscription", "id")

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
func (s *Service) stripeEventCustomerSubscriptionUpdated(e *stripe.Event) error {
	// Also received when subscription plan is upgraded or downgraded

	return nil
}

// customer.subscription.deleted
func (s *Service) stripeEventCustomerSubscriptionDeleted(e *stripe.Event) error {
	// When customer subscription is cancelled immediately

	return nil
}
