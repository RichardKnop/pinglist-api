package notifications

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface
	FindEndpointByUserIDAndApplicationARN(userID uint, applicationARN string) (*Endpoint, error)
	PublishMessage(endpointARN, msg string, opt map[string]interface{}) (string, error)

	// Needed for the newRoutes to be able to register handlers
	registerDeviceHandler(w http.ResponseWriter, r *http.Request)
}
