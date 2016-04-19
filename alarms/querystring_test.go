package alarms

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeRangeParamsFromQueryString(t *testing.T) {
	var (
		r        *http.Request
		from, to *time.Time
		err      error
	)

	// Let's try without any params first
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foobar", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	from, to, err = GetTimeRangeParamsFromQueryString(r)

	// Check error is nil and correct values were returned
	if assert.NoError(t, err) {
		assert.Nil(t, from)
		assert.Nil(t, to)
	}

	// Let's try with all params now
	expectedFrom := "2016-02-08T00:00:00Z"
	expectedTo := "2016-03-08T00:00:00Z"
	r, err = http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://1.2.3.4/v1/foobar?from=%s&to=%s",
			expectedFrom,
			expectedTo,
		),
		nil,
	)
	assert.NoError(t, err, "Request setup should not get an error")
	from, to, err = GetTimeRangeParamsFromQueryString(r)

	// Check error is nil and correct values were returned
	if assert.NoError(t, err) {
		assert.Equal(t, expectedFrom, from.UTC().Format(time.RFC3339))
		assert.Equal(t, expectedTo, to.UTC().Format(time.RFC3339))
	}
}
