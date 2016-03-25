# Subscriptions

* [Checkout Button](#checkout-button)
* [Create Subscription](#create-subscription)
* [Update Subscription](#update-subscription)
* [Cancel Subscription](#cancel-subscription)
* [List Subscriptions](#list-subscriptions)

## Checkout Button

Checkout supports two different integrations:

- *Simple*: The (simple integration)[https://stripe.com/docs/checkout#integration-simple] provides a blue "Pay with card" button and submits your payment form with a Stripe token in a hidden input field.
- *Custom*: The [custom integration](https://stripe.com/docs/checkout#integration-custom) lets you create a custom button and passes a Stripe token to a JavaScript callback.

Simple integration example:

```html
<form action="" method="POST">
	<script
		src="https://checkout.stripe.com/checkout.js" class="stripe-button"
		data-key="stripe_publishable_key"
		data-amount="250"
		data-currency="GBP"
		data-name="The name of your company or website"
		data-description="A description of the product or service being purchased"
		data-locale="auto"
		data-email="If you already know the email address of your user, you can provide it to Checkout to be pre-filled">
	</script>
	<input type="hidden" name="planID" value="1">
</form>
```

When submitted, the above form with add `stripeToken` and `stripeEmail` parameters to the request data.

## Create Subscription

Example request:

```
curl --compressed -v localhost:8080/v1/subscriptions \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"plan_id": 1,
		"card_id": 1
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/subscriptions/1"
		}
	},
	"_embedded": {
		"plan":	{
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
		"card": {
			"_links": {
				"self": {
					"href": "/v1/cards/1"
				}
			},
			"id": 1,
			"brand": "Visa",
			"last_four": "4242",
			"created_at": "2016-01-14T13:52:24Z",
			"updated_at": "2016-01-14T13:52:24Z"
		}
	},
	"id": 1,
	"subscription_id": "sub_7z94rezxDE9frw",
	"started_at": "2016-01-14T13:52:24Z",
	"cancelled_at": "",
	"ended_at": "",
	"period_start": "2016-01-14T13:52:24Z",
	"period_end": "2016-02-14T13:52:24Z",
	"trial_start": "",
	"trial_end": "",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Update Subscription

Example request:

```
curl -XPUT --compressed -v localhost:8080/v1/subscriptions \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"plan_id": 2,
		"card_id": 2
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/subscriptions/1"
		}
	},
	"_embedded": {
		"plan": {
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
		"card": {
			"_links": {
				"self": {
					"href": "/v1/cards/2"
				}
			},
			"id": 1,
			"brand": "Visa",
			"last_four": "4343",
			"created_at": "2016-01-14T13:52:24Z",
			"updated_at": "2016-01-14T13:52:24Z"
		}
	},
	"id": 1,
	"subscription_id": "sub_7z94rezxDE9frw",
	"started_at": "2016-01-14T13:52:24Z",
	"cancelled_at": "",
	"ended_at": "",
	"period_start": "2016-01-14T13:52:24Z",
	"period_end": "2016-02-14T13:52:24Z",
	"trial_start": "",
	"trial_end": "",
	"created_at": "2016-01-14T13:52:24Z",
	"updated_at": "2016-01-14T13:52:24Z"
}
```

## Cancel Subscription

Example request:

```
curl -XDELETE --compressed -v localhost:8080/v1/subscriptions/1 \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Returns `204` empty response on success.

## List Subscriptions

Example request:

```
curl --compressed -v "localhost:8080/v1/subscriptions" \
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
			"href": "/v1/subscriptions?page=1"
		},
		"last": {
			"href": "/v1/subscriptions?page=1"
		},
		"next": {
			"href": ""
		},
		"prev": {
			"href": ""
		},
		"self": {
			"href": "/v1/subscriptions"
		}
	},
	"_embedded": {
		"subscriptions": [
			{
				"_links": {
					"self": {
						"href": "/v1/subscriptions/1"
					}
				},
				"_embedded": {
					"plan":	{
						"_links": {
							"self": {
								"href": "/v1/plans/1"
							}
						},
						"id": 1,
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
					"card": {
						"_links": {
							"self": {
								"href": "/v1/cards/1"
							}
						},
						"id": 1,
						"brand": "Visa",
						"last_four": "4242",
						"created_at": "2016-01-14T13:52:24Z",
						"updated_at": "2016-01-14T13:52:24Z"
					}
				},
				"id": 1,
				"subscription_id": "sub_7z94rezxDE9frw",
				"started_at": "2016-02-14T13:52:24Z",
				"cancelled_at": "",
				"ended_at": "",
				"period_start": "2016-02-14T13:52:24Z",
				"period_end": "2016-03-14T13:52:24Z",
				"trial_start": "",
				"trial_end": "",
				"created_at": "2016-02-14T13:52:24Z",
				"updated_at": "2016-02-14T13:52:24Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/subscriptions/2"
					}
				},
				"_embedded": {
					"plan":	{
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
					"card": {
						"_links": {
							"self": {
								"href": "/v1/cards/1"
							}
						},
						"id": 1,
						"brand": "Visa",
						"last_four": "4242",
						"created_at": "2016-01-14T13:52:24Z",
						"updated_at": "2016-01-14T13:52:24Z"
					}
				},
				"id": 2,
				"subscription_id": "sub_87HIdrte99poeq",
				"started_at": "2016-01-14T13:52:24Z",
				"cancelled_at": "2016-01-18T13:52:24Z",
				"ended_at": "2016-02-14T13:52:24Z",
				"period_start": "2016-01-14T13:52:24Z",
				"period_end": "2016-02-14T13:52:24Z",
				"trial_start": "",
				"trial_end": "",
				"created_at": "2016-01-14T13:52:24Z",
				"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 2,
	"page": 1
}
```
