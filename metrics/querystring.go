package metrics

import (
	"errors"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/pagination"
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

	// Get from / to time range params
	from, to, err := pagination.GetFromTo(r)
	if err != nil {
		return "", nil, nil, err
	}

	// Get "date_trunc" param from the querystring
	if r.URL.Query().Get("date_trunc") != "" {
		dateTrunc = r.URL.Query().Get("date_trunc")
		if ok, _ := AllowedDateTruncMap[dateTrunc]; !ok {
			return "", nil, nil, ErrInvalidDateTrunc
		}
	}

	return dateTrunc, from, to, nil
}
