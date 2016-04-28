package alarms

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
)

// Handles calls to list alarm regions (GET /v1/regions)
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
	count, page := len(regions), 1
	self, first, last := util.GetCurrentURL(r), util.GetCurrentURL(r), util.GetCurrentURL(r)
	next, previous := "", ""
	listRegionsResponse, err := NewListRegionsResponse(
		count, page,
		self, first, last, previous, next,
		regions,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listRegionsResponse, http.StatusOK)
}
