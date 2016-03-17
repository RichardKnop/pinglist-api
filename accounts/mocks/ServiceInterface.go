package mocks

import "github.com/RichardKnop/pinglist-api/accounts"
import "github.com/stretchr/testify/mock"

import "net/http"
import "github.com/RichardKnop/pinglist-api/config"
import "github.com/RichardKnop/pinglist-api/oauth"
import "github.com/jinzhu/gorm"

type ServiceInterface struct {
	mock.Mock
}

func (_m *ServiceInterface) GetConfig() *config.Config {
	ret := _m.Called()

	var r0 *config.Config
	if rf, ok := ret.Get(0).(func() *config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	return r0
}
func (_m *ServiceInterface) GetOauthService() oauth.ServiceInterface {
	ret := _m.Called()

	var r0 oauth.ServiceInterface
	if rf, ok := ret.Get(0).(func() oauth.ServiceInterface); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(oauth.ServiceInterface)
	}

	return r0
}
func (_m *ServiceInterface) FindAccountByOauthClientID(oauthClientID uint) (*accounts.Account, error) {
	ret := _m.Called(oauthClientID)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(uint) *accounts.Account); ok {
		r0 = rf(oauthClientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(oauthClientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindAccountByID(accountID uint) (*accounts.Account, error) {
	ret := _m.Called(accountID)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(uint) *accounts.Account); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateAccount(name string, description string, key string, secret string, redirectURI string) (*accounts.Account, error) {
	ret := _m.Called(name, description, key, secret, redirectURI)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) *accounts.Account); ok {
		r0 = rf(name, description, key, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string) error); ok {
		r1 = rf(name, description, key, secret, redirectURI)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByOauthUserID(oauthUserID uint) (*accounts.User, error) {
	ret := _m.Called(oauthUserID)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(uint) *accounts.User); ok {
		r0 = rf(oauthUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(oauthUserID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByEmail(email string) (*accounts.User, error) {
	ret := _m.Called(email)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(string) *accounts.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByID(userID uint) (*accounts.User, error) {
	ret := _m.Called(userID)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(uint) *accounts.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
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
func (_m *ServiceInterface) FindUserByFacebookID(facebookID string) (*accounts.User, error) {
	ret := _m.Called(facebookID)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(string) *accounts.User); ok {
		r0 = rf(facebookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(facebookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateUser(account *accounts.Account, userRequest *accounts.UserRequest) (*accounts.User, error) {
	ret := _m.Called(account, userRequest)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*accounts.Account, *accounts.UserRequest) *accounts.User); ok {
		r0 = rf(account, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Account, *accounts.UserRequest) error); ok {
		r1 = rf(account, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateUserTx(tx *gorm.DB, account *accounts.Account, userRequest *accounts.UserRequest) (*accounts.User, error) {
	ret := _m.Called(tx, account, userRequest)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, *accounts.Account, *accounts.UserRequest) *accounts.User); ok {
		r0 = rf(tx, account, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *accounts.Account, *accounts.UserRequest) error); ok {
		r1 = rf(tx, account, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) UpdateUser(user *accounts.User, userRequest *accounts.UserRequest) error {
	ret := _m.Called(user, userRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accounts.User, *accounts.UserRequest) error); ok {
		r0 = rf(user, userRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) FindConfirmationByReference(reference string) (*accounts.Confirmation, error) {
	ret := _m.Called(reference)

	var r0 *accounts.Confirmation
	if rf, ok := ret.Get(0).(func(string) *accounts.Confirmation); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Confirmation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(reference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) ConfirmUser(user *accounts.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accounts.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) CreateFacebookUser(account *accounts.Account, facebookID string, userRequest *accounts.UserRequest) (*accounts.User, error) {
	ret := _m.Called(account, facebookID, userRequest)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*accounts.Account, string, *accounts.UserRequest) *accounts.User); ok {
		r0 = rf(account, facebookID, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Account, string, *accounts.UserRequest) error); ok {
		r1 = rf(account, facebookID, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateSuperuser(account *accounts.Account, email string, password string) (*accounts.User, error) {
	ret := _m.Called(account, email, password)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*accounts.Account, string, string) *accounts.User); ok {
		r0 = rf(account, email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Account, string, string) error); ok {
		r1 = rf(account, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetAccountFromQueryString(r *http.Request) (*accounts.Account, error) {
	ret := _m.Called(r)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(*http.Request) *accounts.Account); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetUserFromQueryString(r *http.Request) (*accounts.User, error) {
	ret := _m.Called(r)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*http.Request) *accounts.User); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) createUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) getMyUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
