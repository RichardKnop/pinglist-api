package metrics

import (
	"errors"
	"net/http"
	"time"
)

// AllowedDateTruncMap ...
var AllowedDateTruncMap = map[string]bool{
	"hour": true,
	"day":  true,
}

// ErrInvalidDateTrunc ...
var ErrInvalidDateTrunc = errors.New("Invalid date_trunc value. Use one of: hour, day")

// GetParamsFromQueryString parses querystring and returns params
func GetParamsFromQueryString(r *http.Request) (string, *time.Time, *time.Time, error) {
	var (
		dateTrunc string
		from, to  *time.Time
	)

	// Get "date_trunc" param from the querystring
	if r.URL.Query().Get("date_trunc") != "" {
		dateTrunc = r.URL.Query().Get("date_trunc")
		if ok, _ := AllowedDateTruncMap[dateTrunc]; !ok {
			return "", nil, nil, ErrInvalidDateTrunc
		}
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
