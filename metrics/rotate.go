package metrics

import (
  "time"
)

// RotateAfterHours defines how long to wait before sub tables are rotated away
const RotateAfterHours = 30 * 24 // 30 days

// RotateSubTables deletes sub tables older than RotateAfterHours hours
func (s *Service) RotateSubTables() error {
	var (
		err             error
		rotateAfterDate = time.Now().Add(
			-time.Duration(RotateAfterHours) * time.Hour,
		)
		subTables []*SubTable
	)

	// Fetch result sub table records we want to rotate
	err = s.db.Where("created_at < ?", rotateAfterDate.UTC()).
		Find(&subTables).Error
	if err != nil {
		return err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Delete old result sub tables
	for _, subTable := range subTables {
		if err := tx.DropTable(subTable.Name).Error; err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}
		logger.Infof("Deleted sub table: %s", subTable.Name)
	}

	// Delete old sub table records
	err = s.db.Where("created_at < ?", rotateAfterDate.UTC()).
		Delete(new(SubTable)).Error
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
