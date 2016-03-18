package accounts

import (
	"github.com/RichardKnop/pinglist-api/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewConfirmationEmail(confirmation *Confirmation) *email.Email
	NewPasswordResetEmail(passwordReset *PasswordReset) *email.Email
}
