package subscriptions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/gorilla/mux"
)

var (
	// ErrDeleteCardPermission ...
	ErrDeleteCardPermission = errors.New("Need permission to delete card")
)

// Handles calls to delete a card (DELETE /v1/cards/{id:[0-9]+})
func (s *Service) deleteCardHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := checkDeleteCardPermissions(authenticatedUser, card); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Delete the card
	if err := s.deleteCard(card); err != nil {
		logger.ERROR.Printf("Delete card error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// 204 no content response
	response.NoContent(w)
}

func checkDeleteCardPermissions(authenticatedUser *accounts.User, card *Card) error {
	// Superusers can delete any cards
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can delete their own cards
	if authenticatedUser.ID == card.Customer.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrDeleteCardPermission
}
