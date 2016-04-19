package alarms

import (
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/util"
)

// GetTimeRangeParamsFromQueryString parses querystring and returns from / to params
func GetTimeRangeParamsFromQueryString(r *http.Request) (*time.Time, *time.Time, error) {
	var (
		from, to *time.Time
	)

	// Get "from" param from the querystring
	if r.URL.Query().Get("from") != "" {
		t, err := util.ParseTimestamp(r.URL.Query().Get("from"))
		if err != nil {
			return nil, nil, err
		}
		from = &t
	}

	// Get "to" param from the querystring
	if r.URL.Query().Get("to") != "" {
		t, err := util.ParseTimestamp(r.URL.Query().Get("to"))
		if err != nil {
			return nil, nil, err
		}
		to = &t
	}

	return from, to, nil
}
