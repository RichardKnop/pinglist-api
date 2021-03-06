package alarms

import (
	"reflect"
	"testing"

	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/stretchr/testify/assert"
)

func TestNewListRegionsResponse(t *testing.T) {
	// Some mock Region objects
	regions := []*Region{new(Region), new(Region), new(Region)}

	// Create list response
	response, err := NewListRegionsResponse(
		3,             // count
		1,             // page
		"/v1/regions", // self
		"/v1/regions", // first
		"/v1/regions", // last
		"",            // previous
		"",            // next
		regions,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/regions", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/regions", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/regions", lastLink.Href)
	}

	// Test previous link
	previousLink, err := response.GetLink("prev")
	if assert.Nil(t, err) {
		assert.Equal(t, "", previousLink.Href)
	}

	// Test next link
	nextLink, err := response.GetLink("next")
	if assert.Nil(t, err) {
		assert.Equal(t, "", nextLink.Href)
	}

	// Test embedded alarms
	embeddedPlans, err := response.GetEmbedded("regions")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedPlans)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(RegionResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 3, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(3), response.Count)
	assert.Equal(t, uint(1), response.Page)
}

func TestNewListAlarmsResponse(t *testing.T) {
	// Some mock Alarm objects
	alarms := []*Alarm{new(Alarm), new(Alarm)}

	// Create list response
	response, err := NewListAlarmsResponse(
		10,                  // count
		2,                   // page
		"/v1/alarms?page=2", // self
		"/v1/alarms?page=1", // first
		"/v1/alarms?page=5", // last
		"/v1/alarms?page=1", // previous
		"/v1/alarms?page=3", // next
		alarms,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms?page=2", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms?page=1", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms?page=5", lastLink.Href)
	}

	// Test previous link
	previousLink, err := response.GetLink("prev")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms?page=1", previousLink.Href)
	}

	// Test next link
	nextLink, err := response.GetLink("next")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms?page=3", nextLink.Href)
	}

	// Test embedded alarms
	embeddedAlarms, err := response.GetEmbedded("alarms")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedAlarms)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(AlarmResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(10), response.Count)
	assert.Equal(t, uint(2), response.Page)
}

func TestNewListIncidentsResponse(t *testing.T) {
	// Some mock Incident objects
	incidents := []*Incident{new(Incident), new(Incident)}

	// Create list response
	response, err := NewListIncidentsResponse(
		10, // count
		2,  // page
		"/v1/alarms/1/incidents?page=2", // self
		"/v1/alarms/1/incidents?page=1", // first
		"/v1/alarms/1/incidents?page=5", // last
		"/v1/alarms/1/incidents?page=1", // previous
		"/v1/alarms/1/incidents?page=3", // next
		incidents,
	)
	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/incidents?page=2", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/incidents?page=1", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/incidents?page=5", lastLink.Href)
	}

	// Test previous link
	previousLink, err := response.GetLink("prev")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/incidents?page=1", previousLink.Href)
	}

	// Test next link
	nextLink, err := response.GetLink("next")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/incidents?page=3", nextLink.Href)
	}

	// Test embedded incidents
	embeddedIncidents, err := response.GetEmbedded("incidents")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedIncidents)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(IncidentResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(10), response.Count)
	assert.Equal(t, uint(2), response.Page)
}

func TestNewListResponseTimesResponse(t *testing.T) {
	// Some mock metrics.ResponseTime objects
	responseTimes := []*metrics.ResponseTime{
		&metrics.ResponseTime{Value: 2},
		&metrics.ResponseTime{Value: 3},
	}

	// Mock aggregate incident type counts
	incidentTypeCounts := map[string]int{
		incidenttypes.Slow:    0,
		incidenttypes.Timeout: 1,
		incidenttypes.BadCode: 2,
		incidenttypes.Other:   3,
	}

	// Create list response
	response, err := NewListResponseTimesResponse(
		10, // count
		2,  // page
		"/v1/alarms/1/response-times?page=2", // self
		"/v1/alarms/1/response-times?page=1", // first
		"/v1/alarms/1/response-times?page=5", // last
		"/v1/alarms/1/response-times?page=1", // previous
		"/v1/alarms/1/response-times?page=3", // next
		responseTimes,
		99.99, // uptime
		incidentTypeCounts,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/response-times?page=2", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/response-times?page=1", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/response-times?page=5", lastLink.Href)
	}

	// Test previous link
	previousLink, err := response.GetLink("prev")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/response-times?page=1", previousLink.Href)
	}

	// Test next link
	nextLink, err := response.GetLink("next")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/alarms/1/response-times?page=3", nextLink.Href)
	}

	// Test embedded response times
	embeddedResponseTimes, err := response.GetEmbedded("response_times")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedResponseTimes)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(metrics.MetricResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, 99.99, response.Uptime)
	assert.Equal(t, 2.5, response.Average)
	assert.Equal(t, uint(10), response.Count)
	assert.Equal(t, uint(2), response.Page)
}
