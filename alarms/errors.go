package alarms

import (
  "net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrMaxAlarmsLimitReached: http.StatusBadRequest,
		ErrRegionNotFound:        http.StatusBadRequest,
		ErrAlarmStateNotFound:    http.StatusBadRequest,
	}
)
