package subscriptions

import (
	"encoding/json"
	"time"

	"github.com/stretchr/testify/assert"
	stripe "github.com/stripe/stripe-go"
)

var customerSubscriptionUpdatedEvent = `{
  "id": "evt_17rUbtKkL3BsdwCi5snf6sTM",
  "livemode": false,
  "created": 1458572785,
  "data": {
    "object": {
      "id": "sub_87k6Km3feYTNdZ",
      "object": "subscription",
      "application_fee_percent": null,
      "cancel_at_period_end": false,
      "canceled_at": null,
      "current_period_end": 1461164783,
      "current_period_start": 1458572783,
      "customer": "cus_87k6TL1r0rosr5",
      "discount": null,
      "ended_at": null,
      "metadata": {},
      "plan": {
        "id": "professional",
        "object": "plan",
        "amount": 1000,
        "created": 1458548167,
        "currency": "usd",
        "interval": "month",
        "interval_count": 1,
        "livemode": false,
        "metadata": {},
        "name": "Professional",
        "statement_descriptor": null,
        "trial_period_days": 30
      },
      "quantity": 1,
      "start": 1458572785,
      "status": "trialing",
      "tax_percent": null,
      "trial_end": 1461164783,
      "trial_start": 1458572783
    },
    "previous_attributes": {
      "plan": {
        "amount": 250,
        "created": 1458548150,
        "id": "personal",
        "name": "Personal"
      },
      "start": 1458572783
    },
    "Obj": {
      "application_fee_percent": null,
      "cancel_at_period_end": false,
      "canceled_at": null,
      "current_period_end": 1461164783,
      "current_period_start": 1458572783,
      "customer": "cus_87k6TL1r0rosr5",
      "discount": null,
      "ended_at": null,
      "id": "sub_87k6Km3feYTNdZ",
      "metadata": {},
      "object": "subscription",
      "plan": {
        "amount": 1000,
        "created": 1458548167,
        "currency": "usd",
        "id": "professional",
        "interval": "month",
        "interval_count": 1,
        "livemode": false,
        "metadata": {},
        "name": "Professional",
        "object": "plan",
        "statement_descriptor": null,
        "trial_period_days": 30
      },
      "quantity": 1,
      "start": 1458572785,
      "status": "trialing",
      "tax_percent": null,
      "trial_end": 1461164783,
      "trial_start": 1458572783
    }
  },
  "pending_webhooks": 0,
  "type": "customer.subscription.updated",
  "request": "req_87k6cdJPtBZRks",
  "user_id": ""
}`

func (suite *SubscriptionsTestSuite) TestStripeEventCustomerSubscriptionUpdated() {
	// Unmarshal the JSON into a Stripe event
	stripeEvent := new(stripe.Event)
	// stripeEvent.Data = new(stripe.EventData)
	// err = stripeEvent.Data.UnmarshalJSON(contents)
	err := json.Unmarshal([]byte(customerSubscriptionUpdatedEvent), stripeEvent)
	assert.NoError(suite.T(), err, "Failed unmarshaling mock JSON into an event")

	// Create a test customer
	testCustomer := NewCustomer(suite.users[1], stripeEvent.GetObjValue("customer"))
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Failed to insert a test customer")

	// Create a test subscription
	startedAt := time.Unix(1458572785, 0)
	periodStart, periodEnd := time.Unix(1458572783, 0), time.Unix(1461164783, 0)
	trialStart, trialEnd := time.Unix(1458572783, 0), time.Unix(1461164783, 0)

	testSubscription := NewSubscription(
		testCustomer,
		suite.plans[0],
		stripeEvent.GetObjValue("id"),
		&startedAt,
		nil, // cancelled at
		nil, // ended at
		&periodStart,
		&periodEnd,
		&trialStart,
		&trialEnd,
	)
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Failed to insert a test subscription")

	// Fire off the event processing
	err = suite.service.stripeEventCustomerSubscriptionUpdated(stripeEvent)
	assert.Nil(suite.T(), err)

	// Fetch the updated subscription
	subscription := new(Subscription)
	notFound := suite.db.Preload("Customer.User").Preload("Plan").
		First(subscription, testSubscription.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Subscription plan and timestamps should have been updated
	assert.Equal(suite.T(), suite.plans[1].ID, uint(subscription.PlanID.Int64))
	assert.Equal(suite.T(), startedAt.UTC(), subscription.StartedAt.Time.UTC())
	assert.False(suite.T(), subscription.CancelledAt.Valid)
	assert.False(suite.T(), subscription.EndedAt.Valid)
	assert.Equal(suite.T(), periodStart.UTC(), subscription.PeriodStart.Time.UTC())
	assert.Equal(suite.T(), periodEnd.UTC(), subscription.PeriodEnd.Time.UTC())
	assert.Equal(suite.T(), trialStart.UTC(), subscription.TrialStart.Time.UTC())
	assert.Equal(suite.T(), trialEnd.UTC(), subscription.TrialEnd.Time.UTC())
}
