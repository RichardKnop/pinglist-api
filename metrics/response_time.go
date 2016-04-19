package metrics

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// LogResponseTime logs request time metric
func (s *Service) LogResponseTime(timestamp time.Time, referenceID uint, value int64) error {
	ResponseTimeRecord := NewResponseTime(
		getSubTableName(ResponseTimeParentTableName, timestamp),
		referenceID,
		timestamp,
		value,
	)
	return s.db.Create(ResponseTimeRecord).Error
}

// ResponseTimesCount returns a total count of response time records
func (s *Service) ResponseTimesCount(referenceID int, dateTrunc string, from, to *time.Time) (int, error) {
	var count int

	// Get the pagination query
	query := s.responseTimesQuery(referenceID, from, to)

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

// FindPaginatedResponseTimes returns paginated request time records
// Records can optionally be filtered by:
// - reference_id
// - date_trunc (day, hour etc)
// - from
// - to
func (s *Service) FindPaginatedResponseTimes(offset, limit int, orderBy string, referenceID int, dateTrunc string, from, to *time.Time) ([]*ResponseTime, error) {
	var ResponseTimes []*ResponseTime

	// Get the pagination query
	query := s.responseTimesQuery(referenceID, from, to)

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
		// ERROR:  column "metrics_response_times.timestamp" must appear in the GROUP BY clause or be used in an aggregate function
		orderBy = strings.Replace(orderBy, "timestamp", "t", 1)
	}

	// Offset and limit
	query = query.Offset(offset).Limit(limit).Order(orderBy)

	// In case we are not aggregating results, we can just use query.Find
	if dateTrunc == "" {
		if err := query.Find(&ResponseTimes).Error; err != nil {
			return ResponseTimes, err
		}
		return ResponseTimes, nil
	}

	// We are aggregating results, therefor it gets more complicated
	rows, err := query.Rows()
	if err != nil {
		return ResponseTimes, err
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
			return ResponseTimes, err
		}

		// Append correct object to our return slice
		ResponseTimes = append(ResponseTimes, &ResponseTime{
			Timestamp: timestamp,
			Value:     int64(value),
		})
	}

	return ResponseTimes, nil
}

// responseTimesQuery returns a common part of db query for
// fetching response time records
func (s *Service) responseTimesQuery(referenceID int, from, to *time.Time) *gorm.DB {
	// Basic query
	query := s.db.Model(new(ResponseTime))

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
