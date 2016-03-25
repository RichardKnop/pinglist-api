package teams

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrMaxTeamsLimitReached:          http.StatusBadRequest,
		ErrMaxMembersPerTeamLimitReached: http.StatusBadRequest,
	}
)
