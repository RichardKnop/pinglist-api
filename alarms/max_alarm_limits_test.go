package alarms

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetMaxAlarms() {
	var (
		user      *accounts.User
		maxAlarms int
	)

	user = new(accounts.User)
	*user = *suite.users[1]

	// Mock find team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// Without a team and without a subscription, but the user is in free trial
	user.CreatedAt = time.Now()
	maxAlarms = suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), subscriptions.FreeTrialMaxAlarms, maxAlarms)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Mock find team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// Without a team and without a subscription, and no longer in free trial
	user.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
	maxAlarms = suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), 0, maxAlarms)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}
