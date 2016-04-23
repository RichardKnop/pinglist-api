package alarms

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrMaxAlarmsLimitReached: http.StatusBadRequest,
		ErrIntervalTooSmall:      http.StatusBadRequest,
		ErrMaxResponseTimeTooBig: http.StatusBadRequest,
		ErrRegionNotFound:        http.StatusBadRequest,
	}
)
