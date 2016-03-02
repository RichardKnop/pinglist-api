package subscriptions

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListPlansResponse(t *testing.T) {
	// Some mock Plan objects
	plans := []*Plan{new(Plan), new(Plan)}

	// Create list response
	response, err := NewListPlansResponse(plans)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/plans", selfLink.Href)
	}

	// Test embedded alarms
	embeddedPlans, err := response.GetEmbedded("plans")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedPlans)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(PlanResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}
}
