package pagination

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

var (
	// DefaultLimit ...
	DefaultLimit = 25
	// ErrPageTooSmall ...
	ErrPageTooSmall = errors.New("Page must be a positive number")
	// ErrLimitTooSmall ...
	ErrLimitTooSmall = errors.New("Limit must be a positive number")
	// ErrPageTooBig ...
	ErrPageTooBig = errors.New("Page too big")
)

// GetPageLimit parses querystring and returns page and limit
func GetPageLimit(r *http.Request) (int, int, error) {
	var (
		page  = 1 // default page
		limit = DefaultLimit
		err   error
	)

	// Get page from the querystring
	if r.URL.Query().Get("page") != "" {
		// String to int conversion
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			return 0, 0, err
		}

		// Page must be >= 0
		if page < 1 {
			return 0, 0, ErrPageTooSmall
		}
	}

	// Get limit from the querystring
	if r.URL.Query().Get("limit") != "" {
		// String to int conversion
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			return 0, 0, err
		}

		// Limit must be > 0
		if limit < 1 {
			return 0, 0, ErrLimitTooSmall
		}
	}

	return page, limit, nil
}

// GetPaginationLinks returns links for first, last, previous and next page
func GetPaginationLinks(urlObject *url.URL, count, page, limit int) (string, string, string, string, error) {
	var (
		first    string
		last     string
		previous string
		next     string
		q        url.Values
	)

	// Number of pages
	nuPages := int(math.Ceil(float64(count) / float64(limit)))
	if nuPages < 1 {
		nuPages = 1
	}

	// Page too big
	if page > nuPages {
		return first, last, previous, next, ErrPageTooBig
	}

	// First page
	q = urlObject.Query()
	q.Set("page", fmt.Sprintf("%d", 1))
	first = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())

	// Last page
	q = urlObject.Query()
	q.Set("page", fmt.Sprintf("%d", nuPages))
	last = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())

	// Previous page
	if page > 1 {
		q := urlObject.Query()
		q.Set("page", fmt.Sprintf("%d", page-1))
		previous = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())
	}

	// Next page
	if page < nuPages {
		q := urlObject.Query()
		q.Set("page", fmt.Sprintf("%d", page+1))
		next = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())
	}

	return first, last, previous, next, nil
}

// GetOffsetForPage returns an offset for a page
func GetOffsetForPage(count, page, limit int) int {
	return limit * (page - 1)
}
