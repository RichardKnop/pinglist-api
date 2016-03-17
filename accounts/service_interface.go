package accounts

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/jinzhu/gorm"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetOauthService() oauth.ServiceInterface
	FindAccountByOauthClientID(oauthClientID uint) (*Account, error)
	FindAccountByID(accountID uint) (*Account, error)
	CreateAccount(name, description, key, secret, redirectURI string) (*Account, error)
	FindUserByOauthUserID(oauthUserID uint) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByID(userID uint) (*User, error)
	FindUserByFacebookID(facebookID string) (*User, error)
	CreateUser(account *Account, userRequest *UserRequest) (*User, error)
	CreateUserTx(tx *gorm.DB, account *Account, userRequest *UserRequest) (*User, error)
	UpdateUser(user *User, userRequest *UserRequest) error
	FindConfirmationByReference(reference string) (*Confirmation, error)
	ConfirmUser(user *User) error
	CreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error)
	CreateSuperuser(account *Account, email, password string) (*User, error)
	GetAccountFromQueryString(r *http.Request) (*Account, error)
	GetUserFromQueryString(r *http.Request) (*User, error)

	// Needed for the newRoutes to be able to register handlers
	createUserHandler(w http.ResponseWriter, r *http.Request)
	getMyUserHandler(w http.ResponseWriter, r *http.Request)
	updateUserHandler(w http.ResponseWriter, r *http.Request)
}
