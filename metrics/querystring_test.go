package metrics

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetParamsFromQueryString(t *testing.T) {
	var (
		r         *http.Request
		dateTrunc string
		from, to  *time.Time
		err       error
	)

	// Let's try without any params first
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foobar", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	dateTrunc, from, to, err = GetParamsFromQueryString(r)

	// Check error is nil and correct values were returned
	if assert.NoError(t, err) {
		assert.Equal(t, "", dateTrunc)
		assert.Nil(t, from)
		assert.Nil(t, to)
	}

	// Let's try with all params now
	expectedDateTrunc := "day"
	expectedFrom := "2016-02-08T00:00:00Z"
	expectedTo := "2016-03-08T00:00:00Z"
	r, err = http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://1.2.3.4/v1/foobar?date_trunc=%s&from=%s&to=%s",
			expectedDateTrunc,
			expectedFrom,
			expectedTo,
		),
		nil,
	)
	assert.NoError(t, err, "Request setup should not get an error")
	dateTrunc, from, to, err = GetParamsFromQueryString(r)

	// Check error is nil and correct values were returned
	if assert.NoError(t, err) {
		assert.Equal(t, expectedDateTrunc, dateTrunc)
		assert.Equal(t, expectedFrom, from.UTC().Format(time.RFC3339))
		assert.Equal(t, expectedTo, to.UTC().Format(time.RFC3339))
	}

	// Let's try with invalid date_trunc
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foobar?date_trunc=bogus", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	dateTrunc, from, to, err = GetParamsFromQueryString(r)

	// Check correct error and zero values were returned
	if assert.Error(t, err) {
		assert.Equal(t, ErrInvalidDateTrunc, err)
		assert.Equal(t, "", dateTrunc)
		assert.Nil(t, from)
		assert.Nil(t, to)
	}
}
