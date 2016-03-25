package subscriptions

import (
	"errors"
	"net/http"
	"net/http/httputil"

	stripe "github.com/stripe/stripe-go"
)

var (
	// ErrStripeEventAlreadyRecevied ...
	ErrStripeEventAlreadyRecevied = errors.New("Stripe event already received")
)

// createStripeEventLog logs a Stripe event in an idempotent way
func (s *Service) createStripeEventLog(stripeEvent *stripe.Event, r *http.Request) (*StripeEventLog, error) {
	// Idempotency check
	notFound := s.db.First(new(StripeEventLog), "event_id = ?", stripeEvent.ID).RecordNotFound()
	if !notFound {
		return nil, ErrStripeEventAlreadyRecevied
	}

	// Get request dump including body (so we can see the payload in the event log table)
	requestDump, err := httputil.DumpRequest(
		r,
		true, // include body
	)
	if err != nil {
		return nil, err
	}

	// Save the event data into our log table
	stripeEventLog := NewStripeEventLog(
		stripeEvent.ID,
		stripeEvent.Type,
		string(requestDump),
	)
	if err := s.db.Create(stripeEventLog).Error; err != nil {
		return nil, err
	}

	return stripeEventLog, nil
}
