package metrics

import (
	"time"

	"github.com/jinzhu/gorm"
)

// LogRequestTime logs request time metric
func (s *Service) LogRequestTime(timestamp time.Time, referenceID uint, value int64) error {
	requestTimeRecord := newRequestTime(
		getSubTableName(RequestTimeParentTableName, timestamp),
		referenceID,
		timestamp,
		value,
	)
	return s.db.Create(requestTimeRecord).Error
}

// PaginatedRequestTimesCount returns a total count of request time records
func (s *Service) PaginatedRequestTimesCount(referenceID uint) (int, error) {
	var count int
	if err := s.paginatedRequestTimesQuery(referenceID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindPaginatedRequestTimes returns paginated request time records
// Records can optionally be filtered by a referenceID
func (s *Service) FindPaginatedRequestTimes(offset, limit int, orderBy string, referenceID uint) ([]*RequestTime, error) {
	var requestTimes []*RequestTime

	// Get the pagination query
	requestTimesQuery := s.paginatedRequestTimesQuery(referenceID)

	// Default ordering
	if orderBy == "" {
		orderBy = "timestamp"
	}

	// Retrieve paginated results from the database
	err := requestTimesQuery.Offset(offset).Limit(limit).Order(orderBy).
		Find(&requestTimes).Error
	if err != nil {
		return requestTimes, err
	}

	return requestTimes, nil
}

// paginatedRequestTimesQuery returns a db query for paginated request time records
func (s *Service) paginatedRequestTimesQuery(referenceID uint) *gorm.DB {
	// Basic query
	requestTimesQuery := s.db.Model(new(RequestTime))

	// Optionally filter by reference ID
	if referenceID > 0 {
		requestTimesQuery = requestTimesQuery.Where("reference_id = ?", referenceID)
	}

	return requestTimesQuery
}
