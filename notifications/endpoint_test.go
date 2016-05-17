package notifications

import (
	"github.com/stretchr/testify/assert"
)

func (suite *NotificationsTestSuite) TestFindEndpointByUserIDAndApplicationARN() {
	var (
		testEndpoint, endpoint *Endpoint
		err                    error
	)

	// Insert a test endpoint
	testEndpoint = NewEndpoint(
		suite.users[0],
		suite.cnf.AWS.APNSPlatformApplicationARN,
		"endpoint_arn",
		"device_token",
		false, // enabled
	)
	err = suite.db.Create(testEndpoint).Error
	assert.NoError(suite.T(), err, "Failed to insert a test endpoint")

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
		assert.Equal(suite.T(), testEndpoint.ID, endpoint.ID)
		assert.Equal(suite.T(), testEndpoint.User.ID, endpoint.User.ID)
		assert.Equal(suite.T(), suite.cnf.AWS.APNSPlatformApplicationARN, endpoint.ApplicationARN)
		assert.Equal(suite.T(), testEndpoint.ARN, endpoint.ARN)
		assert.Equal(suite.T(), testEndpoint.Enabled, endpoint.Enabled)
	}
}
