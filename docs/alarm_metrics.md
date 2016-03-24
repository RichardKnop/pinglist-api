# Alarm Metrics

* [List Alarm Request Times](#list-alarm-request-times)

## List Alarm Request Times

Example request:

```
curl --compressed -v "localhost:8080/v1/alarms/1/requesttimes?page=1" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Use `offset` and `limit` query string parameters to paginate and `order_by` to order the results.

Notice the ampersand is escaped as `\u0026` in the `_links` section.

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/alarms/1/requesttimes?page=1"
		},
		"last": {
			"href": "/v1/alarms/1/requesttimes?page=2"
		},
		"next": {
			"href": "/v1/alarms/1/requesttimes?page=2"
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/alarms/1/requesttimes?page=1"
		}
	},
	"_embedded": {
		"requesttimes": [
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
	"count": 4,
	"page": 1
}
```
