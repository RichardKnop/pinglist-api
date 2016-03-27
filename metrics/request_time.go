package metrics

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// LogRequestTime logs request time metric
func (s *Service) LogRequestTime(timestamp time.Time, referenceID uint, value int64) error {
	requestTimeRecord := NewRequestTime(
		getSubTableName(RequestTimeParentTableName, timestamp),
		referenceID,
		timestamp,
		value,
	)
	return s.db.Create(requestTimeRecord).Error
}

// PaginatedRequestTimesCount returns a total count of request time records
func (s *Service) PaginatedRequestTimesCount(referenceID int, dateTrunc string, from, to *time.Time) (int, error) {
	var count int

	// Get the pagination query
	query := s.paginatedRequestTimesQuery(referenceID, from, to)

	// Are we aggregating data based on some time period (e.g. hourly / daily averages)?
	if dateTrunc != "" {
		query = query.
			Select(fmt.Sprintf("COUNT(DISTINCT(date_trunc('%s', timestamp at time zone 'Z')))", dateTrunc))
		if err := query.Row().Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindPaginatedRequestTimes returns paginated request time records
// Records can optionally be filtered by:
// - reference_id
// - date_trunc (day, hour etc)
// - from
// - to
func (s *Service) FindPaginatedRequestTimes(offset, limit int, orderBy string, referenceID int, dateTrunc string, from, to *time.Time) ([]*RequestTime, error) {
	var requestTimes []*RequestTime

	// Get the pagination query
	query := s.paginatedRequestTimesQuery(referenceID, from, to)

	// Default ordering
	if orderBy == "" {
		orderBy = "timestamp"
	}

	// Are we aggregating data based on some time period (e.g. hourly / daily averages)?
	if dateTrunc != "" {
		query = query.
			Select("date_trunc(?, timestamp at time zone 'Z') t, AVG(value) avg", dateTrunc).
			Group("t")
		// This is needed because if we use "timestamp" in ORDER BY clause,
		// since timestamp is not present in our aggregate function there is an error:
		// ERROR:  column "metrics_request_times.timestamp" must appear in the GROUP BY clause or be used in an aggregate function
		orderBy = strings.Replace(orderBy, "timestamp", "t", 1)
	}

	// Offset and limit
	query = query.Offset(offset).Limit(limit).Order(orderBy)

	// In case we are not aggregating results, we can just use query.Find
	if dateTrunc == "" {
		if err := query.Find(&requestTimes).Error; err != nil {
			return requestTimes, err
		}
		return requestTimes, nil
	}

	// We are aggregating results, therefor it gets more complicated
	rows, err := query.Rows()
	if err != nil {
		return requestTimes, err
	}

	// Iterate over *sql.Rows
	for rows.Next() {
		// Declare vars for copying the data from the row
		var (
			timestamp time.Time
			value     float64
		)

		// Scan the data into our vars
		if err := rows.Scan(&timestamp, &value); err != nil {
			return requestTimes, err
		}

		// Append correct object to our return slice
		requestTimes = append(requestTimes, &RequestTime{
			Timestamp: timestamp,
			Value:     int64(value),
		})
	}

	return requestTimes, nil
}

// paginatedRequestTimesQuery returns a common part of db query for
// paginated request time records
func (s *Service) paginatedRequestTimesQuery(referenceID int, from, to *time.Time) *gorm.DB {
	// Basic query
	query := s.db.Model(new(RequestTime))

	// Optionally filter by reference ID
	if referenceID > 0 {
		query = query.Where("reference_id = ?", referenceID)
	}

	// Optionally limit timestamp to be greater or equal to
	if from != nil {
		query = query.Where("timestamp >= ?", from)
	}

	// Optionally limit timestamp to be less or equal to
	if to != nil {
		query = query.Where("timestamp <= ?", to)
	}

	return query
}
