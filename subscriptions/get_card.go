package subscriptions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/gorilla/mux"
)

var (
	// ErrGetCardPermission ...
	ErrGetCardPermission = errors.New("Need permission to get card")
)

// Handles calls to get a card (GET /v1/cards/{id:[0-9]+})
func (s *Service) getCardHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the id from request URI and type assert it
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the card we want to delete
	card, err := s.FindCardByID(uint(cardID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkGetCardPermissions(authenticatedUser, card); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Create response
	cardResponse, err := NewCardResponse(card)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, cardResponse, http.StatusOK)
}

func checkGetCardPermissions(authenticatedUser *accounts.User, card *Card) error {
	// Superusers can get any cards
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can get their own cards
	if authenticatedUser.ID == card.Customer.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrGetCardPermission
}
