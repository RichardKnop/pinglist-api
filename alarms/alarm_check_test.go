package alarms

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetAlarmsToCheck() {
	var (
		alarms    []*Alarm
		err       error
		watermark time.Time
		interval  = uint(60)
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// First, let's try with no active alarms
	alarms, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 0 alarms
	assert.Equal(suite.T(), 0, len(alarms))

	// Now insert an active test alarm with watermark + interval >= now
	watermark = time.Now().Add(-time.Duration(interval-1) * time.Second)
	err = suite.db.Create(&Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://foo",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		Interval:         interval,
		Active:           true,
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Try again
	alarms, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 0 alarms
	assert.Equal(suite.T(), 0, len(alarms))

	// Now insert an active test alarm with watermark + interval < now
	watermark = time.Now().Add(-time.Duration(interval+1) * time.Second)
	err = suite.db.Create(&Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://bar",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		Interval:         interval,
		Active:           true,
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Try again
	alarms, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 alarms
	assert.Equal(suite.T(), 1, len(alarms))
	assert.Equal(suite.T(), "http://bar", alarms[0].EndpointURL)
}

func (suite *AlarmsTestSuite) TestAlarmCheck() {
	var (
		alarm  = suite.alarms[2]
		err    error
		server *httptest.Server
		client *http.Client
		start  time.Time
	)

	// First, let's test a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockLogRequestTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 0 incidents
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

	// Second, let's test a timeout
	server, client = testServerTimeout()
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockAlarmDownEmail()
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Sleep for the email goroutine to finish
	time.Sleep(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Incidents[0].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].IncidentTypeID.String)
	assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Third, let's test a bad code
	server, client = testServer(&http.Response{StatusCode: 500})
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(suite.alarms[2].ID, alarm.Watermark.Time)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Incidents[1].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
	assert.Equal(suite.T(), "", alarm.Incidents[1].Response.String)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

	// Finally, let's test a return to a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockAlarmUpEmail()
	suite.mockLogRequestTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(suite.alarms[2].ID, alarm.Watermark.Time)

	// Sleep for the email goroutine to finish
	time.Sleep(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// Resolved incidents
	for _, incident := range alarm.Incidents {
		assert.True(suite.T(), incident.ResolvedAt.Valid)
	}

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())
}

func (suite *AlarmsTestSuite) TestAlarmCheckIdempotency() {
	// Fetch the alarm
	alarm := new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Prepare test server and client
	server, client := testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start := time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}

	// Just one request time metric will be logged
	suite.mockLogRequestTime(start, alarm.ID, nil)

	concurrency := 4

	// Trigger multiple parallel alarm checks with the same watermark
	errChan := make(chan error)
	for i := 0; i < concurrency; i++ {
		go suite.alarmCheckWrapper(alarm.ID, alarm.Watermark.Time, errChan)
		time.Sleep(10 * time.Millisecond)
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
		First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 0 incidents
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())
}

func (suite *AlarmsTestSuite) alarmCheckWrapper(alarmID uint, watermark time.Time, errChan chan error) {
	errChan <- suite.service.CheckAlarm(alarmID, watermark)
}
