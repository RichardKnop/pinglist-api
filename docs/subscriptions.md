# Subscriptions

* [Checkout Button](#checkout-button)
* [Subscribe User](#subscribe-user)

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

## Subscribe User

Example request:

```
curl --compressed -v localhost:8080/v1/subscriptions \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"stripe_token": "token",
		"stripe_email": "test@user.com",
		"plan_id": 1
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
		"customer": {
			"_links": {
				"self": {
					"href": "/v1/customers/1"
				}
			},
			"id": 1,
			"user_id": 1,
			"customer_id": "cus_7z94mLsfxLva84",
			"created_at": "2016-01-14T13:52:24Z",
			"updated_at": "2016-01-14T13:52:24Z"
		},
		"plan": {
			"_links": {
				"self": {
					"href": "/v1/plans/1"
				}
			},
			"id": 1,
			"plan_id": "personal",
			"currency": "GBP",
			"amount": 250,
			"trial_period": 30,
			"interval": 30,
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
