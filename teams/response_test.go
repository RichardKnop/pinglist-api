package teams

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListTeamsResponse(t *testing.T) {
	// Some mock Team objects
	teams := []*Team{new(Team), new(Team), new(Team)}

	// Create list response
	response, err := NewListTeamsResponse(
		3,           // count
		1,           // page
		"/v1/teams", // self
		"/v1/teams", // first
		"/v1/teams", // last
		"",          // previous
		"",          // next
		teams,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/teams", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/teams", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/teams", lastLink.Href)
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

	// Test embedded teams
	embeddedTeams, err := response.GetEmbedded("teams")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedTeams)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(TeamResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 3, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(3), response.Count)
	assert.Equal(t, uint(1), response.Page)
}
