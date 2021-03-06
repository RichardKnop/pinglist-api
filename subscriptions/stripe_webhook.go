package subscriptions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/response"
	stripe "github.com/stripe/stripe-go"
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
		logger.ERROR.Printf("Failed to unmarshal stripe event: %s", payload)
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
	stripeEventLog, err := s.createStripeEventLog(stripeEvent, r)

	// We have already received and processed this event, just return OK
	if err != nil && err == ErrStripeEventAlreadyRecevied && stripeEventLog != nil && stripeEventLog.Processed {
		response.WriteJSON(w, stripeEvent, http.StatusOK)
		return
	}

	// Something went wrong logging the event
	if err != nil && err != ErrStripeEventAlreadyRecevied {
		logger.ERROR.Print("Failed to log the stripe event")
		logger.ERROR.Print(err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map events we are interested in to respective handlers
	stripeEventHandlerMap := map[string]func(e *stripe.Event) error{
		"customer.subscription.updated":        s.stripeEventCustomerSubscriptionUpdated,
		"customer.subscription.deleted":        s.stripeEventCustomerSubscriptionDeleted,
		"customer.subscription.trial_will_end": s.stripeEventCustomerSubscriptionTrialWillEnd,
	}

	// Get the event handler
	stripeEventHandler, ok := stripeEventHandlerMap[stripeEvent.Type]

	// No handler defined for this event, we are not interested so just return OK
	if !ok {
		response.WriteJSON(w, stripeEvent, http.StatusOK)
		return
	}

	// Try to process the event
	if err := stripeEventHandler(stripeEvent); err != nil {
		logger.ERROR.Printf("Failed to process stripe event: %v", stripeEvent)
		logger.ERROR.Print(err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the processed flag in the event log table
	if err := s.db.Model(stripeEventLog).UpdateColumns(map[string]interface{}{
		"processed":  true,
		"updated_at": time.Now(),
	}).Error; err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write JSON response
	response.WriteJSON(w, stripeEvent, http.StatusOK)
}
