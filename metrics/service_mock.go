package metrics

import (
  "time"
  "github.com/stretchr/testify/mock"
  "github.com/RichardKnop/pinglist-api/accounts"
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

// PartitionRequestTime ...
func (_m *ServiceMock) PartitionRequestTime(parentTableName string, now time.Time) error {
	ret := _m.Called(parentTableName, now)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Time) error); ok {
		r0 = rf(parentTableName, now)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RotateSubTables ...
func (_m *ServiceMock) RotateSubTables() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LogRequestTime ...
func (_m *ServiceMock) LogRequestTime(timestamp time.Time, referenceID uint, value int64) error {
	ret := _m.Called(timestamp, referenceID, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time, uint, int64) error); ok {
		r0 = rf(timestamp, referenceID, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PaginatedRequestTimesCount ...
func (_m *ServiceMock) PaginatedRequestTimesCount(referenceID uint) (int, error) {
	ret := _m.Called(referenceID)

	var r0 int
	if rf, ok := ret.Get(0).(func(uint) int); ok {
		r0 = rf(referenceID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(referenceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindPaginatedRequestTimes ...
func (_m *ServiceMock) FindPaginatedRequestTimes(offset int, limit int, orderBy string, referenceID uint) ([]*metrics.RequestTime, error) {
	ret := _m.Called(offset, limit, orderBy, referenceID)

	var r0 []*metrics.RequestTime
	if rf, ok := ret.Get(0).(func(int, int, string, uint) []*metrics.RequestTime); ok {
		r0 = rf(offset, limit, orderBy, referenceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*metrics.RequestTime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, string, uint) error); ok {
		r1 = rf(offset, limit, orderBy, referenceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
