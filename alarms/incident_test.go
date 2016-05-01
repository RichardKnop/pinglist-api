package alarms

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/RichardKnop/pinglist-api/notifications"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestIncidents() {
	var (
		testAlarm, alarm *Alarm
		err              error
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
		Active:                 true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Fetch the alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// First, let's open a new slow_response incident
	when1 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when1
	}
	suite.mockAlarmDownEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s is down", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	err = suite.service.openIncident(
		alarm,
		incidenttypes.SlowResponse,
		&http.Response{StatusCode: 200},
		2345, // response time
		"",   // error message
	)

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(10 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state changed, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status changed to Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt should be set
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt should be nil
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 1 incident
		assert.Equal(suite.T(), 1, len(alarm.Incidents))

		// New incident
		assert.Equal(suite.T(), incidenttypes.SlowResponse, alarm.Incidents[0].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(200), alarm.Incidents[0].HTTPCode.Int64)
		assert.Equal(suite.T(), int64(2345), alarm.Incidents[0].ResponseTime.Int64)
		assert.True(suite.T(), alarm.Incidents[0].Response.Valid)
		assert.False(suite.T(), alarm.Incidents[0].ErrorMessage.Valid)
		assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.SlowResponse, alarm.Incidents[0].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(200), alarm.Incidents[0].HTTPCode.Int64)
	assert.Equal(suite.T(), int64(2345), alarm.Incidents[0].ResponseTime.Int64)
	assert.True(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Second, let's try opening another slow_response incident
	// This should not create a new incident entry
	when2 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when2
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.SlowResponse,
		&http.Response{StatusCode: 200},
		3456, // response time
		"",   // error message
	)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state unchanged, no new incidents created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 1 incident
		assert.Equal(suite.T(), 1, len(alarm.Incidents))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// Third, let's open a new timeout incident
	when3 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when3
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.Timeout,
		nil,                // response
		0,                  // response time
		"timeout error...", // error message
	)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state changed, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status changed to Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt should be set
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt should be nil
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 2 incidents
		assert.Equal(suite.T(), 2, len(alarm.Incidents))

		// New incident
		assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[1].IncidentTypeID.String)
		assert.False(suite.T(), alarm.Incidents[1].HTTPCode.Valid)
		assert.False(suite.T(), alarm.Incidents[1].ResponseTime.Valid)
		assert.False(suite.T(), alarm.Incidents[1].Response.Valid)
		assert.Equal(suite.T(), "timeout error...", alarm.Incidents[1].ErrorMessage.String)
		assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[1].IncidentTypeID.String)
	assert.False(suite.T(), alarm.Incidents[1].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ResponseTime.Valid)
	assert.False(suite.T(), alarm.Incidents[1].Response.Valid)
	assert.Equal(suite.T(), "timeout error...", alarm.Incidents[1].ErrorMessage.String)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Next, let's try opening another timeout incident
	// This should not create a new incident entry
	when4 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when4
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.Timeout,
		nil,                // response
		0,                  // response time
		"timeout error...", // error message
	)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state unchanged, no new incidents created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 2 incidents
		assert.Equal(suite.T(), 2, len(alarm.Incidents))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// Next, open a new bad code incident
	when5 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when5
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.BadCode,
		&http.Response{StatusCode: 500},
		1000, // response time
		"",   // error message
	)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state unchanged, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 3 incidents
		assert.Equal(suite.T(), 3, len(alarm.Incidents))

		// New incident
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(500), alarm.Incidents[2].HTTPCode.Int64)
		assert.Equal(suite.T(), int64(1000), alarm.Incidents[2].ResponseTime.Int64)
		assert.True(suite.T(), alarm.Incidents[2].Response.Valid)
		assert.False(suite.T(), alarm.Incidents[2].ErrorMessage.Valid)
		assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[2].HTTPCode.Int64)
	assert.Equal(suite.T(), int64(1000), alarm.Incidents[2].ResponseTime.Int64)
	assert.True(suite.T(), alarm.Incidents[2].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)

	// Next, let's try opening another bad code incident with the same code
	// This should not create a new incident entry
	when6 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when6
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.BadCode,
		&http.Response{StatusCode: 500},
		900, // response time
		"",  // error message
	)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state unchanged, no new incidents created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 3 incidents
		assert.Equal(suite.T(), 3, len(alarm.Incidents))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// Next, open a new bad code incident with a different code
	when7 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when7
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.BadCode,
		&http.Response{StatusCode: 404},
		1000, // response time
		"",   // error message
	)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state unchanged, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 4 incidents
		assert.Equal(suite.T(), 4, len(alarm.Incidents))

		// New incident
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[3].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(404), alarm.Incidents[3].HTTPCode.Int64)
		assert.Equal(suite.T(), int64(1000), alarm.Incidents[3].ResponseTime.Int64)
		assert.True(suite.T(), alarm.Incidents[3].Response.Valid)
		assert.False(suite.T(), alarm.Incidents[3].ErrorMessage.Valid)
		assert.False(suite.T(), alarm.Incidents[3].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 4 incidents
	assert.Equal(suite.T(), 4, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[3].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(404), alarm.Incidents[3].HTTPCode.Int64)
	assert.Equal(suite.T(), int64(1000), alarm.Incidents[3].ResponseTime.Int64)
	assert.True(suite.T(), alarm.Incidents[3].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[3].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[3].ResolvedAt.Valid)

	// Finally, resolve the incidents
	when8 := time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return when8
	}
	suite.mockAlarmUpEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s is up again", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	err = suite.service.resolveIncidents(alarm)

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(10 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil, the alarm state changed, all incidents resolved
	if assert.Nil(suite.T(), err) {
		// Status back to OK
		assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt should be set
		assert.Equal(
			suite.T(),
			when8.Format("2006-01-02T15:04:05Z"),
			alarm.LastUptimeStartedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		)

		// 4 incidents
		assert.Equal(suite.T(), 4, len(alarm.Incidents))

		// Resolved incidents
		for _, incident := range alarm.Incidents {
			assert.True(suite.T(), incident.ResolvedAt.Valid)
		}
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 4 incidents
	assert.Equal(suite.T(), 4, len(alarm.Incidents))

	// Resolved incidents
	for _, incident := range alarm.Incidents {
		assert.True(suite.T(), incident.ResolvedAt.Valid)
	}
}

func (suite *AlarmsTestSuite) TestIncidentTypeCounts() {
	incidentTypeCounts, err := suite.service.incidentTypeCounts(
		nil, // user
		nil, // alarm
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, incidentTypeCounts[incidenttypes.SlowResponse])
		assert.Equal(suite.T(), 2, incidentTypeCounts[incidenttypes.Timeout])
		assert.Equal(suite.T(), 1, incidentTypeCounts[incidenttypes.BadCode])
		assert.Equal(suite.T(), 1, incidentTypeCounts[incidenttypes.Other])
	}
}

func (suite *AlarmsTestSuite) TestIncidentsCount() {
	var (
		err   error
		count int
	)

	// Without any filtering
	count, err = suite.service.incidentsCount(
		nil, // user
		nil, // alarm
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user with 4 incidents
	count, err = suite.service.incidentsCount(
		suite.users[1],
		nil, // alarm
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by alarm with 4 incidents
	count, err = suite.service.incidentsCount(
		nil, // user
		suite.alarms[0],
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user and alarm with 4 incidents
	count, err = suite.service.incidentsCount(
		suite.users[1],
		suite.alarms[0],
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user without incidents
	count, err = suite.service.incidentsCount(
		suite.users[0],
		nil, // alarm
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	// Filter by alarm without incidents
	count, err = suite.service.incidentsCount(
		nil, // user
		suite.alarms[1],
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	// Filter by user without incidents and alarm with incidents
	count, err = suite.service.incidentsCount(
		suite.users[1],
		suite.alarms[1],
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	// Filter by user with incidents and alarm without incidents
	count, err = suite.service.incidentsCount(
		suite.users[0],
		suite.alarms[0],
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	// Filter by incident type with 2 incidents
	timeoutIncidentType := incidenttypes.Timeout
	count, err = suite.service.incidentsCount(
		nil,                  // user
		nil,                  // alarm
		&timeoutIncidentType, // incident type
		nil,                  // from
		nil,                  // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	// Filter by incident type with 0 incidents
	slowIncidentType := incidenttypes.SlowResponse
	count, err = suite.service.incidentsCount(
		nil,               // user
		nil,               // alarm
		&slowIncidentType, // incident type
		nil,               // from
		nil,               // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	// Filter by >= from
	from, err := util.ParseTimestamp("2016-02-27T22:31:12Z")
	assert.NoError(suite.T(), err, "Failed parsing from timestamp")
	count, err = suite.service.incidentsCount(
		nil, // user
		nil, // alarm
		nil, // incident type
		&from,
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, count)
	}

	// Filter by <= to
	to, err := util.ParseTimestamp("2016-02-27T22:32:12Z")
	assert.NoError(suite.T(), err, "Failed parsing to timestamp")
	count, err = suite.service.incidentsCount(
		nil, // user
		nil, // alarm
		nil, // incident type
		nil, // from
		&to,
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, count)
	}

	// Filter by from >= to
	count, err = suite.service.incidentsCount(
		nil, // user
		nil, // alarm
		nil, // incident type
		&from,
		&to,
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}
}

func (suite *AlarmsTestSuite) TestFindPaginatedIncidents() {
	var (
		incidents []*Incident
		err       error
	)

	// This should return all incidents
	incidents, err = suite.service.findPaginatedIncidents(
		0,   // offset
		25,  // limit
		"",  // order by
		nil, // user
		nil, // alarm
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(incidents))
		assert.Equal(suite.T(), suite.incidents[0].ID, incidents[0].ID)
		assert.Equal(suite.T(), suite.incidents[1].ID, incidents[1].ID)
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[2].ID)
		assert.Equal(suite.T(), suite.incidents[3].ID, incidents[3].ID)
	}

	// This should return all incidents ordered by ID desc
	incidents, err = suite.service.findPaginatedIncidents(
		0,         // offset
		25,        // limit
		"id desc", // order by
		nil,       // user
		nil,       // alarm
		nil,       // incident type
		nil,       // from
		nil,       // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(incidents))
		assert.Equal(suite.T(), suite.incidents[3].ID, incidents[0].ID)
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[1].ID)
		assert.Equal(suite.T(), suite.incidents[1].ID, incidents[2].ID)
		assert.Equal(suite.T(), suite.incidents[0].ID, incidents[3].ID)
	}

	// Test offset
	incidents, err = suite.service.findPaginatedIncidents(
		2,   // offset
		25,  // limit
		"",  // order by
		nil, // user
		nil, // alarm
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(incidents))
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[0].ID)
		assert.Equal(suite.T(), suite.incidents[3].ID, incidents[1].ID)
	}

	// Test limit
	incidents, err = suite.service.findPaginatedIncidents(
		2,   // offset
		1,   // limit
		"",  // order by
		nil, // user
		nil, // alarm
		nil, // incident type
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(incidents))
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[0].ID)
	}
}
