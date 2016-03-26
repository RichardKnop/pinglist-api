package alarms

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetMaxAlarmsNoSubscriptionUserInTrialPeriod() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User is not a member of a team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// User does not have a personal subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// User is in a free trial period
	user.CreatedAt = time.Now()

	// Max alarms should default to the free trial constant
	maxAlarms := suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), subscriptions.FreeTrialMaxAlarms, maxAlarms)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetMaxAlarmsNoSubscriptionUserNotInTrialPeriod() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User is not a member of a team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// User does not have a personal subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// User is no longer in a free trial period
	user.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)

	// Max alarms should be zero
	maxAlarms := suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), 0, maxAlarms)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetMaxAlarmsTeamWithSubscriptionUserWithoutSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User is a member of a team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		&teams.Team{Owner: &accounts.User{Model: gorm.Model{ID: 123}}},
		nil,
	)

	// The team has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		123,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxAlarms: 100,
			},
		},
		nil,
	)

	// Max alarms should be taken from the team subscription
	maxAlarms := suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), 100, maxAlarms)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetMaxAlarmsTeamWithoutSubscriptionUserWithSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User is a member of a team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		&teams.Team{Owner: &accounts.User{Model: gorm.Model{ID: 123}}},
		nil,
	)

	// The team has no active subscription
	suite.mockFindActiveSubscriptionByUserID(
		123,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// User has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxAlarms: 10,
			},
		},
		nil,
	)

	// Max alarms should be taken from the user subscription
	maxAlarms := suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), 10, maxAlarms)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetMaxAlarmsUserWithSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User is not a member of a team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// User has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxAlarms: 10,
			},
		},
		nil,
	)

	// Max alarms should be taken from the user subscription
	maxAlarms := suite.service.getMaxAlarms(user)
	assert.Equal(suite.T(), 10, maxAlarms)
}
