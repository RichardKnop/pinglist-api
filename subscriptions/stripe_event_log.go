package subscriptions

import (
	"net/http"
	"net/http/httputil"

	stripe "github.com/stripe/stripe-go"
)

func (s *Service) logStripeEvent(stripeEvent *stripe.Event, r *http.Request) (*StripeEventLog, error) {
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
