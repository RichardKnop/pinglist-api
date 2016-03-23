package subscriptions

import (
	"testing"
	"time"

	"github.com/RichardKnop/pinglist-api/util"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionIsActive(t *testing.T) {
	var (
		endedAt      time.Time
		subscription *Subscription
	)

	// Subscription.EndedAt not set, therefor IsActive should return true
	subscription = new(Subscription)
	assert.True(t, subscription.IsActive())

	// Subscription.EndedAt is in the future, therefor IsActive should return true
	endedAt = time.Now().Add(+1 * time.Hour)
	subscription = &Subscription{
		EndedAt: util.TimeOrNull(&endedAt),
	}
	assert.True(t, subscription.IsActive())

	// Subscription.EndedAt is in the past, therefor IsActive should return false
	endedAt = time.Now().Add(-1 * time.Hour)
	subscription = &Subscription{
		EndedAt: util.TimeOrNull(&endedAt),
	}
	assert.False(t, subscription.IsActive())
}

func TestSubscriptionIsCancelled(t *testing.T) {
	var (
		cancelledAt  time.Time
		subscription *Subscription
	)

	// Subscription.CancelledAt not set, therefor IsCancelled should return false
	subscription = new(Subscription)
	assert.False(t, subscription.IsCancelled())

	// Subscription.CancelledAt is in the future, therefor IsCancelled should return false
	cancelledAt = time.Now().Add(+1 * time.Hour)
	subscription = &Subscription{
		CancelledAt: util.TimeOrNull(&cancelledAt),
	}
	assert.False(t, subscription.IsCancelled())

	// Subscription.CancelledAt is in thepast, therefor IsCancelled should return true
	cancelledAt = time.Now().Add(-1 * time.Hour)
	subscription = &Subscription{
		CancelledAt: util.TimeOrNull(&cancelledAt),
	}
	assert.True(t, subscription.IsActive())
}

func (suite *SubscriptionsTestSuite) TestFindSubscriptionByID() {
	var (
		subscription *Subscription
		err          error
	)

	// When we try to find a subscription with a bogus ID
	subscription, err = suite.service.FindSubscriptionByID(12345)

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrSubscriptionNotFound, err)
	}

	// When we try to find a plan with a valid ID
	subscription, err = suite.service.FindSubscriptionByID(suite.subscriptions[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), suite.subscriptions[0].ID, subscription.ID)
		assert.Equal(suite.T(), suite.customers[0].ID, subscription.Customer.ID)
		assert.Equal(suite.T(), suite.users[0].ID, subscription.Customer.User.ID)
		assert.Equal(suite.T(), suite.plans[0].ID, subscription.Plan.ID)
	}
}

func (suite *SubscriptionsTestSuite) TestFindSubscriptionBySubscriptionID() {
	var (
		subscription *Subscription
		err          error
	)

	// When we try to find a subscription with a bogus subscription ID
	subscription, err = suite.service.FindSubscriptionBySubscriptionID("bogus")

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrSubscriptionNotFound, err)
	}

	// When we try to find a plan with a valid subscription ID
	subscription, err = suite.service.FindSubscriptionBySubscriptionID(suite.subscriptions[0].SubscriptionID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct subscription object should be returned
	if assert.NotNil(suite.T(), subscription) {
		assert.Equal(suite.T(), suite.subscriptions[0].ID, subscription.ID)
		assert.Equal(suite.T(), suite.customers[0].ID, subscription.Customer.ID)
		assert.Equal(suite.T(), suite.users[0].ID, subscription.Customer.User.ID)
		assert.Equal(suite.T(), suite.plans[0].ID, subscription.Plan.ID)
	}
}

func (suite *SubscriptionsTestSuite) TestFindActiveUserSubscription() {
	var (
		subscription *Subscription
		err          error
	)

	// First, try a user with an active subscription
	subscription, err = suite.service.FindActiveUserSubscription(suite.users[0].ID)

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
	subscription, err = suite.service.FindActiveUserSubscription(suite.users[1].ID)

	// Subscription object should be nil
	assert.Nil(suite.T(), subscription)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserHasNoActiveSubscription, err)
	}
}

func (suite *SubscriptionsTestSuite) TestFindUserSubscriptions() {
	var (
		subscriptions []*Subscription
		err           error
	)

	// This should return 4 subscriptions, 1 active, 3 ended
	subscriptions, err = suite.service.findUserSubscriptions(suite.users[0].ID)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(subscriptions))
		assert.False(suite.T(), subscriptions[0].IsActive())
		assert.False(suite.T(), subscriptions[1].IsActive())
		assert.False(suite.T(), subscriptions[2].IsActive())
		assert.True(suite.T(), subscriptions[3].IsActive())
	}

	// This should return no subscriptions
	subscriptions, err = suite.service.findUserSubscriptions(suite.users[1].ID)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, len(subscriptions))
	}
}

func (suite *SubscriptionsTestSuite) TestPaginatedSubscriptionsCount() {
	var (
		count int
		err   error
	)

	count, err = suite.service.paginatedSubscriptionsCount(nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedSubscriptionsCount(suite.users[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedSubscriptionsCount(suite.users[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *SubscriptionsTestSuite) TestFindPaginatedSubscriptions() {
	var (
		subscriptions []*Subscription
		err           error
	)

	// This should return all subscriptions
	subscriptions, err = suite.service.findPaginatedSubscriptions(0, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(subscriptions))
		assert.Equal(suite.T(), suite.subscriptions[0].ID, subscriptions[0].ID)
		assert.Equal(suite.T(), suite.subscriptions[1].ID, subscriptions[1].ID)
		assert.Equal(suite.T(), suite.subscriptions[2].ID, subscriptions[2].ID)
		assert.Equal(suite.T(), suite.subscriptions[3].ID, subscriptions[3].ID)
	}

	// This should return all agencies ordered by ID desc
	subscriptions, err = suite.service.findPaginatedSubscriptions(0, 25, "id desc", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(subscriptions))
		assert.Equal(suite.T(), suite.subscriptions[3].ID, subscriptions[0].ID)
		assert.Equal(suite.T(), suite.subscriptions[2].ID, subscriptions[1].ID)
		assert.Equal(suite.T(), suite.subscriptions[1].ID, subscriptions[2].ID)
		assert.Equal(suite.T(), suite.subscriptions[0].ID, subscriptions[3].ID)
	}

	// Test offset
	subscriptions, err = suite.service.findPaginatedSubscriptions(2, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(subscriptions))
		assert.Equal(suite.T(), suite.subscriptions[2].ID, subscriptions[0].ID)
		assert.Equal(suite.T(), suite.subscriptions[3].ID, subscriptions[1].ID)
	}

	// Test limit
	subscriptions, err = suite.service.findPaginatedSubscriptions(2, 1, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(subscriptions))
		assert.Equal(suite.T(), suite.subscriptions[2].ID, subscriptions[0].ID)
	}
}
