package notifications

import (
	"github.com/stretchr/testify/assert"
)

func (suite *NotificationsTestSuite) TestFindEndpointByUserIDAndApplicationARN() {
	var (
		endpoint *Endpoint
		err      error
	)

	// When we try to find an endpoint with a bogus user ID and application ARN
	endpoint, err = suite.service.FindEndpointByUserIDAndApplicationARN(
		12345,
		"bogus",
	)

	// Endpoint object should be nil
	assert.Nil(suite.T(), endpoint)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrEndpointNotFound, err)
	}

	// When we try to find an endpoint with a valid user ID and bogus application ARN
	endpoint, err = suite.service.FindEndpointByUserIDAndApplicationARN(
		suite.users[0].ID,
		"bogus",
	)

	// Endpoint object should be nil
	assert.Nil(suite.T(), endpoint)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrEndpointNotFound, err)
	}

	// When we try to find an endpoint with a bogus user ID and valid application ARN
	endpoint, err = suite.service.FindEndpointByUserIDAndApplicationARN(
		12345,
		suite.cnf.AWS.APNSPlatformApplicationARN,
	)

	// Endpoint object should be nil
	assert.Nil(suite.T(), endpoint)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrEndpointNotFound, err)
	}

	// When we try to find an endpoint with a valid user ID and valid application ARN
	endpoint, err = suite.service.FindEndpointByUserIDAndApplicationARN(
		suite.users[0].ID,
		suite.cnf.AWS.APNSPlatformApplicationARN,
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct endpoint object should be returned with preloaded data
	if assert.NotNil(suite.T(), endpoint) {
		assert.Equal(suite.T(), suite.endpoints[0].ID, endpoint.ID)
		assert.Equal(suite.T(), suite.users[0].ID, uint(endpoint.UserID.Int64))
		assert.Equal(suite.T(), suite.cnf.AWS.APNSPlatformApplicationARN, endpoint.ApplicationARN)
		assert.Equal(suite.T(), "the_arn_1", endpoint.ARN)
		assert.True(suite.T(), endpoint.Enabled)
	}
}
