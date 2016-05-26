package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetAlarmLimitsFreeTier() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User does not have an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	alarmLimits := suite.service.getAlarmLimits(
		nil, // team
		user,
	)
	// Max alarms should default to the free tier value
	assert.Equal(suite.T(), FreeTierMaxAlarms, alarmLimits.maxAlarms)
	// Min alarm interval should default to the free tier value
	assert.Equal(suite.T(), FreeTierMinAlarmInterval, alarmLimits.minAlarmInterval)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestgetAlarmLimitsTeamWithSubscriptionUserWithoutSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// The team has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		123,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxAlarms:        100,
				MinAlarmInterval: 40,
			},
		},
		nil,
	)

	alarmLimits := suite.service.getAlarmLimits(
		&teams.Team{Owner: &accounts.User{Model: gorm.Model{ID: 123}}},
		user,
	)
	// Max alarms should be taken from the team plan
	assert.Equal(suite.T(), uint(100), alarmLimits.maxAlarms)
	// Min alarm interval should be taken from the team plan
	assert.Equal(suite.T(), uint(40), alarmLimits.minAlarmInterval)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestgetAlarmLimitsTeamWithoutSubscriptionUserWithSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

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
				MaxAlarms:        10,
				MinAlarmInterval: 50,
			},
		},
		nil,
	)

	alarmLimits := suite.service.getAlarmLimits(
		&teams.Team{Owner: &accounts.User{Model: gorm.Model{ID: 123}}},
		user,
	)
	// Max alarms should be taken from the user plan
	assert.Equal(suite.T(), uint(10), alarmLimits.maxAlarms)
	// Min alarm interval should be taken from the user plan
	assert.Equal(suite.T(), uint(50), alarmLimits.minAlarmInterval)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetAlarmLimitsUserWithSubscription() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User has an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxAlarms:        20,
				MinAlarmInterval: 40,
			},
		},
		nil,
	)

	alarmLimits := suite.service.getAlarmLimits(
		nil, // team
		user,
	)
	// Max alarms should be taken from the user plan
	assert.Equal(suite.T(), uint(20), alarmLimits.maxAlarms)
	// Min alarm interval should be taken from the user plan
	assert.Equal(suite.T(), uint(40), alarmLimits.minAlarmInterval)
}
