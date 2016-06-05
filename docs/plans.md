# Plans

* [List Plans](#list-plans)

## List Plans

Example request:

```
curl --compressed -v "localhost:8080/v1/plans" \
	-u test_client_1:test_secret
```

Example response:

```json
{
    "_links": {
        "first": {
            "href": "/v1/plans"
        },
        "last": {
            "href": "/v1/plans"
        },
        "next": {
            "href": ""
        },
        "prev": {
            "href": ""
        },
        "self": {
            "href": "/v1/plans"
        }
    },
    "_embedded": {
        "plans": [
            {
                "_links": {
                    "self": {
                        "href": "/v1/plans/1"
                    }
                },
                "id": 1,
                "plan_id": "personal",
                "name": "Personal",
                "description": "Personal website and/or a blog.",
                "currency": "USD",
                "amount": 400,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 5,
                "max_teams": 0,
                "max_members_per_team": 0,
								"min_alarm_interval": 60,
                "unlimited_emails": false,
                "max_emails_per_interval": 100,
                "unlimited_push_notifications": true,
                "max_push_notifications_per_interval": null,
                "slack_alerts": false,
                "created_at": "2016-01-14T13:52:24Z",
                "updated_at": "2016-01-14T13:52:24Z"
            },
            {
                "_links": {
                    "self": {
                        "href": "/v1/plans/2"
                    }
                },
                "id": 2,
                "plan_id": "professional",
                "name": "Professional",
                "description": "Monitor up to 10 different websites or APIs.",
                "currency": "USD",
                "amount": 1000,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 15,
                "max_teams": 0,
                "max_members_per_team": 0,
								"min_alarm_interval": 60,
								"unlimited_emails": false,
                "max_emails_per_interval": 300,
                "unlimited_push_notifications": true,
                "max_push_notifications_per_interval": null,
                "slack_alerts": false,
                "created_at": "2016-01-14T13:52:24Z",
                "updated_at": "2016-01-14T13:52:24Z"
            },
            {
                "_links": {
                    "self": {
                        "href": "/v1/plans/3"
                    }
                },
                "id": 3,
                "plan_id": "startup",
                "name": "Startup",
                "description": "Create a team of 10 members and monitor up to 100 APIs.",
                "currency": "USD",
                "amount": 4000,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 75,
                "max_teams": 1,
                "max_members_per_team": 10,
								"min_alarm_interval": 30,
								"unlimited_emails": true,
                "max_emails_per_interval": null,
                "unlimited_push_notifications": true,
                "max_push_notifications_per_interval": null,
                "slack_alerts": true,
                "created_at": "2016-01-14T13:52:24Z",
                "updated_at": "2016-01-14T13:52:24Z"
            },
            {
                "_links": {
                    "self": {
                        "href": "/v1/plans/4"
                    }
                },
                "id": 4,
                "plan_id": "business",
                "name": "Business",
                "description": "Create 5 teams of 10 members each and monitor up to 200 APIs.",
                "currency": "USD",
                "amount": 12000,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 300,
                "max_teams": 10,
                "max_members_per_team": 20,
								"min_alarm_interval": 30,
								"unlimited_emails": true,
                "max_emails_per_interval": null,
                "unlimited_push_notifications": true,
                "max_push_notifications_per_interval": null,
                "slack_alerts": true,
                "created_at": "2016-01-14T13:52:24Z",
                "updated_at": "2016-01-14T13:52:24Z"
            }
        ]
    },
    "count": 5,
    "page": 1
}
```
