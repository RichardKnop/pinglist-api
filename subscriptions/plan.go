package subscriptions

import (
	"errors"
)

var (
	// ErrPlanNotFound ...
	ErrPlanNotFound = errors.New("Plan not found")
)

// FindPlanByID looks up a plan by an ID and returns it
func (s *Service) FindPlanByID(planID uint) (*Plan, error) {
	// Fetch the plan from the database
	plan := new(Plan)
	notFound := s.db.First(plan, planID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPlanNotFound
	}

	return plan, nil
}

// FindPlanByPlanID looks up a plan by a plan ID and returns it
func (s *Service) FindPlanByPlanID(planID string) (*Plan, error) {
	// Fetch the plan from the database
	plan := new(Plan)
	notFound := s.db.Where("plan_id = ?", planID).First(plan).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPlanNotFound
	}

	return plan, nil
}
