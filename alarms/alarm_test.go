package alarms

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/stretchr/testify/assert"
)

func TestHasOpenIncident(t *testing.T) {
	alarm := &Alarm{
		Incidents: []*Incident{
			&Incident{
				IncidentTypeID: util.StringOrNull(incidenttypes.SlowResponse),
			},
			&Incident{
				IncidentTypeID: util.StringOrNull(incidenttypes.Timeout),
				ErrorMessage:   util.StringOrNull("timeout error..."),
			},
			&Incident{
				IncidentTypeID: util.StringOrNull(incidenttypes.Other),
				ErrorMessage:   util.StringOrNull("other error..."),
			},
			&Incident{
				IncidentTypeID: util.StringOrNull(incidenttypes.BadCode),
				HTTPCode:       sql.NullInt64{Valid: true, Int64: 500},
			},
		},
	}

	assert.True(t, alarm.HasOpenIncident(
		incidenttypes.SlowResponse,
		nil, // response
		"",  // error message
	))

	assert.False(t, alarm.HasOpenIncident(
		incidenttypes.Timeout,
		nil, // response
		"",  // error message
	))
	assert.True(t, alarm.HasOpenIncident(
		incidenttypes.Timeout,
		nil,                // response
		"timeout error...", // error message
	))

	assert.False(t, alarm.HasOpenIncident(
		incidenttypes.Other,
		nil, // response
		"",  // error message
	))
	assert.True(t, alarm.HasOpenIncident(
		incidenttypes.Other,
		nil,              // response
		"other error...", // error message
	))

	assert.True(t, alarm.HasOpenIncident(
		incidenttypes.BadCode,
		&http.Response{StatusCode: 500}, // response
		"", // error message
	))
	assert.False(t, alarm.HasOpenIncident(
		incidenttypes.BadCode,
		&http.Response{StatusCode: 404}, // response
		"", // error message
	))
}

func (suite *AlarmsTestSuite) TestFindAlarmById() {
	var (
		alarm *Alarm
		err   error
	)

	// When we try to find an alarm with a bogus ID
	alarm, err = suite.service.FindAlarmByID(12345)

	// Alarm object should be nil
	assert.Nil(suite.T(), alarm)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAlarmNotFound, err)
	}

	// When we try to find an alarm with a valid ID
	alarm, err = suite.service.FindAlarmByID(suite.alarms[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct alarm object should be returned with preloaded data
	if assert.NotNil(suite.T(), alarm) {
		assert.Equal(suite.T(), suite.alarms[0].ID, alarm.ID)
		assert.Equal(suite.T(), suite.users[1].ID, alarm.User.ID)
		assert.Equal(suite.T(), "test@user", alarm.User.OauthUser.Username)

		// Only open incidents should be preloaded
		assert.Equal(suite.T(), 3, len(alarm.Incidents))

		// Timeout incident
		assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].IncidentTypeID.String)
		assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
		assert.False(suite.T(), alarm.Incidents[0].Response.Valid)

		// Bad code incident
		assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].IncidentTypeID.String)
		assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
		assert.Equal(suite.T(), "Internal Server Error", alarm.Incidents[1].Response.String)

		// Other incident
		assert.Equal(suite.T(), incidenttypes.Other, alarm.Incidents[2].IncidentTypeID.String)
		assert.False(suite.T(), alarm.Incidents[2].HTTPCode.Valid)
		assert.False(suite.T(), alarm.Incidents[2].Response.Valid)
	}
}

func (suite *AlarmsTestSuite) TestAlarmsCount() {
	var (
		count int
		err   error
	)

	// Without filtering
	count, err = suite.service.alarmsCount(nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user with 4 alarms
	count, err = suite.service.alarmsCount(suite.users[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user without alarms
	count, err = suite.service.alarmsCount(suite.users[2])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *AlarmsTestSuite) TestFindPaginatedAlarms() {
	var (
		alarms []*Alarm
		err    error
	)

	// This should return all alarms
	alarms, err = suite.service.findPaginatedAlarms(0, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(alarms))
		assert.Equal(suite.T(), suite.alarms[0].ID, alarms[0].ID)
		assert.Equal(suite.T(), suite.alarms[1].ID, alarms[1].ID)
		assert.Equal(suite.T(), suite.alarms[2].ID, alarms[2].ID)
		assert.Equal(suite.T(), suite.alarms[3].ID, alarms[3].ID)
	}

	// This should return all alarms ordered by ID desc
	alarms, err = suite.service.findPaginatedAlarms(0, 25, "id desc", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(alarms))
		assert.Equal(suite.T(), suite.alarms[3].ID, alarms[0].ID)
		assert.Equal(suite.T(), suite.alarms[2].ID, alarms[1].ID)
		assert.Equal(suite.T(), suite.alarms[1].ID, alarms[2].ID)
		assert.Equal(suite.T(), suite.alarms[0].ID, alarms[3].ID)
	}

	// Test offset
	alarms, err = suite.service.findPaginatedAlarms(2, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(alarms))
		assert.Equal(suite.T(), suite.alarms[2].ID, alarms[0].ID)
		assert.Equal(suite.T(), suite.alarms[3].ID, alarms[1].ID)
	}

	// Test limit
	alarms, err = suite.service.findPaginatedAlarms(2, 1, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(alarms))
		assert.Equal(suite.T(), suite.alarms[2].ID, alarms[0].ID)
	}
}
