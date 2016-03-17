package alarms

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
)

// Handles calls to list alarm regions (GET /v1/alarms/regions)
func (s *Service) listRegionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	_, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Fetch the regions
	var regions []*Region
	if err := s.db.Order("id").Find(&regions).Error; err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	listRegionsResponse, err := NewListRegionsResponse(regions)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listRegionsResponse, http.StatusOK)
}
