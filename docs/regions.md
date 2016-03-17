# Regions

* [List Regions](#list-regions)

## List Regions

Example request:

```
curl --compressed -v "localhost:8080/v1/alarms/regions" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/alarms/regions"
		}
	},
	"_embedded": {
		"regions": [
			{
				"_links": {
					"self": {
						"href": "/v1/alarms/regions/1"
					}
				},
				"id": "SGP",
				"name": "Singapore"
			}
		]
	}
}
```
