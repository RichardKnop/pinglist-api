package alarms

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
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
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// First, let's open a new timeout incident
	err = suite.service.openIncident(
		alarm,
		incidenttypes.Timeout,
		nil, // HTTP response
		"",  // error message
	)

	// Error should be nil, the alarm state changed, a new incident created
	if assert.Nil(suite.T(), err) {
		// Status changed to Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

		// 1 incident, 0 results
		assert.Equal(suite.T(), 1, len(alarm.Incidents))
		assert.Equal(suite.T(), 0, len(alarm.Results))

		// New incident
		assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[0].AlarmID.Int64))
		assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].IncidentTypeID.String)
		assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
		assert.False(suite.T(), alarm.Incidents[0].Response.Valid)
		assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident, 0 results
	assert.Equal(suite.T(), 1, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

	// New incident
	assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[0].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].IncidentTypeID.String)
	assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Second, let's try opening another timeout incident
	// This should not create a new incident entry
	err = suite.service.openIncident(
		alarm,
		incidenttypes.Timeout,
		nil, // HTTP response
		"",  // error message
	)

	// Error should be nil, the alarm state unchanged, no new incidents created
	if assert.Nil(suite.T(), err) {
		// Status still Alarm
		assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

		// 1 incident, 0 results
		assert.Equal(suite.T(), 1, len(alarm.Incidents))
		assert.Equal(suite.T(), 0, len(alarm.Results))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident, 0 results
	assert.Equal(suite.T(), 1, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

	// Third, open a new bad code incident
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

		// 2 incidents, 0 results
		assert.Equal(suite.T(), 2, len(alarm.Incidents))
		assert.Equal(suite.T(), 0, len(alarm.Results))

		// New incident
		assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[1].AlarmID.Int64))
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
		assert.Equal(suite.T(), "", alarm.Incidents[1].Response.String)
		assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents, 0 results
	assert.Equal(suite.T(), 2, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

	// New incident
	assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[1].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
	assert.Equal(suite.T(), "", alarm.Incidents[1].Response.String)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Next, let's try opening another bad code incident with the same code
	// This should not create a new incident entry
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

		// 2 incidents, 0 results
		assert.Equal(suite.T(), 2, len(alarm.Incidents))
		assert.Equal(suite.T(), 0, len(alarm.Results))
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents, 0 results
	assert.Equal(suite.T(), 2, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

	// Next, open a new bad code incident with a different code
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

		// 3 incidents, 0 results
		assert.Equal(suite.T(), 3, len(alarm.Incidents))
		assert.Equal(suite.T(), 0, len(alarm.Results))

		// New incident
		assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[2].AlarmID.Int64))
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(404), alarm.Incidents[2].HTTPCode.Int64)
		assert.Equal(suite.T(), "", alarm.Incidents[2].Response.String)
		assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 3 incidents, 0 results
	assert.Equal(suite.T(), 3, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

	// New incident
	assert.Equal(suite.T(), suite.alarms[1].ID, uint(alarm.Incidents[2].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(404), alarm.Incidents[2].HTTPCode.Int64)
	assert.Equal(suite.T(), "", alarm.Incidents[2].Response.String)
	assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)

	// Finally, resolve the incidents
	err = suite.service.resolveIncidentsTx(suite.db, alarm)

	// Error should be nil, the alarm state changed, all incidents resolved
	if assert.Nil(suite.T(), err) {
		// Status back to OK
		assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

		// 3 incidents, 0 results
		assert.Equal(suite.T(), 3, len(alarm.Incidents))
		assert.Equal(suite.T(), 0, len(alarm.Results))

		// Resolved incidents
		for _, incident := range alarm.Incidents {
			assert.True(suite.T(), incident.ResolvedAt.Valid)
		}
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[1].ID).RecordNotFound())

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 3 incidents, 0 results
	assert.Equal(suite.T(), 3, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

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
