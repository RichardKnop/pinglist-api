package alarms

import (
	"github.com/stretchr/testify/mock"
)

// SlackAdapterMock is a mocked object implementing SlackAdapterInterface
type SlackAdapterMock struct {
	mock.Mock
}

// SendMessage ...
func (_m *SlackAdapterMock) SendMessage(channel string, username string, text string, emoji string) error {
	ret := _m.Called(channel, username, text, emoji)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(channel, username, text, emoji)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
