package subscriptions

// CardRequest ...
type CardRequest struct {
	Token string `json:"token"`
}

// SubscriptionRequest ...
type SubscriptionRequest struct {
	PlanID uint   `json:"plan_id"`
	Token  string `json:"token"`
}
