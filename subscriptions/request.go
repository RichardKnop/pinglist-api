package subscriptions

// CardRequest ...
type CardRequest struct {
	Token string `json:"token"`
}

// SubscriptionRequest ...
type SubscriptionRequest struct {
	StripeToken string `json:"stripe_token"`
	StripeEmail string `json:"stripe_email"`
	PlanID      uint   `json:"plan_id"`
}
