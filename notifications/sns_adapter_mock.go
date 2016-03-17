package notifications

import "github.com/stretchr/testify/mock"

// SNSAdapterMock is a mocked object implementing SNSAdapterInterface
type SNSAdapterMock struct {
	mock.Mock
}

// CreateEndpoint ...
func (_m *SNSAdapterMock) CreateEndpoint(applicationARN string, customUserData string, deviceToken string) (string, error) {
	ret := _m.Called(applicationARN, customUserData, deviceToken)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(applicationARN, customUserData, deviceToken)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(applicationARN, customUserData, deviceToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEndpointAttributes ...
func (_m *SNSAdapterMock) GetEndpointAttributes(endpointARN string) (*EndpointAttributes, error) {
	ret := _m.Called(endpointARN)

	var r0 *EndpointAttributes
	if rf, ok := ret.Get(0).(func(string) *EndpointAttributes); ok {
		r0 = rf(endpointARN)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*EndpointAttributes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(endpointARN)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEndpointAttributes ...
func (_m *SNSAdapterMock) SetEndpointAttributes(endpointARN string, endpointAttributes *EndpointAttributes) error {
	ret := _m.Called(endpointARN, endpointAttributes)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *EndpointAttributes) error); ok {
		r0 = rf(endpointARN, endpointAttributes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishMessage ...
func (_m *SNSAdapterMock) PublishMessage(endpointARN string, message string, opt map[string]interface{}) (string, error) {
	ret := _m.Called(endpointARN, message, opt)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, map[string]interface{}) string); ok {
		r0 = rf(endpointARN, message, opt)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, map[string]interface{}) error); ok {
		r1 = rf(endpointARN, message, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
