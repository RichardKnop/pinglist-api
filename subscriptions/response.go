package subscriptions

import (
	"fmt"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/util"
)

// CardResponse ...
type CardResponse struct {
	jsonhal.Hal
	ID        uint   `json:"id"`
	Brand     string `json:"brand"`
	Funding   string `json:"funding"`
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
	MinAlarmInterval  uint   `json:"min_alarm_interval"`
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
	Status         string `json:"status"`
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
		Funding:   card.Funding,
		LastFour:  card.LastFour,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		CreatedAt: util.FormatTime(card.CreatedAt),
		UpdatedAt: util.FormatTime(card.UpdatedAt),
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
		CreatedAt:         util.FormatTime(plan.CreatedAt),
		UpdatedAt:         util.FormatTime(plan.UpdatedAt),
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
		StartedAt:      util.FormatTime(subscription.StartedAt.Time),
		Status:         subscription.Status,
		CreatedAt:      util.FormatTime(subscription.CreatedAt),
		UpdatedAt:      util.FormatTime(subscription.UpdatedAt),
	}
	if subscription.CancelledAt.Valid {
		response.CancelledAt = util.FormatTime(subscription.CancelledAt.Time)
	}
	if subscription.EndedAt.Valid {
		response.EndedAt = util.FormatTime(subscription.EndedAt.Time)
	}
	if subscription.PeriodStart.Valid {
		response.PeriodStart = util.FormatTime(subscription.PeriodStart.Time)
	}
	if subscription.PeriodEnd.Valid {
		response.PeriodEnd = util.FormatTime(subscription.PeriodEnd.Time)
	}
	if subscription.TrialStart.Valid {
		response.TrialStart = util.FormatTime(subscription.TrialStart.Time)
	}
	if subscription.TrialEnd.Valid {
		response.TrialEnd = util.FormatTime(subscription.TrialEnd.Time)
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
