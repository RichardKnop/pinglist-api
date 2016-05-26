package alarms

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrMaxAlarmsLimitReached: http.StatusBadRequest,
		ErrMaxResponseTimeTooBig: http.StatusBadRequest,
		ErrRegionNotFound:        http.StatusBadRequest,
	}
)

func getErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	switch err.(type) {
	case ErrIntervalTooSmall:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
