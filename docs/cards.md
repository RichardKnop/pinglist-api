# Subscriptions

* [Checkout Button](#checkout-button)
* [Create Card](#create-card)
* [Delete Card](#delete-card)
* [List Card](#list-cards)

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
	"card_id": "card_17ssTxKkL3BsdwCiJJPZQc8m",
	"brand": "Visa",
	"last_four": "4242",
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

Use `offset` and `limit` query string parameters to paginate and `order_by` to order the results.

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
      	"card_id": "card_17ssTxKkL3BsdwCiJJPZQc8m",
      	"brand": "Visa",
      	"last_four": "4242",
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
      	"card_id": "card_Jd83fsafH94dIFSF8fasf02b",
      	"brand": "Visa",
      	"last_four": "4343",
      	"created_at": "2016-01-14T13:52:24Z",
      	"updated_at": "2016-01-14T13:52:24Z"
			}
		]
	},
	"count": 2,
	"page": 1
}
```
