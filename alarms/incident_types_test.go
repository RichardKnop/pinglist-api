package alarms

import (
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestFindIncidentTypeByID() {
	var (
		incidentType *IncidentType
		err          error
	)

	// Let's try to find an incident type by a bogus ID
	incidentType, err = suite.service.findIncidentTypeByID("bogus")

	// Incident type should be nil
	assert.Nil(suite.T(), incidentType)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrIncidentTypeNotFound, err)
	}

	// Now let's pass a valid ID
	incidentType, err = suite.service.findIncidentTypeByID("timeout")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct incident type should be returned with preloaded data
	if assert.NotNil(suite.T(), incidentType) {
		assert.Equal(suite.T(), incidenttypes.Timeout, incidentType.ID)
	}
}
