package alarms

import (
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

// paginatedResultsCount returns a total count of results
func (s *Service) paginatedResultsCount(alarm *Alarm) (int, error) {
	var count int
	if err := s.paginatedResultsQuery(alarm).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedResults returns paginated result records
// Results can optionally be filtered by alarm
func (s *Service) findPaginatedResults(offset, limit int, orderBy string, alarm *Alarm) ([]*Result, error) {
	var results []*Result

	// Get the pagination query
	resultsQuery := s.paginatedResultsQuery(alarm)

	// Default ordering
	if orderBy == "" {
		orderBy = "id"
	}

	// Retrieve paginated results from the database
	err := resultsQuery.Offset(offset).Limit(limit).Order(orderBy).
		Find(&results).Error
	if err != nil {
		return results, err
	}

	return results, nil
}

// paginatedResultsQuery returns a db query for paginated results
func (s *Service) paginatedResultsQuery(alarm *Alarm) *gorm.DB {
	// Basic query
	resultsQuery := s.db.Model(new(Result))

	// Optionally filter by alarm
	if alarm != nil {
		resultsQuery = resultsQuery.Where(Result{
			AlarmID: util.PositiveIntOrNull(int64(alarm.ID)),
		})
	}

	return resultsQuery
}
