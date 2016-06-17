package health

import (
	"net/http"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Needed for the newRoutes to be able to register handlers
	healthcheck(w http.ResponseWriter, r *http.Request)
}
