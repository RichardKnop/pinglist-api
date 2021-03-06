package accounts

import (
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/stretchr/testify/mock"
)

// EmailFactoryMock is a mocked object implementing EmailFactoryInterface
type EmailFactoryMock struct {
	mock.Mock
}

// NewConfirmationEmail ...
func (_m *EmailFactoryMock) NewConfirmationEmail(confirmation *Confirmation) *email.Email {
	ret := _m.Called(confirmation)

	var r0 *email.Email
	if rf, ok := ret.Get(0).(func(*Confirmation) *email.Email); ok {
		r0 = rf(confirmation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Email)
		}
	}

	return r0
}

// NewInvitationEmail just records the activity, and returns what the Mock object tells it to
func (_m *EmailFactoryMock) NewInvitationEmail(invitation *Invitation) *email.Email {
	ret := _m.Called(invitation)

	var r0 *email.Email
	if rf, ok := ret.Get(0).(func(*Invitation) *email.Email); ok {
		r0 = rf(invitation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Email)
		}
	}

	return r0
}

// NewPasswordResetEmail ...
func (_m *EmailFactoryMock) NewPasswordResetEmail(passwordReset *PasswordReset) *email.Email {
	ret := _m.Called(passwordReset)

	var r0 *email.Email
	if rf, ok := ret.Get(0).(func(*PasswordReset) *email.Email); ok {
		r0 = rf(passwordReset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Email)
		}
	}

	return r0
}
