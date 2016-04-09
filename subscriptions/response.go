package subscriptions

import (
	"fmt"
	"time"

	"github.com/RichardKnop/jsonhal"
)

// CardResponse ...
type CardResponse struct {
	jsonhal.Hal
	ID        uint   `json:"id"`
	Brand     string `json:"brand"`
	LastFour  string `json:"last_four"`
	ExpMonth  uint   `json:"exp_month"`
	ExpYear   uint   `json:"exp_year"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListCardsResponse ...
type ListCardsResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// PlanResponse ...
type PlanResponse struct {
	jsonhal.Hal
	ID                uint   `json:"id"`
	PlanID            string `json:"plan_id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Currency          string `json:"currency"`
	Amount            uint   `json:"amount"`
	TrialPeriod       uint   `json:"trial_period"`
	Interval          uint   `json:"interval"`
	MaxAlarms         uint   `json:"max_alarms"`
	MaxTeams          uint   `json:"max_teams"`
	MaxMembersPerTeam uint   `json:"max_members_per_team"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// ListPlansResponse ...
type ListPlansResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// SubscriptionResponse ...
type SubscriptionResponse struct {
	jsonhal.Hal
	ID             uint   `json:"id"`
	SubscriptionID string `json:"subscription_id"`
	StartedAt      string `json:"started_at"`
	CancelledAt    string `json:"cancelled_at"`
	EndedAt        string `json:"ended_at"`
	PeriodStart    string `json:"period_start"`
	PeriodEnd      string `json:"period_end"`
	TrialStart     string `json:"trial_start"`
	TrialEnd       string `json:"trial_end"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ListSubscriptionsResponse ...
type ListSubscriptionsResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// NewCardResponse creates new CardResponse instance
func NewCardResponse(card *Card) (*CardResponse, error) {
	response := &CardResponse{
		ID:        card.ID,
		Brand:     card.Brand,
		LastFour:  card.LastFour,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		CreatedAt: card.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: card.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/cards/%d", card.ID), // href
		"", // title
	)

	return response, nil
}

// NewListCardsResponse creates new ListCardsResponse instance
func NewListCardsResponse(count, page int, self, first, last, previous, next string, cards []*Card) (*ListCardsResponse, error) {
	response := &ListCardsResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of card responses
	cardResponses := make([]*CardResponse, len(cards))
	for i, card := range cards {
		cardResponse, err := NewCardResponse(card)
		if err != nil {
			return nil, err
		}
		cardResponses[i] = cardResponse
	}

	// Set embedded cards
	response.SetEmbedded(
		"cards",
		jsonhal.Embedded(cardResponses),
	)

	return response, nil
}

// NewPlanResponse creates new PlanResponse instance
func NewPlanResponse(plan *Plan) (*PlanResponse, error) {
	response := &PlanResponse{
		ID:                plan.ID,
		PlanID:            plan.PlanID,
		Name:              plan.Name,
		Description:       plan.Description.String,
		Currency:          plan.Currency,
		Amount:            plan.Amount,
		TrialPeriod:       plan.TrialPeriod,
		Interval:          plan.Interval,
		MaxAlarms:         plan.MaxAlarms,
		MaxTeams:          plan.MaxTeams,
		MaxMembersPerTeam: plan.MaxMembersPerTeam,
		CreatedAt:         plan.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:         plan.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/plans/%d", plan.ID), // href
		"", // title
	)

	return response, nil
}

// NewListPlansResponse creates new ListPlansResponse instance
func NewListPlansResponse(count, page int, self, first, last, previous, next string, plans []*Plan) (*ListPlansResponse, error) {
	response := &ListPlansResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of plan responses
	planResponses := make([]*PlanResponse, len(plans))
	for i, plan := range plans {
		planResponse, err := NewPlanResponse(plan)
		if err != nil {
			return nil, err
		}
		planResponses[i] = planResponse
	}

	// Set embedded plans
	response.SetEmbedded(
		"plans",
		jsonhal.Embedded(planResponses),
	)

	return response, nil
}

// NewSubscriptionResponse creates new SubscriptionResponse instance
func NewSubscriptionResponse(subscription *Subscription) (*SubscriptionResponse, error) {
	response := &SubscriptionResponse{
		ID:             subscription.ID,
		SubscriptionID: subscription.SubscriptionID,
		StartedAt:      subscription.StartedAt.Time.UTC().Format(time.RFC3339),
		CreatedAt:      subscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      subscription.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if subscription.CancelledAt.Valid {
		response.CancelledAt = subscription.CancelledAt.Time.UTC().Format(time.RFC3339)
	}
	if subscription.EndedAt.Valid {
		response.EndedAt = subscription.EndedAt.Time.UTC().Format(time.RFC3339)
	}
	if subscription.PeriodStart.Valid {
		response.PeriodStart = subscription.PeriodStart.Time.UTC().Format(time.RFC3339)
	}
	if subscription.PeriodEnd.Valid {
		response.PeriodEnd = subscription.PeriodEnd.Time.UTC().Format(time.RFC3339)
	}
	if subscription.TrialStart.Valid {
		response.TrialStart = subscription.TrialStart.Time.UTC().Format(time.RFC3339)
	}
	if subscription.TrialEnd.Valid {
		response.TrialEnd = subscription.TrialEnd.Time.UTC().Format(time.RFC3339)
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/subscriptions/%d", subscription.ID), // href
		"", // title
	)

	// Set embedded plan
	planResponse, err := NewPlanResponse(subscription.Plan)
	if err != nil {
		return nil, err
	}
	response.SetEmbedded(
		"plan",
		jsonhal.Embedded(planResponse),
	)

	// Set embedded card
	cardResponse, err := NewCardResponse(subscription.Card)
	if err != nil {
		return nil, err
	}
	response.SetEmbedded(
		"card",
		jsonhal.Embedded(cardResponse),
	)

	return response, nil
}

// NewListSubscriptionsResponse creates new ListSubscriptionsResponse instance
func NewListSubscriptionsResponse(count, page int, self, first, last, previous, next string, subscriptions []*Subscription) (*ListSubscriptionsResponse, error) {
	response := &ListSubscriptionsResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of subscription responses
	subscriptionResponses := make([]*SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		subscriptionResponse, err := NewSubscriptionResponse(subscription)
		if err != nil {
			return nil, err
		}
		subscriptionResponses[i] = subscriptionResponse
	}

	// Set embedded subscriptions
	response.SetEmbedded(
		"subscriptions",
		jsonhal.Embedded(subscriptionResponses),
	)

	return response, nil
}
