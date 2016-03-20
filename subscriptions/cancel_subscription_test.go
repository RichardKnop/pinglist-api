package subscriptions

import (
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestCancelSubscriptionRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.cancelSubscriptionHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}
