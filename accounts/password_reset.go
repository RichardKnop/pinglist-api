package accounts

import (
	"errors"
)

var (
	// ErrPasswordResetNotFound ...
	ErrPasswordResetNotFound = errors.New("Password reset not found")
)

// FindPasswordResetByReference looks up a password reset by a reference
func (s *Service) FindPasswordResetByReference(reference string) (*PasswordReset, error) {
	// Fetch the password reset from the database
	passwordReset := new(PasswordReset)
	notFound := s.db.Where("reference = ?", reference).
		Preload("User.OauthUser").First(passwordReset).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPasswordResetNotFound
	}

	return passwordReset, nil
}
