package alarms

import (
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/stretchr/testify/mock"
)

// EmailFactoryMock is a mocked object implementing EmailFactoryInterface
type EmailFactoryMock struct {
	mock.Mock
}

// NewAlarmDownEmail ...
func (_m *EmailFactoryMock) NewAlarmDownEmail(alarm *Alarm) *email.Email {
	ret := _m.Called(alarm)

	var r0 *email.Email
	if rf, ok := ret.Get(0).(func(*Alarm) *email.Email); ok {
		r0 = rf(alarm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Email)
		}
	}

	return r0
}

// NewAlarmUpEmail ...
func (_m *EmailFactoryMock) NewAlarmUpEmail(alarm *Alarm) *email.Email {
	ret := _m.Called(alarm)

	var r0 *email.Email
	if rf, ok := ret.Get(0).(func(*Alarm) *email.Email); ok {
		r0 = rf(alarm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Email)
		}
	}

	return r0
}
