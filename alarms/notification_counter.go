package alarms

import (
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

func (s *Service) findNotificationCounter(userID, year, month uint) (*NotificationCounter, error) {
	// Fetch the notification counter from the database
	notificationCounter := new(NotificationCounter)
	err := s.db.FirstOrCreate(
		notificationCounter,
		NotificationCounter{
			UserID: util.PositiveIntOrNull(int64(userID)),
			Year:   year,
			Month:  month,
		},
	).Error

	if err != nil {
		return nil, err
	}

	return notificationCounter, nil
}

func (s *Service) updateNotificationCounterIncrementEmail(userID, year, month uint) error {
	notificationCounter, err := s.findNotificationCounter(userID, year, month)
	if err != nil {
		return err
	}

	return s.db.Model(notificationCounter).UpdateColumn("email", gorm.Expr("email + ?", 1)).Error
}

func (s *Service) updateNotificationCounterIncrementPush(userID, year, month uint) error {
	notificationCounter, err := s.findNotificationCounter(userID, year, month)
	if err != nil {
		return err
	}

	return s.db.Model(notificationCounter).UpdateColumn("push", gorm.Expr("push + ?", 1)).Error
}

func (s *Service) updateNotificationCounterIncrementSlack(userID, year, month uint) error {
	notificationCounter, err := s.findNotificationCounter(userID, year, month)
	if err != nil {
		return err
	}

	return s.db.Model(notificationCounter).UpdateColumn("slack", gorm.Expr("slack + ?", 1)).Error
}
