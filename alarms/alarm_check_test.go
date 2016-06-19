package alarms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/RichardKnop/pinglist-api/notifications"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckSimpleUseCase() {
	var (
		alarmIDs  []uint
		testAlarm *Alarm
		err       error
		testUser  = suite.users[1]
		watermark time.Time
		interval  = uint(60)
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// First, let's try with no active alarms
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 0 alarm IDs
	assert.Equal(suite.T(), 0, len(alarmIDs))

	// Now insert an active test alarm not yet ready to be checked
	watermark = time.Now().Add(-time.Duration(interval-1) * time.Second)
	testAlarm = &Alarm{
		User:             testUser,
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://foo",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		MaxResponseTime:  1000,
		Interval:         interval,
		Active:           true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test alarm failed")

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 0 alarm IDs
	assert.Equal(suite.T(), 0, len(alarmIDs))

	// Now insert an active test alarm ready for check
	watermark = time.Now().Add(-time.Duration(interval+1) * time.Second)
	testAlarm = &Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://bar",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		MaxResponseTime:  1000,
		Interval:         interval,
		Active:           true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test alarm failed")

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 alarm ID
	assert.Equal(suite.T(), 1, len(alarmIDs))
	assert.Equal(suite.T(), testAlarm.ID, alarmIDs[0])
}

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckMultipleAlarms() {
	var (
		alarmIDs  []uint
		testAlarm *Alarm
		err       error
		testUsers = []*accounts.User{
			suite.users[1],
			suite.users[2],
		}
		watermark        time.Time
		interval         = uint(60)
		expectedAlarmIDs []uint
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// Insert an active test alarm ready for check, one for each test user
	for _, testUser := range testUsers {
		watermark = time.Now().Add(-time.Duration(interval+1) * time.Second)
		testAlarm = &Alarm{
			User:             testUser,
			Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
			AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
			EndpointURL:      "http://bar",
			Watermark:        util.TimeOrNull(&watermark),
			ExpectedHTTPCode: 200,
			MaxResponseTime:  1000,
			Interval:         interval,
			Active:           true,
		}
		err = suite.db.Create(testAlarm).Error
		assert.NoError(suite.T(), err, "Inserting test alarm failed")
		expectedAlarmIDs = append(expectedAlarmIDs, testAlarm.ID)
	}

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 2 alarm IDs
	assert.Equal(suite.T(), 2, len(alarmIDs))
	for i, expectedAlarmID := range expectedAlarmIDs {
		assert.Equal(suite.T(), expectedAlarmID, alarmIDs[i])
	}
}

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckFreeTierMaxAlarmsLimit() {
	var (
		alarmIDs         []uint
		testAlarm        *Alarm
		err              error
		testUser         = suite.users[1]
		watermark        time.Time
		interval         = uint(60)
		expectedAlarmIDs []uint
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// Insert multiple active test alarms ready for check, user in free tier
	for i := 0; i < 2; i++ {
		watermark = time.Now().Add(-time.Duration(interval+10) * time.Second)
		testAlarm = &Alarm{
			User:             testUser,
			Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
			AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
			EndpointURL:      "http://bar",
			Watermark:        util.TimeOrNull(&watermark),
			ExpectedHTTPCode: 200,
			MaxResponseTime:  1000,
			Interval:         interval,
			Active:           true,
		}
		err = suite.db.Create(testAlarm).Error
		assert.NoError(suite.T(), err, "Inserting test alarm failed")
		if i == 0 {
			expectedAlarmIDs = append(expectedAlarmIDs, testAlarm.ID)
		}
	}

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 alarm ID
	assert.Equal(suite.T(), 1, len(alarmIDs))
	assert.Equal(suite.T(), expectedAlarmIDs[0], alarmIDs[0])
}

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckExpiredUserSubscription() {
	var (
		alarmIDs         []uint
		testAlarm        *Alarm
		err              error
		watermark        time.Time
		testUser         = suite.users[1]
		testPlan         *subscriptions.Plan
		testCustomer     *subscriptions.Customer
		periodEnd        time.Time
		testSubscription *subscriptions.Subscription
		interval         = uint(60)
		expectedAlarmIDs []uint
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// Insert an expired subscription
	testPlan = &subscriptions.Plan{
		MaxAlarms: uint(2),
	}
	err = suite.db.Create(testPlan).Error
	assert.NoError(suite.T(), err, "Inserting test plan failed")
	testCustomer = &subscriptions.Customer{
		User: testUser,
	}
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Inserting test customer failed")
	periodEnd = time.Now().Add(-10 * time.Second)
	testSubscription = &subscriptions.Subscription{
		Customer:  testCustomer,
		Plan:      testPlan,
		PeriodEnd: util.TimeOrNull(&periodEnd),
	}
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Inserting test subscription failed")

	// Insert multiple active test alarms ready for check
	for i := 0; i < 2; i++ {
		watermark = time.Now().Add(-time.Duration(interval+10) * time.Second)
		testAlarm = &Alarm{
			User:             testUser,
			Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
			AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
			EndpointURL:      "http://bar",
			Watermark:        util.TimeOrNull(&watermark),
			ExpectedHTTPCode: 200,
			MaxResponseTime:  1000,
			Interval:         interval,
			Active:           true,
		}
		err = suite.db.Create(testAlarm).Error
		assert.NoError(suite.T(), err, "Inserting test alarm failed")
		if i == 0 {
			expectedAlarmIDs = append(expectedAlarmIDs, testAlarm.ID)
		}
	}

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 alarm ID
	assert.Equal(suite.T(), 1, len(alarmIDs))
	assert.Equal(suite.T(), expectedAlarmIDs[0], alarmIDs[0])
}

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckActiveUserSubscription() {
	var (
		alarmIDs         []uint
		testAlarm        *Alarm
		err              error
		watermark        time.Time
		testUser         = suite.users[1]
		testPlan         *subscriptions.Plan
		testCustomer     *subscriptions.Customer
		periodEnd        time.Time
		testSubscription *subscriptions.Subscription
		interval         = uint(60)
		expectedAlarmIDs []uint
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// Insert an active subscription
	testPlan = &subscriptions.Plan{
		PlanID:    "test_plan_id",
		Name:      "Test Plan",
		MaxAlarms: uint(2),
	}
	err = suite.db.Create(testPlan).Error
	assert.NoError(suite.T(), err, "Inserting test plan failed")
	testCustomer = &subscriptions.Customer{
		User:       testUser,
		CustomerID: "test_customer_id",
	}
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Inserting test customer failed")
	periodEnd = time.Now().Add(10 * time.Second)
	testSubscription = &subscriptions.Subscription{
		Customer:       testCustomer,
		Plan:           testPlan,
		PeriodEnd:      util.TimeOrNull(&periodEnd),
		SubscriptionID: "test_subscription_id",
	}
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Inserting test subscription failed")

	// Insert multiple active test alarms ready for check
	for i := 0; i < 2; i++ {
		watermark = time.Now().Add(-time.Duration(interval+10) * time.Second)
		testAlarm = &Alarm{
			User:             testUser,
			Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
			AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
			EndpointURL:      "http://bar",
			Watermark:        util.TimeOrNull(&watermark),
			ExpectedHTTPCode: 200,
			MaxResponseTime:  1000,
			Interval:         interval,
			Active:           true,
		}
		err = suite.db.Create(testAlarm).Error
		assert.NoError(suite.T(), err, "Inserting test alarm failed")
		expectedAlarmIDs = append(expectedAlarmIDs, testAlarm.ID)
	}

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 2 alarm IDs
	assert.Equal(suite.T(), 2, len(alarmIDs))
	for i, expectedAlarmID := range expectedAlarmIDs {
		assert.Equal(suite.T(), expectedAlarmID, alarmIDs[i])
	}
}

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckExpiredTeamSubscription() {
	var (
		alarmIDs         []uint
		testAlarm        *Alarm
		err              error
		watermark        time.Time
		testUser         = suite.users[1]
		testTeamOwner    = suite.users[0]
		testTeam         *teams.Team
		testPlan         *subscriptions.Plan
		testCustomer     *subscriptions.Customer
		periodEnd        time.Time
		testSubscription *subscriptions.Subscription
		interval         = uint(60)
		expectedAlarmIDs []uint
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// Insert a member team with an expired subscription
	testTeam = &teams.Team{
		Owner:   testTeamOwner,
		Name:    "Test Team",
		Members: []*accounts.User{testUser},
	}
	err = suite.db.Create(testTeam).Error
	assert.NoError(suite.T(), err, "Inserting test team failed")
	testPlan = &subscriptions.Plan{
		PlanID:    "test_plan_id",
		Name:      "Test Plan",
		MaxAlarms: uint(2),
	}
	err = suite.db.Create(testPlan).Error
	assert.NoError(suite.T(), err, "Inserting test plan failed")
	testCustomer = &subscriptions.Customer{
		User: testUser,
	}
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Inserting test customer failed")
	periodEnd = time.Now().Add(-10 * time.Second)
	testSubscription = &subscriptions.Subscription{
		Customer:  testCustomer,
		Plan:      testPlan,
		PeriodEnd: util.TimeOrNull(&periodEnd),
	}
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Inserting test subscription failed")

	// Insert multiple active test alarms ready for check
	for i := 0; i < 2; i++ {
		watermark = time.Now().Add(-time.Duration(interval+10) * time.Second)
		testAlarm = &Alarm{
			User:             testUser,
			Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
			AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
			EndpointURL:      "http://bar",
			Watermark:        util.TimeOrNull(&watermark),
			ExpectedHTTPCode: 200,
			MaxResponseTime:  1000,
			Interval:         interval,
			Active:           true,
		}
		err = suite.db.Create(testAlarm).Error
		assert.NoError(suite.T(), err, "Inserting test alarm failed")
		if i == 0 {
			expectedAlarmIDs = append(expectedAlarmIDs, testAlarm.ID)
		}
	}

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 alarm ID
	assert.Equal(suite.T(), 1, len(alarmIDs))
	assert.Equal(suite.T(), expectedAlarmIDs[0], alarmIDs[0])
}

func (suite *AlarmsTestSuite) TestGetAlarmsToCheckActiveTeamSubscription() {
	var (
		alarmIDs         []uint
		testAlarm        *Alarm
		err              error
		watermark        time.Time
		testUser         = suite.users[1]
		testTeamOwner    = suite.users[0]
		testTeam         *teams.Team
		testPlan         *subscriptions.Plan
		testCustomer     *subscriptions.Customer
		periodEnd        time.Time
		testSubscription *subscriptions.Subscription
		interval         = uint(60)
		expectedAlarmIDs []uint
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// Insert a member team with an active subscription
	testTeam = &teams.Team{
		Owner:   testTeamOwner,
		Name:    "Test Team",
		Members: []*accounts.User{testUser},
	}
	err = suite.db.Create(testTeam).Error
	assert.NoError(suite.T(), err, "Inserting test team failed")
	testPlan = &subscriptions.Plan{
		PlanID:    "test_plan_id",
		Name:      "Test Plan",
		MaxAlarms: uint(2),
	}
	err = suite.db.Create(testPlan).Error
	assert.NoError(suite.T(), err, "Inserting test plan failed")
	testCustomer = &subscriptions.Customer{
		User: testUser,
	}
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Inserting test customer failed")
	periodEnd = time.Now().Add(10 * time.Second)
	testSubscription = &subscriptions.Subscription{
		Customer:  testCustomer,
		Plan:      testPlan,
		PeriodEnd: util.TimeOrNull(&periodEnd),
	}
	err = suite.db.Create(testSubscription).Error
	assert.NoError(suite.T(), err, "Inserting test subscription failed")

	// Insert multiple active test alarms ready for check
	for i := 0; i < 2; i++ {
		watermark = time.Now().Add(-time.Duration(interval+10) * time.Second)
		testAlarm = &Alarm{
			User:             testUser,
			Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
			AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
			EndpointURL:      "http://bar",
			Watermark:        util.TimeOrNull(&watermark),
			ExpectedHTTPCode: 200,
			MaxResponseTime:  1000,
			Interval:         interval,
			Active:           true,
		}
		err = suite.db.Create(testAlarm).Error
		assert.NoError(suite.T(), err, "Inserting test alarm failed")
		expectedAlarmIDs = append(expectedAlarmIDs, testAlarm.ID)
	}

	// Try again
	alarmIDs, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 2 alarm IDs
	assert.Equal(suite.T(), 2, len(alarmIDs))
	for i, expectedAlarmID := range expectedAlarmIDs {
		assert.Equal(suite.T(), expectedAlarmID, alarmIDs[i])
	}
}

func (suite *AlarmsTestSuite) TestAlarmCheck() {
	var (
		testAlarm, alarm *Alarm
		err              error
		server           *httptest.Server
		client           *http.Client
		start            time.Time
	)

	// Insert a test alarm
	testAlarm = &Alarm{
		User:                   suite.users[1],
		Region:                 &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:             &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:            "http://foobar",
		ExpectedHTTPCode:       200,
		MaxResponseTime:        1000,
		Interval:               60,
		EmailAlerts:            true,
		PushNotificationAlerts: true,
		SlackAlerts:            true,
		Active:                 true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Fetch the alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// First, let's test a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		util.FormatTime(start),
		util.FormatTime(alarm.Watermark.Time),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 0 incidents
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Second, let's test a slow response
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	assert.NoError(
		suite.T(),
		suite.db.Model(alarm).UpdateColumn("max_response_time", 0).Error,
		"Updating max_response_time to 0 failed",
	)
	suite.mockFindTeamByMemberID(
		alarm.User.ID,
		nil,
		teams.ErrTeamNotFound,
	)
	suite.mockFindActiveSubscriptionByUserID(
		alarm.User.ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				UnlimitedEmails: true,
				SlackAlerts:     true,
			},
		},
		nil,
	)
	suite.mockNewIncidentEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s returned slow response", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	suite.mockNewIncidentSlackMessage(suite.users[1])
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)
	assert.NoError(
		suite.T(),
		suite.db.Model(alarm).UpdateColumn("max_response_time", 1000).Error,
		"Updating max_response_time back to 1000 failed",
	)

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(15 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		util.FormatTime(start),
		util.FormatTime(alarm.Watermark.Time),
	)

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.Slow, alarm.Incidents[0].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(200), alarm.Incidents[0].HTTPCode.Int64)
	assert.True(suite.T(), alarm.Incidents[0].ResponseTime.Valid)
	assert.True(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Third, let's test a timeout
	server, client = testServerTimeout()
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockFindTeamByMemberID(
		alarm.User.ID,
		nil,
		teams.ErrTeamNotFound,
	)
	suite.mockFindActiveSubscriptionByUserID(
		alarm.User.ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				UnlimitedEmails: true,
				SlackAlerts:     true,
			},
		},
		nil,
	)
	suite.mockNewIncidentEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s timed out", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	suite.mockNewIncidentSlackMessage(suite.users[1])
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(15 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		util.FormatTime(start),
		util.FormatTime(alarm.Watermark.Time),
	)

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[1].IncidentTypeID.String)
	assert.False(suite.T(), alarm.Incidents[1].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ResponseTime.Valid)
	assert.False(suite.T(), alarm.Incidents[1].Response.Valid)
	assert.True(suite.T(), alarm.Incidents[1].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Next, let's test a bad code
	server, client = testServer(&http.Response{StatusCode: 500})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockFindTeamByMemberID(
		alarm.User.ID,
		nil,
		teams.ErrTeamNotFound,
	)
	suite.mockFindActiveSubscriptionByUserID(
		alarm.User.ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				UnlimitedEmails: true,
				SlackAlerts:     true,
			},
		},
		nil,
	)
	suite.mockNewIncidentEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s returned bad status code", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	suite.mockNewIncidentSlackMessage(suite.users[1])
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(15 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		util.FormatTime(start),
		util.FormatTime(alarm.Watermark.Time),
	)

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[2].HTTPCode.Int64)
	assert.True(suite.T(), alarm.Incidents[2].ResponseTime.Valid)
	assert.True(suite.T(), alarm.Incidents[2].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)

	// Finally, let's test a return to a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockFindTeamByMemberID(
		alarm.User.ID,
		nil,
		teams.ErrTeamNotFound,
	)
	suite.mockFindActiveSubscriptionByUserID(
		alarm.User.ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				UnlimitedEmails: true,
				SlackAlerts:     true,
			},
		},
		nil,
	)
	suite.mockIncidentsResolvedEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s is up and working correctly", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	suite.mockIncidentsResolvedSlackMessage(suite.users[1])
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Sleep for the email & push notification goroutines to finish
	time.Sleep(15 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		util.FormatTime(start),
		util.FormatTime(alarm.Watermark.Time),
	)

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// Resolved incidents
	for _, incident := range alarm.Incidents {
		assert.True(suite.T(), incident.ResolvedAt.Valid)
	}
}

func (suite *AlarmsTestSuite) TestAlarmCheckIdempotency() {
	var (
		testAlarm, alarm *Alarm
		err              error
		server           *httptest.Server
		client           *http.Client
		start            time.Time
	)

	// Insert a test alarm
	testAlarm = &Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://foobar",
		ExpectedHTTPCode: 200,
		MaxResponseTime:  1000,
		Interval:         60,
		EmailAlerts:      true,
		Active:           true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare test server and client
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}

	// Just one response time metric will be logged
	suite.mockLogResponseTime(start, testAlarm.ID, nil)

	// Trigger multiple parallel alarm checks with the same watermark
	var (
		concurrency = 4
		errChan     = make(chan error)
	)
	for i := 0; i < concurrency; i++ {
		go suite.alarmCheckWrapper(testAlarm.ID, testAlarm.Watermark.Time, errChan)
		time.Sleep(15 * time.Millisecond)
	}

	// Receive start times and errors from goroutines
	var (
		successful int
		errs       []error
	)
	for i := 0; i < concurrency; i++ {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		} else {
			successful++
		}
	}

	// One alarm check should have gone through
	assert.Equal(suite.T(), 1, successful)

	// 3 alarm checks should have been stopped by the idempotency check
	assert.Equal(suite.T(), 3, len(errs))
	for _, err := range errs {
		assert.Equal(suite.T(), ErrCheckAlreadyTriggered, err)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		util.FormatTime(start),
		util.FormatTime(alarm.Watermark.Time),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 0 incidents
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(15 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) alarmCheckWrapper(alarmID uint, watermark time.Time, errChan chan error) {
	errChan <- suite.service.CheckAlarm(alarmID, watermark)
}
