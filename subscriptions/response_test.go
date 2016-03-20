package subscriptions

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListPlansResponse(t *testing.T) {
	// Some mock Plan objects
	plans := []*Plan{new(Plan), new(Plan), new(Plan)}

	// Create list response
	response, err := NewListPlansResponse(
		3,           // count
		1,           // page
		"/v1/plans", // self
		"/v1/plans", // first
		"/v1/plans", // last
		"",          // previous
		"",          // next
		plans,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/plans", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/plans", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/plans", lastLink.Href)
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
	embeddedPlans, err := response.GetEmbedded("plans")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedPlans)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(PlanResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 3, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(3), response.Count)
	assert.Equal(t, uint(1), response.Page)
}

func TestNewListSubscriptionsResponse(t *testing.T) {
	// Some mock Subscription objects
	subscriptions := []*Subscription{
		&Subscription{Plan: new(Plan), Customer: new(Customer)},
		&Subscription{Plan: new(Plan), Customer: new(Customer)},
	}

	// Create list response
	response, err := NewListSubscriptionsResponse(
		10, // count
		2,  // page
		"/v1/subscriptions?page=2", // self
		"/v1/subscriptions?page=1", // first
		"/v1/subscriptions?page=5", // last
		"/v1/subscriptions?page=1", // previous
		"/v1/subscriptions?page=3", // next
		subscriptions,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/subscriptions?page=2", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/subscriptions?page=1", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/subscriptions?page=5", lastLink.Href)
	}

	// Test previous link
	previousLink, err := response.GetLink("prev")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/subscriptions?page=1", previousLink.Href)
	}

	// Test next link
	nextLink, err := response.GetLink("next")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/subscriptions?page=3", nextLink.Href)
	}

	// Test embedded subscriptions
	embeddedSubscriptions, err := response.GetEmbedded("subscriptions")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedSubscriptions)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(SubscriptionResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(10), response.Count)
	assert.Equal(t, uint(2), response.Page)
}
