package subscriptions

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
)

// Handles calls to list cards plans (GET /v1/plans)
func (s *Service) listPlansHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated account from the request context
	_, err := accounts.GetAuthenticatedAccount(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Fetch all the plans
	var plans []*Plan
	if err := s.db.Order("id").Find(&plans).Error; err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get page and limit
	page, limit, err := pagination.GetPageLimit(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Count total number of results
	count := len(plans)

	// Get pagination links
	first, last, previous, next, err := pagination.GetPaginationLinks(
		r.URL,
		count,
		page,
		limit,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create response
	self := util.GetCurrentURL(r)
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
