package subscriptions

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrUserCanOnlyHaveOneActiveSubscription: http.StatusBadRequest,
		ErrCustomerNotFound:                     http.StatusBadRequest,
		ErrPlanNotFound:                         http.StatusBadRequest,
		ErrCardNotFound:                         http.StatusBadRequest,
	}
)

func getErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	return http.StatusInternalServerError
}
