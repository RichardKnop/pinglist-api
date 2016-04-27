package facebook

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/response"
)

var (
	// ErrAccountMismatch ...
	ErrAccountMismatch = errors.New("Account mismatch")
)

// Handles requests to login with Facebook access token (POST /v1/facebook/login)
func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated account from the request context
	authenticatedAccount, err := accounts.GetAuthenticatedAccount(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the scope string
	scope, err := s.GetAccountsService().GetOauthService().GetScope(r.Form.Get("scope"))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user data from facebook
	resp, err := s.adapter.GetMe(r.Form.Get("access_token"))
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	logger.Infof("%v", resp)

	// Initialise variables with from facebook
	var (
		facebookID = fmt.Sprintf("%s", resp["id"])
		email      = fmt.Sprintf("%s", resp["email"])
		firstName  = fmt.Sprintf("%s", resp["first_name"])
		lastName   = fmt.Sprintf("%s", resp["last_name"])
		user       *accounts.User
	)

	// Get or create a new user based on facebook ID and other details
	user, err = s.GetAccountsService().GetOrCreateFacebookUser(
		authenticatedAccount,
		facebookID,
		&accounts.UserRequest{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check that the same account is being used
	if authenticatedAccount.ID != user.Account.ID {
		response.UnauthorizedError(w, ErrAccountMismatch.Error())
		return
	}

	// Log in the user
	accessToken, refreshToken, err := s.GetAccountsService().GetOauthService().Login(
		user.Account.OauthClient,
		user.OauthUser,
		scope,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON access token to the response
	accessTokenRespone := &oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    s.cnf.Oauth.AccessTokenLifetime,
		TokenType:    oauth.TokenType,
		Scope:        accessToken.Scope,
		RefreshToken: refreshToken.Token,
	}
	response.WriteJSON(w, accessTokenRespone, 200)
}
