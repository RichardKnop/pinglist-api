package notifications

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrEndpointNotFound ...
	ErrEndpointNotFound = errors.New("Endpoint not found")
)

// FindEndpointByUserIDAndApplicationARN looks up an endoint by user ID and
// platform application ARN and returns it
func (s *Service) FindEndpointByUserIDAndApplicationARN(userID uint, applicationARN string) (*Endpoint, error) {
	// Fetch the endpoint from the database
	endpoint := new(Endpoint)
	notFound := s.db.Where(Endpoint{
		UserID:         util.PositiveIntOrNull(int64(userID)),
		ApplicationARN: applicationARN,
	}).Preload("User").First(endpoint).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrEndpointNotFound
	}

	return endpoint, nil
}

// createOrUpdateEndpoint creates or updates a mobile application endpoint
func (s *Service) createOrUpdateEndpoint(user *accounts.User, applicationARN, deviceToken string) (*Endpoint, error) {
	// Does this user's device already have an endpoint in our database?
	endpoint, err := s.FindEndpointByUserIDAndApplicationARN(user.ID, applicationARN)
	if err != nil {
		// This should never happen, if it does, abort and return
		if err != ErrEndpointNotFound {
			return nil, err
		}

		// This is a first-time registration, create a new endpoint
		endpointARN, err := s.snsAdapter.CreateEndpoint(
			applicationARN,
			user.OauthUser.Username,
			deviceToken,
		)
		if err != nil {
			return nil, err
		}

		// And store the platform endpoint ARN in our database
		endpoint := NewEndpoint(
			user,
			applicationARN,
			endpointARN,
			deviceToken,
			user.OauthUser.Username, // custom user data
			true, // enabled
		)
		if err := s.db.Create(endpoint).Error; err != nil {
			return nil, err
		}

		// And return
		return endpoint, nil
	}

	// Get endpoint attributes
	endpointAttributes, err := s.snsAdapter.GetEndpointAttributes(endpoint.ARN)
	if err != nil {
		return nil, err
	}

	// If the device token in the endpoint does not match the latest one or
	// get endpoint attributes shows the endpoint as disabled
	if !endpointAttributes.Enabled || endpointAttributes.Token != deviceToken {
		// Set the latest device token and then enable the endpoint
		if err := s.snsAdapter.SetEndpointAttributes(
			endpoint.ARN,
			&EndpointAttributes{
				CustomUserData: endpoint.CustomUserData,
				Enabled:        true,
				Token:          deviceToken,
			},
		); err != nil {
			return nil, err
		}

		// Update the endpoint record in our database
		if err := s.db.Model(endpoint).UpdateColumns(Endpoint{
			DeviceToken: deviceToken,
			Enabled:     true,
			Model:       gorm.Model{UpdatedAt: time.Now()},
		}).Error; err != nil {
			return nil, err
		}
	}

	return endpoint, nil
}
