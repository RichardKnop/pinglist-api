package teams

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface
	FindTeamByID(teamID uint) (*Team, error)
	FindTeamByMemberID(memberID uint) (*Team, error)

	// Needed for the newRoutes to be able to register handlers
	createTeamHandler(w http.ResponseWriter, r *http.Request)
	getTeamHandler(w http.ResponseWriter, r *http.Request)
	listTeamsHandler(w http.ResponseWriter, r *http.Request)
	updateTeamHandler(w http.ResponseWriter, r *http.Request)
	deleteTeamHandler(w http.ResponseWriter, r *http.Request)
}
