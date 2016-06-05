package alarms

import "fmt"

func (s *Service) sendNewIncidentPushNotification(alarm *Alarm, incident *Incident) {
	endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
		uint(alarm.UserID.Int64),
		s.cnf.AWS.APNSPlatformApplicationARN,
	)
	if err == nil && endpoint != nil {
		_, err := s.notificationsService.PublishMessage(
			endpoint.ARN,
			fmt.Sprintf(
				newIncidentPushNotificationTemplates[incident.IncidentTypeID.String],
				alarm.EndpointURL,
			),
			map[string]interface{}{},
		)
		if err != nil {
			logger.Errorf("Publish Message Error: %s", err.Error())
		}
	}
}

func (s *Service) sendNewIncidentEmail(incident *Incident) {
	newIncidentEmail := s.emailFactory.NewIncidentEmail(incident)

	// Try to send the new incident email
	if err := s.emailService.Send(newIncidentEmail); err != nil {
		logger.Errorf("Send email error: %s", err)
		return
	}
}

func (s *Service) sendNewIncidentSlackMessage(alarm *Alarm, incident *Incident) {
	// Fetch the user team
	team, _ := s.teamsService.FindTeamByMemberID(alarm.User.ID)

	// Get alarm limits
	alarmLimits := s.getAlarmLimits(team, alarm.User)

	if !alarmLimits.slackAlerts {
		return
	}

	newIncidentMessage := s.slackFactory.NewIncidentMessage(incident)

	if err := s.slackAdapter.SendMessage(
		alarm.User.SlackIncomingWebhook.String,
		alarm.User.SlackChannel.String,
		slackNotificationsUsername,
		slackNotificationsEmoji,
		newIncidentMessage,
	); err != nil {
		logger.Errorf("Send Slack message error: %s", err)
		return
	}
}

func (s *Service) sendIncidentsResolvedPushNotification(alarm *Alarm) {
	endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		s.cnf.AWS.APNSPlatformApplicationARN,
	)
	if err == nil && endpoint != nil {
		_, err := s.notificationsService.PublishMessage(
			endpoint.ARN,
			fmt.Sprintf(incidentsResolvedPushNotificationTemplate, alarm.EndpointURL),
			map[string]interface{}{},
		)
		if err != nil {
			logger.Errorf("Publish Message Error: %s", err.Error())
		}
	}
}

func (s *Service) sendIncidentsResolvedEmail(alarm *Alarm) {
	alarmUpEmail := s.emailFactory.NewIncidentsResolvedEmail(alarm)

	// Try to send the alarm up email email
	if err := s.emailService.Send(alarmUpEmail); err != nil {
		logger.Errorf("Send email error: %s", err)
		return
	}
}

func (s *Service) sendIncidentsResolvedSlackMessage(alarm *Alarm) {
	// Fetch the user team
	team, _ := s.teamsService.FindTeamByMemberID(alarm.User.ID)

	// Get alarm limits
	alarmLimits := s.getAlarmLimits(team, alarm.User)

	if !alarmLimits.slackAlerts {
		return
	}

	newIncidentMessage := s.slackFactory.NewIncidentsResolvedMessage(alarm)

	if err := s.slackAdapter.SendMessage(
		alarm.User.SlackIncomingWebhook.String,
		alarm.User.SlackChannel.String,
		slackNotificationsUsername,
		slackNotificationsEmoji,
		newIncidentMessage,
	); err != nil {
		logger.Errorf("Send Slack message error: %s", err)
		return
	}
}
