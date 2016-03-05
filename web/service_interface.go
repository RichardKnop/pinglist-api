package web

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetAccountsService() accounts.ServiceInterface

	// Needed for the newRoutes to be able to register handlers
	registerForm(w http.ResponseWriter, r *http.Request)
	register(w http.ResponseWriter, r *http.Request)
	confirmEmail(w http.ResponseWriter, r *http.Request)
	authorizeForm(w http.ResponseWriter, r *http.Request)
	authorize(w http.ResponseWriter, r *http.Request)
	loginForm(w http.ResponseWriter, r *http.Request)
	login(w http.ResponseWriter, r *http.Request)
	logout(w http.ResponseWriter, r *http.Request)
}
