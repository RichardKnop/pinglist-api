package notifications

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
)

const (
	// PlatformIOS ...
	PlatformIOS = "iOS"
	// PlatformAndroid ...
	PlatformAndroid = "Android"
)

var (
	errPlatformNotSupported = errors.New("Platform not supported. Use one of: iOS, Android")
)

// Handles calls to register a new device (POST /v1/devices)
func (s *Service) registerDeviceHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Request body cannot be nil
	if r.Body == nil {
		response.Error(w, "Request body cannot be nil", http.StatusBadRequest)
		return
	}

	// Read the request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the request prototype
	deviceRequest := new(DeviceRequest)
	if err := json.Unmarshal(payload, deviceRequest); err != nil {
		logger.Errorf("Failed to unmarshal device request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the correct platform application ARN
	applicationARN, ok := map[string]string{
		PlatformIOS:     s.cnf.AWS.APNSPlatformApplicationARN,
		PlatformAndroid: s.cnf.AWS.GCMPlatformApplicationARN,
	}[deviceRequest.Platform]
	if !ok {
		response.Error(w, errPlatformNotSupported.Error(), http.StatusBadRequest)
		return
	}

	// Register a new endpoint for this device
	_, err = s.createOrUpdateEndpoint(
		authenticatedUser,
		applicationARN,
		deviceRequest.Token,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 204 no content response
	response.NoContent(w)
}
