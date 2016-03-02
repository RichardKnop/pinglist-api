package subscriptions

import (
	"io"
	"net/url"

	"github.com/stretchr/testify/mock"
	stripe "github.com/stripe/stripe-go"
)

// StripeBackendMock is a mocked object implementing stripe.Backend
type StripeBackendMock struct {
	mock.Mock
}

// Call ...
func (_m *StripeBackendMock) Call(method string, path string, key string, body *url.Values, params *stripe.Params, v interface{}) error {
	ret := _m.Called(method, path, key, body, params, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, *url.Values, *stripe.Params, interface{}) error); ok {
		r0 = rf(method, path, key, body, params, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CallMultipart ...
func (_m *StripeBackendMock) CallMultipart(method string, path string, key string, boundary string, body io.Reader, params *stripe.Params, v interface{}) error {
	ret := _m.Called(method, path, key, boundary, body, params, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, io.Reader, *stripe.Params, interface{}) error); ok {
		r0 = rf(method, path, key, boundary, body, params, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
