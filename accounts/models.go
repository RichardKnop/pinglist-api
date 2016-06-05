package accounts

import (
	"database/sql"

	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pborman/uuid"
)

// Account ...
type Account struct {
	gorm.Model
	OauthClientID sql.NullInt64 `sql:"index;not null"`
	OauthClient   *oauth.Client
	Name          string         `sql:"type:varchar(100);unique;not null"`
	Description   sql.NullString `sql:"type:varchar(200)"`
}

// TableName specifies table name
func (p *Account) TableName() string {
	return "account_accounts"
}

// Role is a one of roles user can have
type Role struct {
	database.TimestampModel
	ID   string `gorm:"primary_key" sql:"type:varchar(20)"`
	Name string `sql:"type:varchar(50);unique;not null"`
}

// TableName specifies table name
func (r *Role) TableName() string {
	return "account_roles"
}

// User ...
type User struct {
	gorm.Model
	AccountID    sql.NullInt64  `sql:"index;not null"`
	OauthUserID  sql.NullInt64  `sql:"index;not null"`
	RoleID       sql.NullString `sql:"type:varchar(20);index;not null"`
	Account      *Account
	OauthUser    *oauth.User
	Role         *Role
	FacebookID   sql.NullString `sql:"type:varchar(60);unique"`
	FirstName    sql.NullString `sql:"type:varchar(100)"`
	LastName     sql.NullString `sql:"type:varchar(100)"`
	Confirmed    bool           `sql:"index;not null"`
	SlackAPIKey  sql.NullString `sql:"type:varchar(100)"`
	SlackChannel sql.NullString `sql:"type:varchar(100)"`
}

// TableName specifies table name
func (u *User) TableName() string {
	return "account_users"
}

// Confirmation ...
type Confirmation struct {
	gorm.Model
	UserID      sql.NullInt64 `sql:"index;not null"`
	User        *User
	Reference   string `sql:"type:varchar(40);unique;not null"`
	EmailSent   bool   `sql:"index;not null"`
	EmailSentAt pq.NullTime
}

// TableName specifies table name
func (c *Confirmation) TableName() string {
	return "account_confirmations"
}

// Invitation ...
type Invitation struct {
	gorm.Model
	InvitedUserID   sql.NullInt64 `sql:"index;not null"`
	InvitedByUserID sql.NullInt64 `sql:"index;not null"`
	InvitedByUser   *User
	InvitedUser     *User
	Reference       string `sql:"type:varchar(40);unique;not null"`
	EmailSent       bool   `sql:"index;not null"`
	EmailSentAt     pq.NullTime
}

// TableName specifies table name
func (i *Invitation) TableName() string {
	return "account_invitations"
}

// PasswordReset ...
type PasswordReset struct {
	gorm.Model
	UserID      sql.NullInt64 `sql:"index;not null"`
	User        *User
	Reference   string `sql:"type:varchar(40);unique;not null"`
	EmailSent   bool   `sql:"index;not null"`
	EmailSentAt pq.NullTime
}

// TableName specifies table name
func (p *PasswordReset) TableName() string {
	return "account_password_resets"
}

// NewAccount creates new Account instance
func NewAccount(oauthClient *oauth.Client, name, description string) *Account {
	oauthClientID := util.PositiveIntOrNull(int64(oauthClient.ID))
	account := &Account{
		OauthClientID: oauthClientID,
		Name:          name,
		Description:   util.StringOrNull(description),
	}
	return account
}

// NewUser creates new User instance
func NewUser(account *Account, oauthUser *oauth.User, role *Role, facebookID, firstName, lastName string, confirmed bool, slackAPIKey, slackChannel string) *User {
	accountID := util.PositiveIntOrNull(int64(account.ID))
	oauthUserID := util.PositiveIntOrNull(int64(oauthUser.ID))
	roleID := util.StringOrNull(role.ID)
	user := &User{
		AccountID:    accountID,
		OauthUserID:  oauthUserID,
		RoleID:       roleID,
		FacebookID:   util.StringOrNull(facebookID),
		FirstName:    util.StringOrNull(firstName),
		LastName:     util.StringOrNull(lastName),
		Confirmed:    confirmed,
		SlackAPIKey:  util.StringOrNull(slackAPIKey),
		SlackChannel: util.StringOrNull(slackChannel),
	}
	return user
}

// NewConfirmation creates new Confirmation instance
func NewConfirmation(user *User) *Confirmation {
	userID := util.PositiveIntOrNull(int64(user.ID))
	confirmation := &Confirmation{
		UserID:      userID,
		Reference:   uuid.New(),
		EmailSentAt: util.TimeOrNull(nil),
	}
	return confirmation
}

// NewInvitation creates new Invitation instance
func NewInvitation(invitedUser, invitedByUser *User) *Invitation {
	invitedUserID := util.PositiveIntOrNull(int64(invitedUser.ID))
	invitedByUserID := util.PositiveIntOrNull(int64(invitedByUser.ID))
	invitation := &Invitation{
		InvitedUserID:   invitedUserID,
		InvitedByUserID: invitedByUserID,
		Reference:       uuid.New(),
		EmailSentAt:     util.TimeOrNull(nil),
	}
	return invitation
}

// NewPasswordReset creates new PasswordReset instance
func NewPasswordReset(user *User) *PasswordReset {
	userID := util.PositiveIntOrNull(int64(user.ID))
	passwordReset := &PasswordReset{
		UserID:      userID,
		Reference:   uuid.New(),
		EmailSentAt: util.TimeOrNull(nil),
	}
	return passwordReset
}
