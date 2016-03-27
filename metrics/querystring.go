package metrics

import (
	"net/http"
	"time"
)

// GetParamsFromQueryString parses querystring and returns params
func GetParamsFromQueryString(r *http.Request) (string, *time.Time, *time.Time, error) {
	var (
		dateTrunc string
		from, to  *time.Time
	)

	// Get "date_trunc" param from the querystring
	if r.URL.Query().Get("date_trunc") != "" {
		dateTrunc = r.URL.Query().Get("date_trunc")
	}

	// Get "from" param from the querystring
	if r.URL.Query().Get("from") != "" {
		t, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
		if err != nil {
			return "", nil, nil, err
		}
		from = &t
	}

	// Get "to" param from the querystring
	if r.URL.Query().Get("to") != "" {
		t, err := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
		if err != nil {
			return "", nil, nil, err
		}
		to = &t
	}

	return dateTrunc, from, to, nil
}
