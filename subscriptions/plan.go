package subscriptions

import (
  "errors"
)

var (
	errPlanNotFound = errors.New("Plan not found")
)

// FindPlanByID looks up a plan by an ID and returns it
func (s *Service) FindPlanByID(planID uint) (*Plan, error) {
	// Fetch the plan from the database
	plan := new(Plan)

	// Not found
	if s.db.First(plan, planID).RecordNotFound() {
		return nil, errPlanNotFound
	}

	return plan, nil
}
