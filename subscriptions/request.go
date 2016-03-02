package subscriptions

// SubscriptionRequest ...
type SubscriptionRequest struct {
	StripeToken string `json:"stripe_token"`
	StripeEmail string `json:"stripe_email"`
	PlanID      uint   `json:"plan_id"`
}
