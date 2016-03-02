package accounts

import (
	"net/http"
	"strconv"
)

// GetAccountFromQueryString parses account_id from the query string and
// returns a matching *Account instance or an error (requires user authentication)
func (s *Service) GetAccountFromQueryString(r *http.Request) (*Account, error) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		return nil, err
	}

	// If no account_id query string parameter found, just return
	if r.URL.Query().Get("account_id") == "" {
		return nil, nil
	}

	// If the account_id is present in the query string, try to convert it to int
	accountID, err := strconv.Atoi(r.URL.Query().Get("account_id"))
	if err != nil {
		return nil, err
	}

	// If the account ID matches the authenticated user's account, just return it
	if uint(accountID) == authenticatedUser.Account.ID {
		return authenticatedUser.Account, nil
	}

	// Fetch the account from the database
	account, err := s.FindAccountByID(uint(accountID))
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetUserFromQueryString parses user_id from the query string and
// returns a matching *User instance or an error (requires user authentication)
func (s *Service) GetUserFromQueryString(r *http.Request) (*User, error) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		return nil, err
	}

	// If no user_id query string parameter found, just return
	if r.URL.Query().Get("user_id") == "" {
		return nil, nil
	}

	// If the user_id is present in the query string, try to convert it to int
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return nil, err
	}

	// If the user ID matches the authenticated user, just return it
	if uint(userID) == authenticatedUser.ID {
		return authenticatedUser, nil
	}

	// Fetch the user from the database
	user, err := s.FindUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
