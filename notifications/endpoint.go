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
	var (
		endpoint           *Endpoint
		endpointAttributes *EndpointAttributes
		err                error
	)

	// Does this user's device already have an endpoint in our database?
	endpoint, err = s.FindEndpointByUserIDAndApplicationARN(user.ID, applicationARN)
	if err != nil {
		// This should never happen, if it does, abort and return
		if err != ErrEndpointNotFound {
			return nil, err
		}

		return s.createEndpoint(user, applicationARN, deviceToken)
	}

	// Get endpoint attributes
	endpointAttributes, err = s.snsAdapter.GetEndpointAttributes(endpoint.ARN)
	if err != nil {
		// Not found? Perhaps the endpoint was deleted
		return s.createEndpoint(user, applicationARN, deviceToken)
	}

	// If the device token in the endpoint does not match the latest one or
	// get endpoint attributes shows the endpoint as disabled
	if !endpointAttributes.Enabled || endpointAttributes.Token != deviceToken {
		// Set the latest device token and then enable the endpoint
		if err := s.snsAdapter.SetEndpointAttributes(
			endpoint.ARN,
			&EndpointAttributes{
				Token:   deviceToken,
				Enabled: true,
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

func (s *Service) createEndpoint(user *accounts.User, applicationARN, deviceToken string) (*Endpoint, error) {
	// This is a first-time registration, create a new endpoint
	endpointARN, err := s.snsAdapter.CreateEndpoint(
		applicationARN,
		deviceToken,
	)
	if err != nil {
		return nil, err
	}

	// Begin a transaction
	tx := s.db.Begin()

	var endpoint = new(Endpoint)

	// Grab the first matching endpoint or create a new one
	if err := s.db.Where(map[string]interface{}{
		"user_id":         user.ID,
		"application_arn": applicationARN,
	}).FirstOrCreate(&endpoint).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Update arn, device token and set enabled to true
	if err := s.db.Model(endpoint).UpdateColumns(map[string]interface{}{
		"arn":          endpointARN,
		"device_token": deviceToken,
		"enabled":      true,
		"updated_at":   time.Now(),
	}).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return endpoint, nil
}
