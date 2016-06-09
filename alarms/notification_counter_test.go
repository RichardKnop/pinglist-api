package alarms

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestFindNotificationCounter() {
	var (
		countBefore, countAfter int
		notificationCounter     *NotificationCounter
		err                     error
	)

	// Count before
	suite.db.Model(new(NotificationCounter)).Count(&countBefore)

	// This should create a new notification counter as there isn't one for May yet
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.May),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), int64(suite.users[0].ID), notificationCounter.UserID.Int64)
		assert.Equal(suite.T(), uint(2016), notificationCounter.Year)
		assert.Equal(suite.T(), uint(time.May), notificationCounter.Month)
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}

	// Count after
	suite.db.Model(new(NotificationCounter)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Count before
	suite.db.Model(new(NotificationCounter)).Count(&countBefore)

	// This should just fetch an existing row, as we already have a counter for May
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.May),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), int64(suite.users[0].ID), notificationCounter.UserID.Int64)
		assert.Equal(suite.T(), uint(2016), notificationCounter.Year)
		assert.Equal(suite.T(), uint(time.May), notificationCounter.Month)
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}

	// Count after
	suite.db.Model(new(NotificationCounter)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Count before
	suite.db.Model(new(NotificationCounter)).Count(&countBefore)

	// This should insert a new row as month changed to June from May
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), int64(suite.users[0].ID), notificationCounter.UserID.Int64)
		assert.Equal(suite.T(), uint(2016), notificationCounter.Year)
		assert.Equal(suite.T(), uint(time.June), notificationCounter.Month)
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}

	// Count after
	suite.db.Model(new(NotificationCounter)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
}

func (suite *AlarmsTestSuite) TestUpdateNotificationCounterIncrementEmail() {
	var (
		notificationCounter *NotificationCounter
		err                 error
	)

	// Increment email counter for the first time
	err = suite.service.updateNotificationCounterIncrementEmail(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)

	// Email counter should be at 1
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), uint(1), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}

	// Increment email counter for the second time
	err = suite.service.updateNotificationCounterIncrementEmail(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)

	// Email counter should be at 2
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), uint(2), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}
}

func (suite *AlarmsTestSuite) TestUpdateNotificationCounterIncrementPush() {
	var (
		notificationCounter *NotificationCounter
		err                 error
	)

	// Increment push counter for the first time
	err = suite.service.updateNotificationCounterIncrementPush(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)

	// Push counter should be at 1
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(1), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}

	// Increment push counter for the second time
	err = suite.service.updateNotificationCounterIncrementPush(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)

	// Push counter should be at 2
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(2), notificationCounter.Push)
		assert.Equal(suite.T(), uint(0), notificationCounter.Slack)
	}
}

func (suite *AlarmsTestSuite) TestUpdateNotificationCounterIncrementSlack() {
	var (
		notificationCounter *NotificationCounter
		err                 error
	)

	// Increment slack counter for the first time
	err = suite.service.updateNotificationCounterIncrementSlack(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)

	// Slack counter should be at 1
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(1), notificationCounter.Slack)
	}

	// Increment slack counter for the second time
	err = suite.service.updateNotificationCounterIncrementSlack(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)

	// Slack counter should be at 2
	notificationCounter, err = suite.service.findNotificationCounter(
		suite.users[0].ID,
		uint(2016),
		uint(time.June),
	)
	assert.NoError(suite.T(), err)
	if assert.NotNil(suite.T(), notificationCounter) {
		assert.Equal(suite.T(), uint(0), notificationCounter.Email)
		assert.Equal(suite.T(), uint(0), notificationCounter.Push)
		assert.Equal(suite.T(), uint(2), notificationCounter.Slack)
	}
}
