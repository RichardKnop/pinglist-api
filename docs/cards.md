# Cards

* [Create Card](#create-card)
* [Get Card](#get-card)
* [Delete Card](#delete-card)
* [List Card](#list-cards)

## Create Card

Example request:

```
curl --compressed -v localhost:8080/v1/cards \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"token": "token"
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/cards/1"
		}
	},
	"id": 1,
	"brand": "Visa",
	"funding": "credit",
	"last_four": "4242",
	"exp_month": 10,
	"exp_year": 2020,
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Get Card

Example request:

```
curl --compressed -v localhost:8080/v1/cards/1 \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/cards/1"
		}
	},
	"id": 1,
	"brand": "Visa",
	"funding": "credit",
	"last_four": "4242",
	"exp_month": 10,
	"exp_year": 2020,
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Delete Card

Example request:

```
curl -XDELETE --compressed -v localhost:8080/v1/cards/1 \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Returns `204` empty response on success.

## List Cards

Example request:

```
curl --compressed -v "localhost:8080/v1/cards" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Use `page` and `limit` query string parameters to paginate and `order_by` to order the results.

Optionally filter results with `user_id` query string parameter.

Notice the ampersand is escaped as `\u0026` in the `_links` section.

Example response:

```json
{
	"_links": {
		"first": {
			"href": "/v1/cards?page=1"
		},
		"last": {
			"href": "/v1/cards?page=1"
		},
		"next": {
			"href": ""
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/cards"
		}
	},
	"_embedded": {
		"cards": [
			{
				"_links": {
					"self": {
						"href": "/v1/cards/1"
					}
				},
        "id": 1,
      	"brand": "Visa",
				"funding": "credit",
      	"last_four": "4242",
				"exp_month": 10,
				"exp_year": 2020,
      	"created_at": "2016-01-14T13:52:24Z",
      	"updated_at": "2016-01-14T13:52:24Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/cards/2"
					}
				},
        "id": 1,
      	"brand": "MasterCard",
				"funding": "debit",
      	"last_four": "4444",
				"exp_month": 10,
				"exp_year": 2020,
      	"created_at": "2016-01-14T13:52:24Z",
      	"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 2,
	"page": 1
}
```
