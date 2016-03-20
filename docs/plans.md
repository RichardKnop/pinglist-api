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
				"description": "description ...",
				"currency": "USD",
				"amount": 250,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 1,
				"max_team_members": 1,
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
				"description": "description ...",
				"currency": "USD",
				"amount": 1000,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 5,
				"max_team_members": 1,
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
				"description": "description ...",
				"currency": "USD",
				"amount": 35000,
				"trial_period": 30,
				"interval": 30,
				"max_alarms": 100,
				"max_team_members": 20,
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 3,
	"page": 1
}
```
