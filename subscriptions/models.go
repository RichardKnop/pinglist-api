package subscriptions

import (
	"database/sql"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// StripeEventLog ...
type StripeEventLog struct {
	gorm.Model
	EventID     string `sql:"type:varchar(60);unique;not null"`
	EventType   string `sql:"type:varchar(60);index;not null"`
	RequestDump string `sql:"type:text"`
	Processed   bool
}

// TableName specifies table name
func (l *StripeEventLog) TableName() string {
	return "subscription_stripe_event_logs"
}

// Plan ...
type Plan struct {
	gorm.Model
	PlanID            string         `sql:"type:varchar(60);unique;not null"`
	Name              string         `sql:"type:varchar(60);not null"`
	Description       sql.NullString `sql:"type:text"`
	Currency          string         `sql:"type:varchar(3);index;not null"`
	Amount            uint
	TrialPeriod       uint // days
	Interval          uint // days
	MaxAlarms         uint
	MaxTeams          uint
	MaxMembersPerTeam uint
}

// TableName specifies table name
func (p *Plan) TableName() string {
	return "subscription_plans"
}

// Customer ...
type Customer struct {
	gorm.Model
	UserID        sql.NullInt64 `sql:"index;not null"`
	User          *accounts.User
	CustomerID    string `sql:"type:varchar(60);unique;not null"`
	Subscriptions []*Subscription
}

// TableName specifies table name
func (c *Customer) TableName() string {
	return "subscription_customers"
}

// Card ...
type Card struct {
	gorm.Model
	CustomerID sql.NullInt64 `sql:"index;not null"`
	Customer   *Customer
	CardID     string `sql:"type:varchar(60);unique;not null"`
	Brand      string `sql:"type:varchar(20);not null"`
	LastFour   string `sql:"type:varchar(4);not null"`
}

// TableName specifies table name
func (c *Card) TableName() string {
	return "subscription_cards"
}

// Subscription ...
type Subscription struct {
	gorm.Model
	CustomerID     sql.NullInt64 `sql:"index;not null"`
	Customer       *Customer
	PlanID         sql.NullInt64 `sql:"index;not null"`
	Plan           *Plan
	CardID         sql.NullInt64 `sql:"index;not null"`
	Card           *Card
	SubscriptionID string      `sql:"type:varchar(60);unique;not null"`
	StartedAt      pq.NullTime `sql:"index"`
	CancelledAt    pq.NullTime `sql:"index"`
	EndedAt        pq.NullTime `sql:"index"`
	PeriodStart    pq.NullTime `sql:"index"`
	PeriodEnd      pq.NullTime `sql:"index"`
	TrialStart     pq.NullTime `sql:"index"`
	TrialEnd       pq.NullTime `sql:"index"`
}

// TableName specifies table name
func (s *Subscription) TableName() string {
	return "subscription_subscriptions"
}

// NewStripeEventLog creates new StripeEventLog instance
func NewStripeEventLog(eventID, eventType, requestDump string) *StripeEventLog {
	return &StripeEventLog{
		EventID:     eventID,
		EventType:   eventType,
		RequestDump: requestDump,
	}
}

// NewCustomer creates new Customer instance
func NewCustomer(user *accounts.User, customerID string) *Customer {
	userID := util.PositiveIntOrNull(int64(user.ID))
	customer := &Customer{
		UserID:     userID,
		CustomerID: customerID,
	}
	if userID.Valid {
		customer.User = user
	}
	return customer
}

// NewCard creates new Card instance
func NewCard(customer *Customer, cardID, brand, lastFour string) *Card {
	customerID := util.PositiveIntOrNull(int64(customer.ID))
	card := &Card{
		CustomerID: customerID,
		CardID:     cardID,
		Brand:      brand,
		LastFour:   lastFour,
	}
	if customerID.Valid {
		card.Customer = customer
	}
	return card
}

// NewSubscription creates new Subscription instance
func NewSubscription(customer *Customer, plan *Plan, card *Card, subscriptionID string, startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd *time.Time) *Subscription {
	customerID := util.PositiveIntOrNull(int64(customer.ID))
	planID := util.PositiveIntOrNull(int64(plan.ID))
	cardID := util.PositiveIntOrNull(int64(card.ID))
	subscription := &Subscription{
		CustomerID:     customerID,
		PlanID:         planID,
		CardID:         cardID,
		SubscriptionID: subscriptionID,
		StartedAt:      util.TimeOrNull(startedAt),
		CancelledAt:    util.TimeOrNull(cancelledAt),
		EndedAt:        util.TimeOrNull(endedAt),
		PeriodStart:    util.TimeOrNull(periodStart),
		PeriodEnd:      util.TimeOrNull(periodEnd),
		TrialStart:     util.TimeOrNull(trialStart),
		TrialEnd:       util.TimeOrNull(trialEnd),
	}
	if customerID.Valid {
		subscription.Customer = customer
	}
	if planID.Valid {
		subscription.Plan = plan
	}
	if cardID.Valid {
		subscription.Card = card
	}
	return subscription
}
