package alarms

import (
	"github.com/stretchr/testify/mock"
)

// SlackAdapterMock is a mocked object implementing SlackAdapterInterface
type SlackAdapterMock struct {
	mock.Mock
}

// SendMessage ...
func (_m *SlackAdapterMock) SendMessage(incomingWebhook string, channel string, username string, emoji string, text string) error {
	ret := _m.Called(incomingWebhook, channel, username, emoji, text)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) error); ok {
		r0 = rf(incomingWebhook, channel, username, emoji, text)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
