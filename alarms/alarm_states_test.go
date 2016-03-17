package alarms

import (
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestFindAlarmStateByID() {
	var (
		alarmState *AlarmState
		err        error
	)

	// Let's try to find an alarm state by a bogus ID
	alarmState, err = suite.service.findAlarmStateByID("bogus")

	// Alarm state should be nil
	assert.Nil(suite.T(), alarmState)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAlarmStateNotFound, err)
	}

	// Now let's pass a valid ID
	alarmState, err = suite.service.findAlarmStateByID("ok")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct alarm state should be returned with preloaded data
	if assert.NotNil(suite.T(), alarmState) {
		assert.Equal(suite.T(), alarmstates.OK, alarmState.ID)
	}
}
