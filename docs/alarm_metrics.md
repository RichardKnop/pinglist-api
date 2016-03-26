# Alarm Metrics

* [List Alarm Request Times](#list-alarm-request-times)

## List Alarm Request Times

Example request:

```
curl --compressed -v "localhost:8080/v1/alarms/1/request-times" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Use `offset` and `limit` query string parameters to paginate and `order_by` to order the results.

Notice the ampersand is escaped as `\u0026` in the `_links` section.

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/alarms/1/request-times?page=1"
		},
		"last": {
			"href": "/v1/alarms/1/request-times?page=2"
		},
		"next": {
			"href": "/v1/alarms/1/request-times?page=2"
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/alarms/1/request-times"
		}
	},
	"_embedded": {
		"request_times": [
			{
				"timestamp": "2016-01-14T13:52:24Z",
				"request_time": 12345
			},
      {
				"timestamp": "2016-01-14T13:53:24Z",
				"request_time": 12345
			},
		]
	},
	"count": 2,
	"page": 1
}
```
