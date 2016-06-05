package alarms

import (
	"github.com/stretchr/testify/mock"
)

// SlackFactoryMock is a mocked object implementing SlackFactoryInterface
type SlackFactoryMock struct {
	mock.Mock
}

// NewIncidentMessage ...
func (_m *SlackFactoryMock) NewIncidentMessage(incident *Incident) string {
	ret := _m.Called(incident)

	var r0 string
	if rf, ok := ret.Get(0).(func(*Incident) string); ok {
		r0 = rf(incident)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewIncidentsResolvedMessage ...
func (_m *SlackFactoryMock) NewIncidentsResolvedMessage(alarm *Alarm) string {
	ret := _m.Called(alarm)

	var r0 string
	if rf, ok := ret.Get(0).(func(*Alarm) string); ok {
		r0 = rf(alarm)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
