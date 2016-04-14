package notifications

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/stretchr/testify/mock"
)

// ServiceMock is a mocked object implementing ServiceInterface
type ServiceMock struct {
	mock.Mock
}

// GetAccountsService ...
func (_m *ServiceMock) GetAccountsService() accounts.ServiceInterface {
	ret := _m.Called()

	var r0 accounts.ServiceInterface
	if rf, ok := ret.Get(0).(func() accounts.ServiceInterface); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(accounts.ServiceInterface)
	}

	return r0
}

// FindEndpointByUserIDAndApplicationARN ...
func (_m *ServiceMock) FindEndpointByUserIDAndApplicationARN(userID uint, applicationARN string) (*Endpoint, error) {
	ret := _m.Called(userID, applicationARN)

	var r0 *Endpoint
	if rf, ok := ret.Get(0).(func(uint, string) *Endpoint); ok {
		r0 = rf(userID, applicationARN)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Endpoint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, string) error); ok {
		r1 = rf(userID, applicationARN)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishMessage ...
func (_m *ServiceMock) PublishMessage(endpointARN string, msg string, opt map[string]interface{}) (string, error) {
	ret := _m.Called(endpointARN, msg, opt)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, map[string]interface{}) string); ok {
		r0 = rf(endpointARN, msg, opt)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, map[string]interface{}) error); ok {
		r1 = rf(endpointARN, msg, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ServiceMock) registerDeviceHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
