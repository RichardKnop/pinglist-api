package teams

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

// FindTeamByID ...
func (_m *ServiceMock) FindTeamByID(teamID uint) (*Team, error) {
	ret := _m.Called(teamID)

	var r0 *Team
	if rf, ok := ret.Get(0).(func(uint) *Team); ok {
		r0 = rf(teamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTeamByMemberID ...
func (_m *ServiceMock) FindTeamByMemberID(memberID uint) (*Team, error) {
	ret := _m.Called(memberID)

	var r0 *Team
	if rf, ok := ret.Get(0).(func(uint) *Team); ok {
		r0 = rf(memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(memberID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ServiceMock) createTeamHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) listTeamsHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) updateTeamHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
