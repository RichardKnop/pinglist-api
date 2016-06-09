package alarms

import (
	"fmt"
	"time"
)

func (s *Service) sendNewIncidentPushNotification(alarm *Alarm, incident *Incident) {
	now := time.Now()

	// Find SNS endpoint
	endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
		uint(alarm.UserID.Int64),
		s.cnf.AWS.APNSPlatformApplicationARN,
	)
	if err != nil {
		logger.Errorf("Find endpoint by user ID and application ARN error: %s", err.Error())
		return
	}

	// Send push notification
	_, err = s.notificationsService.PublishMessage(
		endpoint.ARN,
		fmt.Sprintf(
			newIncidentPushNotificationTemplates[incident.IncidentTypeID.String],
			alarm.EndpointURL,
		),
		map[string]interface{}{},
	)
	if err != nil {
		logger.Errorf("Publish message error: %s", err.Error())
		return
	}

	// Increment push notifications counter
	if err := s.updateNotificationCounterIncrementPush(
		incident.Alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	); err != nil {
		logger.Errorf("Increment push notifications counter error: %s", err)
	}
}

func (s *Service) sendNewIncidentEmail(incident *Incident) {
	now := time.Now()

	// Get alarm limits
	alarmLimits := s.getAlarmLimits(incident.Alarm.User)

	// Fetch the notification counter
	notificationCounter, err := s.findNotificationCounter(
		incident.Alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	)
	if err != nil {
		logger.Errorf("Find notification counter error: %s", err.Error())
		return
	}
	if !alarmLimits.unlimitedEmails && notificationCounter.Email > alarmLimits.maxEmailsPerInterval {
		logger.Errorf(
			"User has already reached maximum emails per month limit: %d/%d",
			notificationCounter.Email,
			alarmLimits.maxEmailsPerInterval,
		)
		return
	}

	newIncidentEmail := s.emailFactory.NewIncidentEmail(incident)

	// Send the email
	if err := s.emailService.Send(newIncidentEmail); err != nil {
		logger.Errorf("Send email error: %s", err)
		return
	}

	// Increment email notifications counter
	if err := s.updateNotificationCounterIncrementEmail(
		incident.Alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	); err != nil {
		logger.Errorf("Increment email notifications counter error: %s", err)
	}
}

func (s *Service) sendNewIncidentSlackMessage(alarm *Alarm, incident *Incident) {
	now := time.Now()

	// Get alarm limits
	alarmLimits := s.getAlarmLimits(alarm.User)

	if !alarmLimits.slackAlerts {
		return
	}

	newIncidentMessage := s.slackFactory.NewIncidentMessage(incident)

	// Send slack message
	if err := s.slackAdapter.SendMessage(
		alarm.User.SlackIncomingWebhook.String,
		alarm.User.SlackChannel.String,
		slackNotificationsUsername,
		slackNotificationsEmoji,
		newIncidentMessage,
	); err != nil {
		logger.Errorf("Send slack message error: %s", err)
		return
	}

	// Increment slack notifications counter
	if err := s.updateNotificationCounterIncrementSlack(
		incident.Alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	); err != nil {
		logger.Errorf("Increment slack notifications counter error: %s", err)
	}
}

func (s *Service) sendIncidentsResolvedPushNotification(alarm *Alarm) {
	now := time.Now()

	// Find SNS endpoint
	endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		s.cnf.AWS.APNSPlatformApplicationARN,
	)
	if err != nil {
		logger.Errorf("Find endpoint by user ID and application ARN error: %s", err.Error())
		return
	}

	// Send push notification
	_, err = s.notificationsService.PublishMessage(
		endpoint.ARN,
		fmt.Sprintf(incidentsResolvedPushNotificationTemplate, alarm.EndpointURL),
		map[string]interface{}{},
	)
	if err != nil {
		logger.Errorf("Publish message error: %s", err.Error())
		return
	}

	// Increment push notifications counter
	if err := s.updateNotificationCounterIncrementPush(
		alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	); err != nil {
		logger.Errorf("Increment push notifications counter error: %s", err)
	}
}

func (s *Service) sendIncidentsResolvedEmail(alarm *Alarm) {
	now := time.Now()

	// Get alarm limits
	alarmLimits := s.getAlarmLimits(alarm.User)

	// Fetch the notification counter
	notificationCounter, err := s.findNotificationCounter(
		alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	)
	if err != nil {
		logger.Errorf("Find notification counter error: %s", err.Error())
		return
	}
	if !alarmLimits.unlimitedEmails && notificationCounter.Email > alarmLimits.maxEmailsPerInterval {
		logger.Errorf(
			"User has already reached maximum emails per month limit: %d/%d",
			notificationCounter.Email,
			alarmLimits.maxEmailsPerInterval,
		)
		return
	}

	alarmUpEmail := s.emailFactory.NewIncidentsResolvedEmail(alarm)

	// Send the email
	if err := s.emailService.Send(alarmUpEmail); err != nil {
		logger.Errorf("Send email error: %s", err)
		return
	}

	// Increment email notifications counter
	if err := s.updateNotificationCounterIncrementEmail(
		alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	); err != nil {
		logger.Errorf("Increment email notifications counter error: %s", err)
	}
}

func (s *Service) sendIncidentsResolvedSlackMessage(alarm *Alarm) {
	now := time.Now()

	// Get alarm limits
	alarmLimits := s.getAlarmLimits(alarm.User)

	if !alarmLimits.slackAlerts {
		return
	}

	newIncidentMessage := s.slackFactory.NewIncidentsResolvedMessage(alarm)

	// Send slack message
	if err := s.slackAdapter.SendMessage(
		alarm.User.SlackIncomingWebhook.String,
		alarm.User.SlackChannel.String,
		slackNotificationsUsername,
		slackNotificationsEmoji,
		newIncidentMessage,
	); err != nil {
		logger.Errorf("Send slack message error: %s", err)
		return
	}

	// Increment slack notifications counter
	if err := s.updateNotificationCounterIncrementSlack(
		alarm.User.ID,
		uint(now.Year()),
		uint(now.Month()),
	); err != nil {
		logger.Errorf("Increment slack notifications counter error: %s", err)
	}
}
