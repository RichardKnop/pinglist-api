# Regions

* [List Regions](#list-regions)

## List Regions

Example request:

```
curl --compressed -v "localhost:8080/v1/regions" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/regions"
		},
		"last": {
			"href": "/v1/regions"
		},
		"next": {
			"href": ""
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/regions"
		}
	},
	"_embedded": {
		"regions": [
			{
				"_links": {
					"self": {
						"href": "/v1/regions/1"
					}
				},
				"id": "SGP",
				"name": "Singapore"
			}
		]
	},
	"count": 4,
	"page": 1
}
```
