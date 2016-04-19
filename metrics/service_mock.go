package metrics

import (
	"time"

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

// PartitionResponseTime ...
func (_m *ServiceMock) PartitionResponseTime(parentTableName string, now time.Time) error {
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

// LogResponseTime ...
func (_m *ServiceMock) LogResponseTime(timestamp time.Time, referenceID uint, value int64) error {
	ret := _m.Called(timestamp, referenceID, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time, uint, int64) error); ok {
		r0 = rf(timestamp, referenceID, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PaginatedResponseTimesCount ...
func (_m *ServiceMock) PaginatedResponseTimesCount(referenceID int, dateTrunc string, from *time.Time, to *time.Time) (int, error) {
	ret := _m.Called(referenceID, dateTrunc, from, to)

	var r0 int
	if rf, ok := ret.Get(0).(func(int, string, *time.Time, *time.Time) int); ok {
		r0 = rf(referenceID, dateTrunc, from, to)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string, *time.Time, *time.Time) error); ok {
		r1 = rf(referenceID, dateTrunc, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindPaginatedResponseTimes ...
func (_m *ServiceMock) FindPaginatedResponseTimes(offset int, limit int, orderBy string, referenceID int, dateTrunc string, from *time.Time, to *time.Time) ([]*ResponseTime, error) {
	ret := _m.Called(offset, limit, orderBy, referenceID, dateTrunc, from, to)

	var r0 []*ResponseTime
	if rf, ok := ret.Get(0).(func(int, int, string, int, string, *time.Time, *time.Time) []*ResponseTime); ok {
		r0 = rf(offset, limit, orderBy, referenceID, dateTrunc, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ResponseTime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, string, int, string, *time.Time, *time.Time) error); ok {
		r1 = rf(offset, limit, orderBy, referenceID, dateTrunc, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
