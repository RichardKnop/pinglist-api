package subscriptions

// CardRequest ...
type CardRequest struct {
	Token string `json:"token"`
}

// SubscriptionRequest ...
type SubscriptionRequest struct {
	PlanID uint `json:"plan_id"`
	CardID uint `json:"card_id"`
}
