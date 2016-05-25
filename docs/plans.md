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
                "amount": 250,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 2,
                "max_teams": 0,
                "max_members_per_team": 0,
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
                "max_alarms": 10,
                "max_teams": 0,
                "max_members_per_team": 0,
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
                "amount": 7500,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 100,
                "max_teams": 1,
                "max_members_per_team": 10,
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
                "amount": 15000,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 200,
                "max_teams": 5,
                "max_members_per_team": 10,
                "created_at": "2016-01-14T13:52:24Z",
                "updated_at": "2016-01-14T13:52:24Z"
            },
            {
                "_links": {
                    "self": {
                        "href": "/v1/plans/5"
                    }
                },
                "id": 5,
                "plan_id": "enterprise",
                "name": "Enterprise",
                "description": "Create 10 teams of 10 members each and monitor up to 500 APIs.",
                "currency": "USD",
                "amount": 35000,
                "trial_period": 30,
                "interval": 30,
                "max_alarms": 500,
                "max_teams": 10,
                "max_members_per_team": 10,
                "created_at": "2016-01-14T13:52:24Z",
                "updated_at": "2016-01-14T13:52:24Z"
            }
        ]
    },
    "count": 5,
    "page": 1
}
```
