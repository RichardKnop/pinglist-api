package subscriptions

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestListCardsRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.listCardsHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *SubscriptionsTestSuite) TestListCards() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		"http://1.2.3.4/v1/cards",
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_cards", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[0])
	suite.mockUserFiltering(suite.users[0])

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	var cards []*Card
	err = suite.db.Preload("Customer.User").Order("id").Find(&cards).Error
	assert.NoError(suite.T(), err, "Fetching data failed")

	cardResponses := make([]*CardResponse, len(cards))
	for i, card := range cards {
		cardResponse, err := NewCardResponse(card)
		assert.NoError(suite.T(), err, "Creating response object failed")
		cardResponses[i] = cardResponse
	}

	expected := &ListCardsResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: "/v1/cards",
				},
				"first": &jsonhal.Link{
					Href: "/v1/cards?page=1",
				},
				"last": &jsonhal.Link{
					Href: "/v1/cards?page=1",
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"cards": jsonhal.Embedded(cardResponses),
			},
		},
		Count: 4,
		Page:  1,
	}

	expectedJSON, err := json.Marshal(expected)
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
		)
	}
}
