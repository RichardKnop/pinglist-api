package pagination

import (
	"log"
	"net/http"
	"net/url"
	"testing"

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

	// Test limit <= 0
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=1&limit=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, err = GetPageLimit(r)
	if assert.NotNil(t, err) {
		assert.Equal(t, ErrLimitTooSmall, err)
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
