package alarms

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrMaxAlarmsLimitReached: http.StatusBadRequest,
		ErrIntervalTooSmall:      http.StatusBadRequest,
		ErrRegionNotFound:        http.StatusBadRequest,
	}
)
