package subscriptions

import (
	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestFindActiveSubscriptionByUserID() {
	var (
		subscription *Subscription
		err          error
	)

	// First, try a user with an active subscription
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), suite.subscriptions[0].ID, subscription.ID)
		assert.Equal(suite.T(), suite.customers[0].ID, subscription.Customer.ID)
		assert.Equal(suite.T(), suite.users[0].ID, subscription.Customer.User.ID)
		assert.Equal(suite.T(), suite.plans[0].ID, subscription.Plan.ID)
	}

	// Second, try a user without subscription
	subscription, err = suite.service.FindActiveSubscriptionByUserID(suite.users[1].ID)

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserHasNoActiveSubscription, err)
	}
}
