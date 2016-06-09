package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetAlarmLimitsFreeTier() {
	user := new(accounts.User)
	*user = *suite.users[1]

	// User is not a member of a team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// User does not have an active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	alarmLimits := suite.service.getAlarmLimits(user)

	// Max alarms should default to the free tier value
	assert.Equal(suite.T(), FreeTierMaxAlarms, alarmLimits.maxAlarms)
	// Min alarm interval should default to the free tier value
	assert.Equal(suite.T(), FreeTierMinAlarmInterval, alarmLimits.minAlarmInterval)
	// Unlimited emails should default to false
	assert.False(suite.T(), alarmLimits.unlimitedEmails)
	// Max emails per interval value should default to the free tier value
	assert.Equal(suite.T(), FreeTierMaxEmailsPerInterval, alarmLimits.maxEmailsPerInterval)
	// Slack alerts should default to false
	assert.False(suite.T(), alarmLimits.slackAlerts)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetAlarmLimitsTeamWithSubscriptionUserWithoutSubscription() {
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
				MaxAlarms:            100,
				MinAlarmInterval:     40,
				UnlimitedEmails:      true,
				MaxEmailsPerInterval: util.PositiveIntOrNull(0),
				SlackAlerts:          true,
			},
		},
		nil,
	)

	alarmLimits := suite.service.getAlarmLimits(user)

	// Max alarms should be taken from the team plan
	assert.Equal(suite.T(), uint(100), alarmLimits.maxAlarms)
	// Min alarm interval should be taken from the team plan
	assert.Equal(suite.T(), uint(40), alarmLimits.minAlarmInterval)
	// Unlimited emails should be taken from the team plan
	assert.True(suite.T(), alarmLimits.unlimitedEmails)
	// Max emails per interval should be taken from the team plan
	assert.Equal(suite.T(), uint(0), alarmLimits.maxEmailsPerInterval)
	// Slack alerts should be taken from the team plan
	assert.True(suite.T(), alarmLimits.slackAlerts)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestgetAlarmLimitsTeamWithoutSubscriptionUserWithSubscription() {
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
				MaxAlarms:            10,
				MinAlarmInterval:     50,
				UnlimitedEmails:      true,
				MaxEmailsPerInterval: util.PositiveIntOrNull(30),
				SlackAlerts:          true,
			},
		},
		nil,
	)

	alarmLimits := suite.service.getAlarmLimits(user)

	// Max alarms should be taken from the user plan
	assert.Equal(suite.T(), uint(10), alarmLimits.maxAlarms)
	// Min alarm interval should be taken from the user plan
	assert.Equal(suite.T(), uint(50), alarmLimits.minAlarmInterval)
	// Unlimited emails should be taken from the user plan
	assert.True(suite.T(), alarmLimits.unlimitedEmails)
	// Max emails per interval should be taken from the user plan
	assert.Equal(suite.T(), uint(30), alarmLimits.maxEmailsPerInterval)
	// Slack alerts should be taken from the user plan
	assert.True(suite.T(), alarmLimits.slackAlerts)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) TestGetAlarmLimitsUserWithSubscription() {
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
				MaxAlarms:            20,
				MinAlarmInterval:     40,
				UnlimitedEmails:      false,
				MaxEmailsPerInterval: util.PositiveIntOrNull(20),
				SlackAlerts:          false,
			},
		},
		nil,
	)

	alarmLimits := suite.service.getAlarmLimits(user)

	// Max alarms should be taken from the user plan
	assert.Equal(suite.T(), uint(20), alarmLimits.maxAlarms)
	// Min alarm interval should be taken from the user plan
	assert.Equal(suite.T(), uint(40), alarmLimits.minAlarmInterval)
	// Unlimited emails should be taken from the user plan
	assert.False(suite.T(), alarmLimits.unlimitedEmails)
	// Max emails per interval should be taken from the user plan
	assert.Equal(suite.T(), uint(20), alarmLimits.maxEmailsPerInterval)
	// Slack alerts should be taken from the user plan
	assert.False(suite.T(), alarmLimits.slackAlerts)
}
