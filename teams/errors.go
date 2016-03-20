package teams

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrUserCanOnlyCreateOneTeam:   http.StatusBadRequest,
		ErrMaxTeamMembersLimitReached: http.StatusBadRequest,
	}
)
