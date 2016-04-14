# Alarm Incidents

* [List Alarm Incidents](#list-alarm-incidents)

## List Alarm Incidents

Example request:

```
curl --compressed -v "localhost:8080/v1/alarms/1/incidents" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Use `offset` and `limit` query string parameters to paginate and `order_by` to order the results.

Notice the ampersand is escaped as `\u0026` in the `_links` section.

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/alarms/1/incidents?page=1"
		},
		"last": {
			"href": "/v1/alarms/1/incidents?page=2"
		},
		"next": {
			"href": "/v1/alarms/1/incidents?page=2"
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/alarms/1/incidents"
		}
	},
	"_embedded": {
		"incidents": [
			{
				"_links": {
					"self": {
						"href": "/v1/alarms/1/incidents/1"
					}
				},
				"id": 1,
				"alarm_id": 1,
				"type": "timeout",
				"http_code": null,
				"response_time": null,
				"response": null,
				"error_message": "timeout error...",
				"resolved_at": null,
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			},
      			{
				"_links": {
					"self": {
						"href": "/v1/alarms/1/incidents/2"
					}
				},
				"id": 2,
				"alarm_id": 1,
				"type": "bad_code",
				"http_code": 500,
				"response_time": 1426,
				"response": "Internal Server Error",
				"error_message": null,
				"resolved_at": null,
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			},
		]
	},
	"count": 2,
	"page": 1
}
```
