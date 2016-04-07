package accounts

import (
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrSuperuserOnlyManually ...
	ErrSuperuserOnlyManually = errors.New("Superusers can only be created manually")
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("User not found")
	// ErrEmailTaken ...
	ErrEmailTaken = errors.New("Email already taken")
	// ErrEmailCannotBeChanged ...
	ErrEmailCannotBeChanged = errors.New("Email cannot be changed")
)

// GetName returns user's full name
func (u *User) GetName() string {
	if u.FirstName.Valid && u.LastName.Valid {
		return fmt.Sprintf("%s %s", u.FirstName.String, u.LastName.String)
	}
	return ""
}

// FindUserByOauthUserID looks up a user by oauth user ID and returns it
func (s *Service) FindUserByOauthUserID(oauthUserID uint) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Where(User{
		OauthUserID: util.PositiveIntOrNull(int64(oauthUserID)),
	}).Preload("Account.OauthClient").Preload("OauthUser").Preload("Role").
		First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByEmail looks up a user by email and returns it
func (s *Service) FindUserByEmail(email string) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Joins("inner join oauth_users on oauth_users.id = account_users.oauth_user_id").
		Where("oauth_users.username = ?", email).Preload("Account.OauthClient").
		Preload("OauthUser").Preload("Role").First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByID looks up a user by ID and returns it
func (s *Service) FindUserByID(userID uint) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Preload("Account.OauthClient").Preload("OauthUser").
		Preload("Role").First(user, userID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByFacebookID looks up a user by a Facebook ID and returns it
func (s *Service) FindUserByFacebookID(facebookID string) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Where("facebook_id = ?", facebookID).
		Preload("Account.OauthClient").Preload("OauthUser").Preload("Role").
		First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CreateUser creates a new oauth user and a new account user
func (s *Service) CreateUser(account *Account, userRequest *UserRequest) (*User, error) {
	// Superusers can only be created manually
	if userRequest.Role == roles.Superuser {
		return nil, ErrSuperuserOnlyManually
	}

	// Begin a transaction
	tx := s.db.Begin()

	user, err := s.createUserCommon(
		tx,
		account,
		userRequest,
		"",    // facebook ID
		false, // confirmed
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return user, nil
}

// CreateUserTx creates a new oauth user and a new account user in a transaction
func (s *Service) CreateUserTx(tx *gorm.DB, account *Account, userRequest *UserRequest) (*User, error) {
	// Superusers can only be created manually
	if userRequest.Role == roles.Superuser {
		return nil, ErrSuperuserOnlyManually
	}

	return s.createUserCommon(tx, account, userRequest, "", false)
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(user *User, userRequest *UserRequest) error {
	// Check if email is already taken if
	if user.OauthUser.Username != userRequest.Email {
		return ErrEmailCannotBeChanged
	}

	// Update basic metadata
	if err := s.db.Model(user).UpdateColumns(User{
		FirstName: util.StringOrNull(userRequest.FirstName),
		LastName:  util.StringOrNull(userRequest.LastName),
		Model:     gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		return err
	}

	return nil
}

// GetOrCreateFacebookUser either returns an existing user
// or updates an existing email user with facebook ID or creates a new user
func (s *Service) GetOrCreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error) {
	// Does a user with this facebook ID already exist?
	user, err := s.FindUserByFacebookID(facebookID)

	// User with this facebook ID alraedy exists, return
	if err == nil {
		return user, nil
	}

	// Does a user with this email already exist?
	user, err = s.FindUserByEmail(userRequest.Email)

	// User with this email already exists, update the record and return
	if err == nil {
		// Set the facebook ID and first / last name
		err = s.db.Model(user).UpdateColumns(User{
			FacebookID: util.StringOrNull(facebookID),
			FirstName:  util.StringOrNull(userRequest.FirstName),
			LastName:   util.StringOrNull(userRequest.LastName),
			Confirmed:  true,
			Model:      gorm.Model{UpdatedAt: time.Now()},
		}).Error
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	// Fetch the user role from the database
	role, err := s.findRoleByID(roles.User)
	if err != nil {
		return nil, err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Create a new oauth user
	oauthUser, err := s.GetOauthService().CreateUserTx(
		tx,
		userRequest.Email,
		"", // no password
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new user
	user = NewUser(
		account,
		oauthUser,
		role,
		facebookID,
		userRequest.FirstName,
		userRequest.LastName,
		true, // confirmed
	)

	// Save the user to the database
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Update the meta user ID field
	err = tx.Model(oauthUser).UpdateColumn(oauth.User{MetaUserID: user.ID}).Error
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return user, nil
}

// CreateSuperuser creates a new superuser account
func (s *Service) CreateSuperuser(account *Account, email, password string) (*User, error) {
	// Begin a transaction
	tx := s.db.Begin()

	userRequest := &UserRequest{
		Email:    email,
		Password: password,
		Role:     roles.Superuser,
	}
	user, err := s.createUserCommon(
		tx,
		account,
		userRequest,
		"",   // facebook ID
		true, // confirmed
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return user, nil
}

func (s *Service) createUserCommon(db *gorm.DB, account *Account, userRequest *UserRequest, facebookID string, confirmed bool) (*User, error) {
	// Check if email is already taken
	if s.GetOauthService().UserExists(userRequest.Email) {
		return nil, ErrEmailTaken
	}

	// If a role is not defined in the user request,
	// fall back to the user role
	if userRequest.Role == "" {
		userRequest.Role = roles.User
	}

	// Fetch the role object
	role, err := s.findRoleByID(userRequest.Role)
	if err != nil {
		return nil, err
	}

	// Create a new oauth user
	oauthUser, err := s.GetOauthService().CreateUserTx(
		db,
		userRequest.Email,
		userRequest.Password,
	)
	if err != nil {
		return nil, err
	}

	// Create a new user
	user := NewUser(
		account,
		oauthUser,
		role,
		facebookID,
		userRequest.FirstName,
		userRequest.LastName,
		confirmed,
	)

	// Save the user to the database
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	// Update the meta user ID field
	err = db.Model(oauthUser).UpdateColumn(oauth.User{MetaUserID: user.ID}).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
