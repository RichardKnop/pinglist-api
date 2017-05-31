package metrics

import (
	"fmt"
	"time"

	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/util"
)

// PartitionResponseTime creates a new request time sub table if needed
func (s *Service) PartitionResponseTime(parentTableName string, now time.Time) error {
	var (
		today = time.Date(
			now.UTC().Year(), now.UTC().Month(), now.UTC().Day(),
			0, 0, 0, 0,
			time.UTC,
		)
		tomorrow = today.Add(24 * time.Hour)
		// Generate a new sub table name for today
		todaySubTableName = getSubTableName(parentTableName, today)
		// Generate a new sub table name for tomorrow
		tomorrowSubTableName = getSubTableName(parentTableName, tomorrow)
	)

	// If a sub table for today doesn't exist, create it
	if !s.db.HasTable(todaySubTableName) {
		subTable, err := s.createResponseTimeSubTable(
			parentTableName,
			todaySubTableName,
			today,
			tomorrow,
		)
		if err != nil {
			return err
		}
		logger.INFO.Printf("Created a new sub table: %s", subTable.Name)
	}

	// If a sub table for tomorrow doesn't exist, create it
	if !s.db.HasTable(tomorrowSubTableName) {
		dayAfterTomorrow := tomorrow.Add(24 * time.Hour)
		subTable, err := s.createResponseTimeSubTable(
			parentTableName,
			tomorrowSubTableName,
			tomorrow,
			dayAfterTomorrow,
		)
		if err != nil {
			return err
		}
		logger.INFO.Printf("Created a new sub table: %s", subTable.Name)
	}

	return nil
}

// createResponseTimeSubTable creates a new request time sub table inheriting
// from the parent table with a check constraint to limit span of the data
// to a period of time
func (s *Service) createResponseTimeSubTable(parentTableName, subTableName string, from, to time.Time) (*SubTable, error) {
	// Begin a transaction
	tx := s.db.Begin()

	// Let's create a new request time sub table
	ResponseTimeTable := &ResponseTime{Table: subTableName}
	if err := tx.CreateTable(ResponseTimeTable).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	var sql string

	// Add CHECK CONSTRAINT to limit the sub table's span to just one day
	sql = fmt.Sprintf(
		"ALTER TABLE %s ADD CONSTRAINT timestamp_check CHECK (timestamp >= '%s' AND timestamp < '%s') NO INHERIT",
		subTableName,
		util.FormatTime(from),
		util.FormatTime(to),
	)
	if err := tx.Exec(sql).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Inherit from the parent table
	sql = fmt.Sprintf(
		"ALTER TABLE %s INHERIT %s",
		subTableName,
		parentTableName,
	)
	if err := tx.Exec(sql).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Keep track of the new sub table
	subTableRecord := &SubTable{ParentTable: parentTableName, Name: subTableName}
	if err := tx.Create(subTableRecord).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return subTableRecord, nil
}
