package accounts

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrPasswordResetNotFound ...
	ErrPasswordResetNotFound = errors.New("Password reset not found")
)

// FindPasswordResetByReference looks up a password reset by a reference
func (s *Service) FindPasswordResetByReference(reference string) (*PasswordReset, error) {
	// Fetch the password reset from the database
	passwordReset := new(PasswordReset)
	validFor := time.Duration(s.cnf.Pinglist.PasswordResetLifetime) * time.Second
	notFound := s.db.Where(
		"reference = ? AND created_at > ?",
		reference,
		time.Now().Add(-validFor),
	).Preload("User.OauthUser").First(passwordReset).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPasswordResetNotFound
	}

	return passwordReset, nil
}

// ResetPassword sets a new password and deletes the password reset record
func (s *Service) ResetPassword(passwordReset *PasswordReset, password string) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Set the new password
	err := s.oauthService.SetPasswordTx(tx, passwordReset.User.OauthUser, password)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Soft delete the password reset
	if err := tx.Delete(passwordReset).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

// findUserPasswordReset returns the first password reset for a user
func (s *Service) findUserPasswordReset(user *User) (*PasswordReset, error) {
	// Fetch the password reset from the database
	passwordReset := new(PasswordReset)
	notFound := s.db.Where(PasswordReset{
		UserID: util.PositiveIntOrNull(int64(user.ID)),
	}).Preload("User.OauthUser").First(passwordReset).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPasswordResetNotFound
	}

	return passwordReset, nil
}

func (s *Service) createPasswordReset(user *User) (*PasswordReset, error) {
	// Begin a transaction
	tx := s.db.Begin()

	// Soft delete old password resets
	if err := tx.Where("user_id = ?", user.ID).Delete(new(PasswordReset)).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new password reset
	passwordReset := NewPasswordReset(user)

	// Save the password reset to the database
	if err := tx.Create(passwordReset).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Assign related object
	passwordReset.User = user

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Send password reset email
	go func() {
		passwordResetEmail := s.emailFactory.NewPasswordResetEmail(passwordReset)

		// Try to send the password reset email
		if err := s.emailService.Send(passwordResetEmail); err != nil {
			logger.ERROR.Printf("Send email error: %s", err)
			return
		}

		// If the email was sent successfully, update the email_sent flag
		now := time.Now()
		s.db.Model(passwordReset).UpdateColumns(PasswordReset{
			EmailSent:   true,
			EmailSentAt: util.TimeOrNull(&now),
			Model:       gorm.Model{UpdatedAt: now},
		})
	}()

	return passwordReset, nil
}
