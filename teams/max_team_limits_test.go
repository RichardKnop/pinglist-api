package teams

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/stretchr/testify/assert"
)

func (suite *TeamsTestSuite) TestGetMaxTeamLimitsUserWithoutSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// Max teams should be 0 and max members per team should be 0
	maxTeams, maxMembersPerTeam := suite.service.getMaxTeamLimits(user)
	assert.Equal(suite.T(), 0, maxTeams)
	assert.Equal(suite.T(), 0, maxMembersPerTeam)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *TeamsTestSuite) TestGetMaxTeamLimitsUserWithSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxTeams:          3,
				MaxMembersPerTeam: 10,
			},
		},
		nil,
	)

	// Max teams should be 3 and max members per team should be 10
	maxTeams, maxMembersPerTeam := suite.service.getMaxTeamLimits(user)
	assert.Equal(suite.T(), 3, maxTeams)
	assert.Equal(suite.T(), 10, maxMembersPerTeam)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}
