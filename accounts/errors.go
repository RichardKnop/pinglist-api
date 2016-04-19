package accounts

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/oauth"
)

var (
	errStatusCodeMap = map[error]int{
		ErrSuperuserOnlyManually: http.StatusBadRequest,
		ErrRoleNotFound:          http.StatusBadRequest,
		oauth.ErrUsernameTaken:   http.StatusBadRequest,
	}
)
