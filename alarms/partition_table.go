package alarms

import (
	"fmt"
	"time"
)

const (
	// RotateAfterHours defines how long to wait before sub tables are rotated away
	RotateAfterHours = 30 * 24 // 30 days
	// NearingTomorrowHours defines how many hours before tomorrow triggers table partitioning
	NearingTomorrowHours = 1
)

// PartitionTable creates a new result sub table if needed
func (s *Service) PartitionTable(parentTableName string, now time.Time) error {
	today := time.Date(
		now.UTC().Year(), now.UTC().Month(), now.UTC().Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	tomorrow := today.Add(24 * time.Hour)

	// Generate a new sub table name for today
	todaySubTableName := getSubtableName(parentTableName, today)

	// If a sub table for today doesn't exist, create it
	if !s.db.HasTable(todaySubTableName) {
		resultSubTable, err := s.createSubTable(
			parentTableName,
			todaySubTableName,
			today,
			tomorrow,
		)
		if err != nil {
			return err
		}
		logger.Infof("Created new result sub table: %s", resultSubTable.Name)
	}

	// If we are not nearing tomorrow yet, just return
	// (less than NearingTomorrowHours hours from now)
	if now.Add(NearingTomorrowHours * time.Hour).Before(tomorrow) {
		return nil
	}

	// Generate a new sub table name for tomorrow
	tomorrowSubTableName := getSubtableName(parentTableName, tomorrow)

	// If a sub table for tomorrow doesn't exist, create it
	if !s.db.HasTable(tomorrowSubTableName) {
		dayAfterTomorrow := tomorrow.Add(24 * time.Hour)
		resultSubTable, err := s.createSubTable(
			parentTableName,
			tomorrowSubTableName,
			tomorrow,
			dayAfterTomorrow,
		)
		if err != nil {
			return err
		}
		logger.Infof("Created new result sub table: %s", resultSubTable.Name)
	}

	return nil
}

// RotateSubTables deletes sub tables older than RotateAfterHours hours
func (s *Service) RotateSubTables() error {
	var (
		err             error
		rotateAfterDate = time.Now().Add(
			-time.Duration(RotateAfterHours) * time.Hour,
		)
		resultSubTables []*ResultSubTable
	)

	// Fetch result sub table records we want to rotate
	err = s.db.Where("created_at < ?", rotateAfterDate.UTC()).
		Find(&resultSubTables).Error
	if err != nil {
		return err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Delete old result sub tables
	for _, resultSubTable := range resultSubTables {
		if err := tx.DropTable(resultSubTable.Name).Error; err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}
		logger.Infof("Deleted result sub table: %s", resultSubTable.Name)
	}

	// Delete old sub table records
	err = s.db.Where("created_at < ?", rotateAfterDate.UTC()).
		Delete(new(ResultSubTable)).Error
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

// getSubtableName generates a new sub table name for a given date
func getSubtableName(parentTableName string, now time.Time) string {
	return fmt.Sprintf("%s_%s", parentTableName, now.UTC().Format("2006_01_02"))
}

// createSubTable creates a new result sub table inheriting from the parent
// table with a check constraint to limit span of the data to a period of time
func (s *Service) createSubTable(parentTableName, subTableName string, from, to time.Time) (*ResultSubTable, error) {
	// Begin a transaction
	tx := s.db.Begin()

	// Let's create a new sub table, since it doesn't exist yet
	subTable := &Result{Table: subTableName}
	if err := tx.CreateTable(subTable).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	var sql string

	// Add CHECK CONSTRAINT to limit the sub table's span to just one day
	sql = fmt.Sprintf(
		"ALTER TABLE %s ADD CONSTRAINT timestamp_check CHECK (timestamp >= '%s' AND timestamp < '%s') NO INHERIT",
		subTableName,
		from.UTC().Format("2006-01-02T15:04:05Z"),
		to.UTC().Format("2006-01-02T15:04:05Z"),
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
	resultSubTable := &ResultSubTable{Name: subTableName}
	if err := tx.Create(resultSubTable).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return resultSubTable, nil
}
