# Alarm Metrics

* [List Alarm Response Times](#list-alarm-response-times)

## List Alarm Response Times

Example request:

```
curl --compressed -v "localhost:8080/v1/alarms/1/response-times?date_trunc=day&from=2016-02-08T00:00:00Z&to=2016-02-09T00:00:00Z" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Use `page` and `limit` query string parameters to paginate and `order_by` to order the results.

Use `date_trunc` to query for average hourly/daily results:

- `hour`: aggregated hourly results
- `day`: aggregated daily results
- etc

Use `from` and `to` parameters to query a specific time range. Pass timestamps formatted according to `RFC3339`.

Notice the ampersand is escaped as `\u0026` in the `_links` section.

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/alarms/1/response-times?page=1"
		},
		"last": {
			"href": "/v1/alarms/1/response-times?page=2"
		},
		"next": {
			"href": "/v1/alarms/1/response-times?page=2"
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/alarms/1/response-times"
		}
	},
	"_embedded": {
		"response_times": [
			{
				"timestamp": "2016-01-14T13:52:24Z",
				"value": 12345
			},
			{
				"timestamp": "2016-01-14T13:53:24Z",
				"value": 12345
			}
		]
	},
	"uptime": 99.99,
	"average": 12345.0,
	"incident_type_counts": {
		"slow_response": 0,
		"timeout": 2,
		"bad_code": 1,
		"other": 1,
	},
	"count": 2,
	"page": 1
}
```
