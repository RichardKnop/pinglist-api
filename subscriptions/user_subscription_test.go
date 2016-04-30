package subscriptions

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestFindActiveSubscriptionByUserID() {
	var (
		subscription     *Subscription
		err              error
		testCustomer     *Customer
		testSubscription *Subscription
		end              = time.Now()
		start            = end.Add(-30 * 24 * time.Hour)
	)

	// Try user without any subscription
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[1].ID)

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserHasNoActiveSubscription, err)
	}

	// Now let's insert a test customer
	testCustomer = NewCustomer(suite.users[1], "new_customer_id")
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Failed to insert a test customer")

	// And also a test subscription
	testSubscription = NewSubscription(
		testCustomer,
		suite.plans[0],
		"new_subscription_id",
		&start, // started at
		nil,    // cancelled at
		nil,    // ended at
		&start, // period start
		&end,   // period end
		&start, // trial start
		&end,   // trial end
		"trialing",
	)
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Failed to insert a test subscription")

	// This time the user should have an active subscription
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), testSubscription.ID, subscription.ID)
	}

	// Second, try a user without subscription
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), testSubscription.ID, subscription.ID)
	}

	// A cancelled subscription is still active until the period end
	err = suite.db.Model(testSubscription).UpdateColumn("cancelled_at", start).Error
	assert.NoError(suite.T(), err, "Failed to update the test subscription")
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), testSubscription.ID, subscription.ID)
	}

	// Once a cancelled subscription reaches the period end, it's not longer active
	err = suite.db.Model(testSubscription).UpdateColumn("ended_at", end).Error
	assert.NoError(suite.T(), err, "Failed to update the test subscription")
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[1].ID)

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserHasNoActiveSubscription, err)
	}
}

func (suite *SubscriptionsTestSuite) TestCalculateTrialPeriodDuration() {
	var (
		trialPeriodDuration time.Duration
		testCustomer        *Customer
		testSubscription    *Subscription
		err                 error
		end                 = time.Now()
		start               = end.Add(-30 * 24 * time.Hour)
		startPlusFiveDays   = start.Add(5 * 24 * time.Hour)
	)

	// Insert a test customer
	testCustomer = NewCustomer(suite.users[1], "new_customer_id")
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Failed to insert a test customer")

	// Calculate the trial period duration
	trialPeriodDuration, err = suite.service.calculateTrialPeriodDuration(
		testCustomer,
		suite.plans[0],
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Full trial period duration should be returned
	if assert.NotNil(suite.T(), trialPeriodDuration) {
		assert.Equal(
			suite.T(),
			time.Duration(suite.plans[0].TrialPeriod)*24*time.Hour,
			trialPeriodDuration,
		)
	}

	// Insert a test subscription
	testSubscription = NewSubscription(
		testCustomer,
		suite.plans[0],
		"new_subscription_id",
		&start,             // started at
		&startPlusFiveDays, // cancelled at
		&end,               // ended at
		&start,             // period start
		&end,               // period end
		&start,             // trial start
		&end,               // trial end
		"trialing",
	)
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Failed to insert a test subscription")

	// Calculate the trial period duration
	trialPeriodDuration, err = suite.service.calculateTrialPeriodDuration(
		testCustomer,
		suite.plans[0],
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Trial period duration should be 5 days shorter now
	if assert.NotNil(suite.T(), trialPeriodDuration) {
		assert.Equal(
			suite.T(),
			time.Duration(suite.plans[0].TrialPeriod-5)*24*time.Hour,
			trialPeriodDuration,
		)
	}
}
