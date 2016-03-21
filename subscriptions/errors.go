package subscriptions

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrUserCanOnlyHaveOneActiveSubscription: http.StatusBadRequest,
		ErrPlanNotFound:                         http.StatusBadRequest,
	}
)
