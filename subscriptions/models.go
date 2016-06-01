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
	MinAlarmInterval  uint
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
	Cards         []*Card
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
	Funding    string `sql:"type:varchar(10);not null"`
	LastFour   string `sql:"type:varchar(4);not null"`
	ExpMonth   uint   `sql:"not null"`
	ExpYear    uint   `sql:"not null"`
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
	SubscriptionID string      `sql:"type:varchar(60);unique;not null"`
	StartedAt      pq.NullTime `sql:"index"`
	CancelledAt    pq.NullTime `sql:"index"`
	EndedAt        pq.NullTime `sql:"index"`
	PeriodStart    pq.NullTime `sql:"index"`
	PeriodEnd      pq.NullTime `sql:"index"`
	TrialStart     pq.NullTime `sql:"index"`
	TrialEnd       pq.NullTime `sql:"index"`
	// trialing, active, past_due, canceled, or unpaid
	Status string `sql:"type:varchar(20);index;not null"`
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
	return customer
}

// NewCard creates new Card instance
func NewCard(customer *Customer, cardID, brand, funding, lastFour string, expMonth, expYear uint) *Card {
	customerID := util.PositiveIntOrNull(int64(customer.ID))
	card := &Card{
		CustomerID: customerID,
		CardID:     cardID,
		Brand:      brand,
		Funding:    funding,
		LastFour:   lastFour,
		ExpMonth:   expMonth,
		ExpYear:    expYear,
	}
	return card
}

// NewSubscription creates new Subscription instance
func NewSubscription(customer *Customer, plan *Plan, subscriptionID string, startedAt, cancelledAt, endedAt, periodStart, periodEnd, trialStart, trialEnd *time.Time, status string) *Subscription {
	customerID := util.PositiveIntOrNull(int64(customer.ID))
	planID := util.PositiveIntOrNull(int64(plan.ID))
	subscription := &Subscription{
		CustomerID:     customerID,
		PlanID:         planID,
		SubscriptionID: subscriptionID,
		StartedAt:      util.TimeOrNull(startedAt),
		CancelledAt:    util.TimeOrNull(cancelledAt),
		EndedAt:        util.TimeOrNull(endedAt),
		PeriodStart:    util.TimeOrNull(periodStart),
		PeriodEnd:      util.TimeOrNull(periodEnd),
		TrialStart:     util.TimeOrNull(trialStart),
		TrialEnd:       util.TimeOrNull(trialEnd),
		Status:         status,
	}
	return subscription
}
