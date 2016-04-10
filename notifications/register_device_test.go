package notifications

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *NotificationsTestSuite) TestRegisterDeviceRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.registerDeviceHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *NotificationsTestSuite) TestRegisterIOSDeviceFirstTime() {
	// Prepare a request
	payload, err := json.Marshal(&DeviceRequest{
		Platform: PlatformIOS,
		Token:    "some_device_token",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/devices",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "register_device", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Mock endpoint creation
	suite.mockCreateEndpoint(
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		suite.users[1].OauthUser.Username,
		"some_device_token",
		"new_endpoint_arn",
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Endpoint)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.snsAdapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Endpoint)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created endpoint
	endpoint := new(Endpoint)
	notFound := suite.db.Preload("User").Last(endpoint).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), suite.users[1].ID, uint(endpoint.UserID.Int64))
	assert.Equal(suite.T(), suite.service.cnf.AWS.APNSPlatformApplicationARN, endpoint.ApplicationARN)
	assert.Equal(suite.T(), "new_endpoint_arn", endpoint.ARN)
	assert.Equal(suite.T(), "some_device_token", endpoint.DeviceToken)
	assert.Equal(suite.T(), suite.users[1].OauthUser.Username, endpoint.CustomUserData)
	assert.True(suite.T(), endpoint.Enabled)

	// Check the response body
	assert.Equal(suite.T(), "", strings.TrimRight(w.Body.String(), "\n"))
}

func (suite *NotificationsTestSuite) TestRegisterIOSDeviceWhenAlreadyRegisteredAndNoChanges() {
	// Prepare a request
	payload, err := json.Marshal(&DeviceRequest{
		Platform: PlatformIOS,
		Token:    "the_device_token_1",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/devices",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "register_device", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[0])

	// Mock getting endpoint attributes
	suite.mockGetAttributes(
		suite.endpoints[0].ARN,
		&EndpointAttributes{
			CustomUserData: "the_custom_user_data_1",
			Enabled:        true,
			Token:          "the_device_token_1",
		},
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Endpoint)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.snsAdapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Endpoint)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	assert.Equal(suite.T(), "", strings.TrimRight(w.Body.String(), "\n"))
}

func (suite *NotificationsTestSuite) TestRegisterIOSDeviceWhenAlreadyRegisteredAndEndpointDisabled() {
	// Disable the endpoint
	err := suite.db.Model(suite.endpoints[0]).UpdateColumn("enabled", false).Error
	assert.NoError(suite.T(), err, "Failed to disable the endpoint")

	// Prepare a request
	payload, err := json.Marshal(&DeviceRequest{
		Platform: PlatformIOS,
		Token:    suite.endpoints[0].DeviceToken,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/devices",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "register_device", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[0])

	// Mock getting endpoint attributes
	suite.mockGetAttributes(
		suite.endpoints[0].ARN,
		&EndpointAttributes{
			CustomUserData: "the_custom_user_data_1",
			Enabled:        false,
			Token:          "the_device_token_1",
		},
		nil,
	)

	// Mock setting endpoint attributes
	suite.mockSetAttributes(
		suite.endpoints[0].ARN,
		&EndpointAttributes{
			CustomUserData: "the_custom_user_data_1",
			Enabled:        true,
			Token:          "the_device_token_1",
		},
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Endpoint)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.snsAdapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Endpoint)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated endpoint
	endpoint := new(Endpoint)
	notFound := suite.db.Preload("User").First(endpoint, suite.endpoints[0].ID).
		RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), "the_device_token_1", endpoint.DeviceToken)
	assert.True(suite.T(), endpoint.Enabled)

	// Check the response body
	assert.Equal(suite.T(), "", strings.TrimRight(w.Body.String(), "\n"))
}

func (suite *NotificationsTestSuite) TestRegisterIOSDeviceWhenAlreadyRegisteredAndTokenChanged() {
	// Disable the endpoint
	err := suite.db.Model(suite.endpoints[0]).UpdateColumn("enabled", false).Error
	assert.NoError(suite.T(), err, "Failed to disable the endpoint")

	// Prepare a request
	payload, err := json.Marshal(&DeviceRequest{
		Platform: PlatformIOS,
		Token:    "changed_device_token",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/devices",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "register_device", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[0])

	// Mock getting endpoint attributes
	suite.mockGetAttributes(
		suite.endpoints[0].ARN,
		&EndpointAttributes{
			CustomUserData: "the_custom_user_data_1",
			Enabled:        true,
			Token:          "the_device_token_1",
		},
		nil,
	)

	// Mock setting endpoint attributes
	suite.mockSetAttributes(
		suite.endpoints[0].ARN,
		&EndpointAttributes{
			CustomUserData: "the_custom_user_data_1",
			Enabled:        true,
			Token:          "changed_device_token",
		},
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Endpoint)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.snsAdapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Endpoint)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated endpoint
	endpoint := new(Endpoint)
	notFound := suite.db.Preload("User").First(endpoint, suite.endpoints[0].ID).
		RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), "changed_device_token", endpoint.DeviceToken)
	assert.True(suite.T(), endpoint.Enabled)

	// Check the response body
	assert.Equal(suite.T(), "", strings.TrimRight(w.Body.String(), "\n"))
}
