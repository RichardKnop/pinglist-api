package subscriptions

import (
	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestFindSubscriptionByID() {
	var (
		subscription *Subscription
		err  error
	)

	// When we try to find a subscription with a bogus ID
	subscription, err = suite.service.FindSubscriptionByID(12345)

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errSubscriptionNotFound, err)
	}

	// When we try to find a plan with a valid ID
	subscription, err = suite.service.FindSubscriptionByID(suite.subscriptions[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), suite.subscriptions[0].ID, subscription.ID)
    assert.Equal(suite.T(), suite.customers[0].ID, subscription.Customer.ID)
    assert.Equal(suite.T(), suite.plans[0].ID, subscription.Plan.ID)
	}
}
