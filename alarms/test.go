package alarms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

// testServer returns a mock HTTP server for unit testing
func testServer(resp *http.Response) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(resp.StatusCode)
		for key := range resp.Header {
			w.Header().Set(key, resp.Header.Get(key))
		}
		fmt.Fprintln(w, resp.Body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	return server, &http.Client{Transport: transport}
}

// testServer returns a mock HTTP server for unit testing a timeout
func testServerTimeout() (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Nanosecond)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	return server, &http.Client{Transport: transport, Timeout: 1 * time.Nanosecond}
}
