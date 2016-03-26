package alarms

import (
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestIncidents() {
	var (
		alarm *Alarm
		err   error
	)

	// Fetch the test alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// First, let's open a new timeout incident
	suite.mockAlarmDownEmail()
	when1 := time.Now()
	gorm.NowFunc = func() time.Time {
		return when1
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.Timeout,
		nil,                // HTTP response
		"timeout error...", // error message
	)

	// Sleep for the email goroutine to finish
	time.Sleep(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

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
		assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[0].AlarmID.Int64))
		assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].IncidentTypeID.String)
		assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
		assert.False(suite.T(), alarm.Incidents[0].Response.Valid)
		assert.Equal(suite.T(), "timeout error...", alarm.Incidents[0].ErrorMessage.String)
		assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[0].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].IncidentTypeID.String)
	assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.Equal(suite.T(), "timeout error...", alarm.Incidents[0].ErrorMessage.String)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Second, let's try opening another timeout incident
	// This should not create a new incident entry
	when2 := time.Now()
	gorm.NowFunc = func() time.Time {
		return when2
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.Timeout,
		nil,                // HTTP response
		"timeout error...", // error message
	)

	// Error should be nil, the alarm state unchanged, no new incidents created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 1 incident
		assert.Equal(suite.T(), 1, len(alarm.Incidents))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// Third, open a new bad code incident
	when3 := time.Now()
	gorm.NowFunc = func() time.Time {
		return when3
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.BadCode,
		&http.Response{StatusCode: 500},
		"", // error message
	)

	// Error should be nil, the alarm state unchanged, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 2 incidents
		assert.Equal(suite.T(), 2, len(alarm.Incidents))

		// New incident
		assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[1].AlarmID.Int64))
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
		assert.False(suite.T(), alarm.Incidents[1].Response.Valid)
		assert.False(suite.T(), alarm.Incidents[1].ErrorMessage.Valid)
		assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[1].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
	assert.False(suite.T(), alarm.Incidents[1].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Next, let's try opening another bad code incident with the same code
	// This should not create a new incident entry
	when4 := time.Now()
	gorm.NowFunc = func() time.Time {
		return when4
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.BadCode,
		&http.Response{StatusCode: 500},
		"", // error message
	)

	// Error should be nil, the alarm state unchanged, no new incidents created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 2 incidents
		assert.Equal(suite.T(), 2, len(alarm.Incidents))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// Next, open a new bad code incident with a different code
	when5 := time.Now()
	gorm.NowFunc = func() time.Time {
		return when5
	}
	err = suite.service.openIncident(
		alarm,
		incidenttypes.BadCode,
		&http.Response{StatusCode: 404},
		"", // error message
	)

	// Error should be nil, the alarm state unchanged, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)
		// LastDowntimeStartedAt unchanged
		assert.Equal(
			suite.T(),
			when1.Format("2006-01-02T15:04:05Z"),
			alarm.LastDowntimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt still nill
		assert.False(suite.T(), alarm.LastUptimeStartedAt.Valid)

		// 3 incidents
		assert.Equal(suite.T(), 3, len(alarm.Incidents))

		// New incident
		assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[2].AlarmID.Int64))
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(404), alarm.Incidents[2].HTTPCode.Int64)
		assert.False(suite.T(), alarm.Incidents[2].Response.Valid)
		assert.False(suite.T(), alarm.Incidents[2].ErrorMessage.Valid)
		assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[2].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(404), alarm.Incidents[2].HTTPCode.Int64)
	assert.False(suite.T(), alarm.Incidents[2].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)

	// Finally, resolve the incidents
	when6 := time.Now()
	gorm.NowFunc = func() time.Time {
		return when6
	}
	suite.mockAlarmUpEmail()
	err = suite.service.resolveIncidents(alarm)

	// Sleep for the email goroutine to finish
	time.Sleep(5 * time.Millisecond)

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
			alarm.LastDowntimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)
		// LastUptimeStartedAt should be set
		assert.Equal(
			suite.T(),
			when6.Format("2006-01-02T15:04:05Z"),
			alarm.LastUptimeStartedAt.Time.Format("2006-01-02T15:04:05Z"),
		)

		// 3 incidents
		assert.Equal(suite.T(), 3, len(alarm.Incidents))

		// Resolved incidents
		for _, incident := range alarm.Incidents {
			assert.True(suite.T(), incident.ResolvedAt.Valid)
		}
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// Resolved incidents
	for _, incident := range alarm.Incidents {
		assert.True(suite.T(), incident.ResolvedAt.Valid)
	}
}

func (suite *AlarmsTestSuite) TestPaginatedIncidentsCount() {
	var (
		err   error
		count int
	)

	count, err = suite.service.paginatedIncidentsCount(nil, nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedIncidentsCount(suite.users[1], nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedIncidentsCount(nil, suite.alarms[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedIncidentsCount(suite.users[1], suite.alarms[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedIncidentsCount(suite.users[0], nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	count, err = suite.service.paginatedIncidentsCount(nil, suite.alarms[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	count, err = suite.service.paginatedIncidentsCount(suite.users[1], suite.alarms[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	count, err = suite.service.paginatedIncidentsCount(suite.users[0], suite.alarms[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *AlarmsTestSuite) TestFindPaginatedIncidents() {
	var (
		incidents []*Incident
		err       error
	)

	// This should return all incidents
	incidents, err = suite.service.findPaginatedIncidents(0, 25, "", nil, nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(incidents))
		assert.Equal(suite.T(), suite.incidents[0].ID, incidents[0].ID)
		assert.Equal(suite.T(), suite.incidents[1].ID, incidents[1].ID)
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[2].ID)
		assert.Equal(suite.T(), suite.incidents[3].ID, incidents[3].ID)
	}

	// This should return all incidents ordered by ID desc
	incidents, err = suite.service.findPaginatedIncidents(0, 25, "id desc", nil, nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(incidents))
		assert.Equal(suite.T(), suite.incidents[3].ID, incidents[0].ID)
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[1].ID)
		assert.Equal(suite.T(), suite.incidents[1].ID, incidents[2].ID)
		assert.Equal(suite.T(), suite.incidents[0].ID, incidents[3].ID)
	}

	// Test offset
	incidents, err = suite.service.findPaginatedIncidents(2, 25, "", nil, nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(incidents))
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[0].ID)
		assert.Equal(suite.T(), suite.incidents[3].ID, incidents[1].ID)
	}

	// Test limit
	incidents, err = suite.service.findPaginatedIncidents(2, 1, "", nil, nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(incidents))
		assert.Equal(suite.T(), suite.incidents[2].ID, incidents[0].ID)
	}
}
