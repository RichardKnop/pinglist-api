package subscriptions

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
)

// Handles calls to list subscription plans (GET /v1/plans)
func (s *Service) listPlansHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	_, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Fetch the plans
	var plans []*Plan
	if err := s.db.Order("id").Find(&plans).Error; err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	count, page := len(plans), 1
	self, first, last := util.GetCurrentURL(r), util.GetCurrentURL(r), util.GetCurrentURL(r)
	next, previous := "", ""
	listPlansResponse, err := NewListPlansResponse(
		count, page,
		self, first, last, next, previous,
		plans,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listPlansResponse, http.StatusOK)
}
