# Teams

* [Create Team](#create-team)
* [Get Team](#get-team)
* [Update Team](#update-team)
* [List Teams](#list-teams)

## Create Team

Example request:

```
curl --compressed -v localhost:8080/v1/teams \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"name": "Test Team 1",
		"members": [
			{
				"id": 3
			}
		]
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/teams/1"
		}
	},
	"_embedded": {
		"members": [
			{
				"_links": {
					"self": {
						"href": "/v1/accounts/users/3"
					}
				},
				"id": 3,
				"email": "john@reese",
				"first_name": "John",
				"last_name": "Reese",
				"role": "user",
				"confirmed": true,
				"created_at": "2015-12-17T06:17:54Z",
				"updated_at": "2015-12-17T06:17:54Z"
			}
		]
	},
	"id": 1,
	"name": "Test Team 1",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Get Team

Example request:

```
curl --compressed -v localhost:8080/v1/teams/1 \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/teams/1"
		}
	},
	"_embedded": {
		"members": [
			{
				"_links": {
					"self": {
						"href": "/v1/accounts/users/3"
					}
				},
				"id": 3,
				"email": "john@reese",
				"first_name": "John",
				"last_name": "Reese",
				"role": "user",
				"confirmed": true,
				"created_at": "2015-12-17T06:17:54Z",
				"updated_at": "2015-12-17T06:17:54Z"
			}
		]
	},
	"id": 1,
	"name": "Test Team 1",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Update Team

Example request:

```
curl -XPUT --compressed -v localhost:8080/v1/teams/1 \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"name": "Test Team 1 Updated",
		"members": [
			{
				"id": 3
			},
			{
				"id": 4
			}
		]
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/teams/1"
		}
	},
	"_embedded": {
		"members": [
			{
				"_links": {
					"self": {
						"href": "/v1/accounts/users/3"
					}
				},
				"id": 3,
				"email": "john@reese",
				"first_name": "John",
				"last_name": "Reese",
				"role": "user",
				"confirmed": true,
				"created_at": "2015-12-17T06:17:54Z",
				"updated_at": "2015-12-17T06:17:54Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/accounts/users/4"
					}
				},
				"id": 4,
				"email": "harold@finch",
				"first_name": "Harold",
				"last_name": "Finch",
				"role": "user",
				"confirmed": true,
				"created_at": "2015-12-17T06:17:54Z",
				"updated_at": "2015-12-17T06:17:54Z"
			}
		]
	},
	"id": 1,
	"name": "Test Team 1 Updated",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## List Teams

Example request:

```
curl --compressed -v localhost:8080/v1/teams \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Use `offset` and `limit` query string parameters to paginate and `order_by` to order the results.

Optionally filter results with `user_id` query string parameter.

Notice the ampersand is escaped as `\u0026` in the `_links` section.

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/teams?page=1"
		},
		"last": {
			"href": "/v1/teams?page=1"
		},
		"next": {
			"href": ""
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/teams"
		}
	},
	"_embedded": {
		"teams": [
			{
				"_links": {
					"self": {
						"href": "/v1/teams/1"
					}
				},
				"_embedded": {
					"members": [
						{
							"_links": {
								"self": {
									"href": "/v1/accounts/users/3"
								}
							},
							"id": 3,
							"email": "john@reese",
							"first_name": "John",
							"last_name": "Reese",
							"role": "user",
							"confirmed": false,
							"created_at": "2015-12-17T06:17:54Z",
							"updated_at": "2015-12-17T06:17:54Z"
						}
					]
				},
				"id": 1,
				"name": "Test Team 1",
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/teams/2"
					}
				},
				"_embedded": {
					"members": [
						{
							"_links": {
								"self": {
									"href": "/v1/accounts/users/4"
								}
							},
							"id": 4,
							"email": "harold@finch",
							"first_name": "Harold",
							"last_name": "Finch",
							"role": "user",
							"confirmed": false,
							"created_at": "2015-12-17T06:17:54Z",
							"updated_at": "2015-12-17T06:17:54Z"
						}
					]
				},
				"id": 2,
				"name": "Test Team 2",
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 2,
	"page": 1
}
```
