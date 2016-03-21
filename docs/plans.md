# Plans

* [List Plans](#list-plans)

## List Plans

Example request:

```
curl --compressed -v "localhost:8080/v1/plans" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
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
				"amount": 500,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 2,
				"max_team_members": 0,
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
				"amount": 2000,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 15,
				"max_team_members": 0,
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
				"description": "Create a team of 5 members with 10 alarms each.",
				"currency": "USD",
				"amount": 2000,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 10,
				"max_team_members": 5,
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
				"description": "Create a team of 10 members with 10 alarms each.",
				"currency": "USD",
				"amount": 15000,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 10,
				"max_team_members": 10,
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
				"plan_id": "enterprise",
				"name": "Enterprise",
				"description": "Create a team of 30 members with 10 alarms each.",
				"currency": "USD",
				"amount": 35000,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 10,
				"max_team_members": 30,
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 3,
	"page": 1
}
```
