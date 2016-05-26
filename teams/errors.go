package teams

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrMaxTeamsLimitReached:          http.StatusBadRequest,
		ErrMaxMembersPerTeamLimitReached: http.StatusBadRequest,
		ErrCannotAddYourself:             http.StatusBadRequest,
	}
)

func getErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	switch err.(type) {
	case ErrUserCanOnlyBeMemberOfOneTeam:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
