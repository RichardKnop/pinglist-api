package subscriptions

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrUserCanOnlyHaveOneActiveSubscription:         http.StatusBadRequest,
		ErrCustomerNotFound:                             http.StatusBadRequest,
		ErrPlanNotFound:                                 http.StatusBadRequest,
		ErrCardNotFound:                                 http.StatusBadRequest,
		ErrCardCanOnlyBeDeletedFromCancelledSubsription: http.StatusBadRequest,
	}
)
