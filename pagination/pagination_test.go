package pagination

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPageLimit(t *testing.T) {
	var (
		page  int
		limit int
		r     *http.Request
		err   error
	)

	// Test default values
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, err = GetPageLimit(r)
	if assert.Nil(t, err) {
		assert.Equal(t, 1, page)
		assert.Equal(t, 25, limit)
	}

	// Test page <= 0
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, err = GetPageLimit(r)
	if assert.NotNil(t, err) {
		assert.Equal(t, ErrPageTooSmall, err)
	}

	// Test limit too small
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=1&limit=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, err = GetPageLimit(r)
	if assert.NotNil(t, err) {
		assert.Equal(t, ErrLimitTooSmall, err)
	}

	// Test limit too big
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=1&limit=1000", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, err = GetPageLimit(r)
	if assert.NotNil(t, err) {
		assert.Equal(t, ErrLimitTooBig, err)
	}

	// Test valid page and limit
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=10&limit=50", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, err = GetPageLimit(r)
	if assert.Nil(t, err) {
		assert.Equal(t, 10, page)
		assert.Equal(t, 50, limit)
	}
}

func TestGetPaginationLinks(t *testing.T) {
	var (
		urlObject *url.URL
		first     string
		last      string
		previous  string
		next      string
		err       error
	)

	// Test with both absolute and relative URI
	for _, testURL := range []string{
		"https://foo.bar/foobar?hello=world",
		"/foobar?hello=world",
	} {
		// Test URL object
		urlObject, err = url.Parse(testURL)
		if err != nil {
			log.Fatal(err)
		}

		// Test with zero results
		first, last, previous, next, err = GetPaginationLinks(
			urlObject,
			0, // count
			1, // page
			2, // limit
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "/foobar?hello=world&page=1", first)
			assert.Equal(t, "/foobar?hello=world&page=1", last)
			assert.Equal(t, "", previous)
			assert.Equal(t, "", next)
		}

		// Test first page
		first, last, previous, next, err = GetPaginationLinks(
			urlObject,
			10, // count
			1,  // page
			2,  // limit
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "/foobar?hello=world&page=1", first)
			assert.Equal(t, "/foobar?hello=world&page=5", last)
			assert.Equal(t, "", previous)
			assert.Equal(t, "/foobar?hello=world&page=2", next)
		}

		// Test middle page
		first, last, previous, next, err = GetPaginationLinks(
			urlObject,
			10, // count
			2,  // page
			2,  // limit
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "/foobar?hello=world&page=1", first)
			assert.Equal(t, "/foobar?hello=world&page=5", last)
			assert.Equal(t, "/foobar?hello=world&page=1", previous)
			assert.Equal(t, "/foobar?hello=world&page=3", next)
		}

		// Test last page
		first, last, previous, next, err = GetPaginationLinks(
			urlObject,
			10, // count
			5,  // page
			2,  // limit
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "/foobar?hello=world&page=1", first)
			assert.Equal(t, "/foobar?hello=world&page=5", last)
			assert.Equal(t, "/foobar?hello=world&page=4", previous)
			assert.Equal(t, "", next)
		}

		// Test page too big
		first, last, previous, next, err = GetPaginationLinks(
			urlObject,
			10, // count
			6,  // page
			2,  // limit
		)
		if assert.NotNil(t, err) {
			assert.Equal(t, ErrPageTooBig, err)
		}

		// Test when limit > count
		first, last, previous, next, err = GetPaginationLinks(
			urlObject,
			10, // count
			1,  // page
			12, // limit
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "/foobar?hello=world&page=1", first)
			assert.Equal(t, "/foobar?hello=world&page=1", last)
			assert.Equal(t, "", previous)
			assert.Equal(t, "", next)
		}
	}
}

func TestGetOffsetForPage(t *testing.T) {
	var offset int

	// First page offset should be zero
	offset = GetOffsetForPage(
		10, // count
		1,  // page
		2,  // limit
	)
	assert.Equal(t, 0, offset)

	// Second page offset should be 2
	offset = GetOffsetForPage(
		10, // count
		2,  // page
		2,  // limit
	)
	assert.Equal(t, 2, offset)

	// Last page offset should be 8
	offset = GetOffsetForPage(
		10, // count
		5,  // page
		2,  // limit
	)
	assert.Equal(t, 8, offset)
}

func TestGetFromTo(t *testing.T) {
	var (
		r        *http.Request
		from, to *time.Time
		err      error
	)

	// Let's try without any params first
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foobar", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	from, to, err = GetFromTo(r)

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
	from, to, err = GetFromTo(r)

	// Check error is nil and correct values were returned
	if assert.NoError(t, err) {
		assert.Equal(t, expectedFrom, from.UTC().Format(time.RFC3339))
		assert.Equal(t, expectedTo, to.UTC().Format(time.RFC3339))
	}
}
