package subscriptions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/response"
	stripe "github.com/stripe/stripe-go"
	"github.com/jinzhu/gorm"
)

// Handles calls to Stripe webhook (POST /v1/stripe-webhook)
func (s *Service) stripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
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
	stripeEventRequest := new(stripe.Event)
	if err := json.Unmarshal(payload, stripeEventRequest); err != nil {
		logger.Errorf("Failed to unmarshal stripe event: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the event by fetching it from Stripe
	stripeEvent, err := s.stripeAdapter.GetEvent(stripeEventRequest.ID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Log the event data in event log table
	stripeEventLog, err := s.logStripeEvent(stripeEvent, r)
	if err != nil {
		logger.Error("Failed to log the stripe event")
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start with a nil error
	err = nil

	// Process event types we are interested in
	switch stripeEvent.Type {
	case "customer.created":
		err = s.stripeEventCustomerCreated(stripeEvent)
	case "customer.subscription.created":
		err = s.stripeEventCustomerSubscriptionCreated(stripeEvent)
	case "customer.subscription.trial_will_end":
		err = s.stripeEventCustomerSubscriptionTrialWillEnd(stripeEvent)
	case "invoice.created":
		err = s.stripeEventInvoiceCreated(stripeEvent)
	case "charge.succeeded":
		err = s.stripeEventChargeSucceeded(stripeEvent)
	case "invoice.payment_succeeded":
		err = s.stripeEventPaymentSucceeded(stripeEvent)
	case "customer.subscription.updated":
		err = s.stripeEventCustomerSubscriptionUpdated(stripeEvent)
	case "customer.subscription.deleted":
		err = s.stripeEventCustomerSubscriptionDeleted(stripeEvent)
	}

	// An error occured while processing an event sent by Stripe
	if err != nil {
		logger.Errorf("Failed to process stripe event: %v", stripeEvent)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Flag the event as processed in the event log table
	if err := s.db.Model(stripeEventLog).UpdateColumns(StripeEventLog{
		Processed: true,
		Model:     gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write JSON response
	response.WriteJSON(w, stripeEvent, http.StatusOK)
}
