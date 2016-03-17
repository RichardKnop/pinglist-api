# Alarms

* [Create Alarm](#create-alarm)
* [Update Alarm](#update-alarm)
* [Delete Alarm](#delete-alarm)
* [List Alarms](#list-alarms)

## Create Alarm

Example request:

```
curl --compressed -v localhost:8080/v1/alarms \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"region": "SGP",
		"endpoint_url": "http://endpoint-1",
		"expected_http_code": 200,
		"interval": 60,
		"active": false
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/alarms/1"
		}
	},
	"id": 1,
	"user_id": 1,
	"region": "SGP",
	"endpoint_url": "http://endpoint-1",
	"expected_http_code": 200,
	"interval": 60,
	"active": false,
	"state": "insufficient_data",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Update Alarm

Example request:

```
curl -XPUT --compressed -v localhost:8080/v1/alarms/1 \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"region": "SGP",
		"endpoint_url": "http://endpoint-1-updated",
		"expected_http_code": 201,
		"interval": 90,
		"interval": true
	}'
```

Example response:

```json
{
  "_links": {
		"self": {
			"href": "/v1/alarms/1"
		}
	},
	"id": 1,
	"user_id": 1,
	"region": "SGP",
	"endpoint_url": "http://endpoint-1-updated",
	"expected_http_code": 201,
	"interval": 90,
	"active": true,
	"state": "insufficient_data",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Delete Alarm

Example request:

```
curl -XDELETE --compressed -v localhost:8080/v1/alarms/1 \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Returns `204` empty response on success.

## List Alarms

Example request:

```
curl --compressed -v "localhost:8080/v1/alarms?page=1" \
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
			"href": "/v1/alarms?page=1"
		},
		"last": {
			"href": "/v1/alarms?page=2"
		},
		"next": {
			"href": "/v1/alarms?page=2"
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/alarms?page=1"
		}
	},
	"_embedded": {
		"alarms": [
			{
				"_links": {
					"self": {
						"href": "/v1/alarms/1"
					}
				},
				"id": 1,
				"user_id": 1,
				"region": "SGP",
				"endpoint_url": "http://endpoint-1",
				"expected_http_code": 200,
				"interval": 60,
				"active": true,
				"state": "ok",
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/alarms/2"
					}
				},
				"id": 2,
				"user_id": 1,
				"region": "SGP",
				"endpoint_url": "http://endpoint-2",
				"expected_http_code": 200,
				"interval": 60,
				"active": true,
				"state": "alarm",
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 4,
	"page": 1
}
```
