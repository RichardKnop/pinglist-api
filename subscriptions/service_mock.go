package subscriptions

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

// FindPlanByID ...
func (_m *ServiceMock) FindPlanByID(planID uint) (*Plan, error) {
	ret := _m.Called(planID)

	var r0 *Plan
	if rf, ok := ret.Get(0).(func(uint) *Plan); ok {
		r0 = rf(planID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Plan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(planID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindPlanByPlanID ...
func (_m *ServiceMock) FindPlanByPlanID(planID string) (*Plan, error) {
	ret := _m.Called(planID)

	var r0 *Plan
	if rf, ok := ret.Get(0).(func(string) *Plan); ok {
		r0 = rf(planID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Plan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(planID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindCustomerByID ...
func (_m *ServiceMock) FindCustomerByID(customerID uint) (*Customer, error) {
	ret := _m.Called(customerID)

	var r0 *Customer
	if rf, ok := ret.Get(0).(func(uint) *Customer); ok {
		r0 = rf(customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Customer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindCustomerByCustomerID ...
func (_m *ServiceMock) FindCustomerByCustomerID(customerID string) (*Customer, error) {
	ret := _m.Called(customerID)

	var r0 *Customer
	if rf, ok := ret.Get(0).(func(string) *Customer); ok {
		r0 = rf(customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Customer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSubscriptionByID ...
func (_m *ServiceMock) FindSubscriptionByID(subscriptionID uint) (*Subscription, error) {
	ret := _m.Called(subscriptionID)

	var r0 *Subscription
	if rf, ok := ret.Get(0).(func(uint) *Subscription); ok {
		r0 = rf(subscriptionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSubscriptionBySubscriptionID ...
func (_m *ServiceMock) FindSubscriptionBySubscriptionID(subscriptionID string) (*Subscription, error) {
	ret := _m.Called(subscriptionID)

	var r0 *Subscription
	if rf, ok := ret.Get(0).(func(string) *Subscription); ok {
		r0 = rf(subscriptionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindActiveUserSubscription ...
func (_m *ServiceMock) FindActiveUserSubscription(userID uint) (*Subscription, error) {
	ret := _m.Called(userID)

	var r0 *Subscription
	if rf, ok := ret.Get(0).(func(uint) *Subscription); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ServiceMock) listPlansHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) subscribeUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) listSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) cancelSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) stripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
