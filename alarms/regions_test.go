package alarms

import (
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestFindRegionByID() {
	var (
		region *Region
		err    error
	)

	// Let's try to find a region by a bogus ID
	region, err = suite.service.findRegionByID("bogus")

	// Region should be nil
	assert.Nil(suite.T(), region)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrRegionNotFound, err)
	}

	// Now let's pass a valid ID
	region, err = suite.service.findRegionByID("SGP")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct region should be returned with preloaded data
	if assert.NotNil(suite.T(), region) {
		assert.Equal(suite.T(), regions.Singapore, region.ID)
	}
}
