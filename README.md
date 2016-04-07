[![Codeship Status for RichardKnop/ping](https://codeship.com/projects/fb4fa9f0-c2bb-0133-461d-4e6bd7c806c7/status?branch=master)](https://codeship.com/projects/137882)

# Ping List

API / website uptime & performance monitoring platform.

# Index

* [Ping List](#ping-list)
* [Index](#index)
* [Docs](../../../ping-list/blob/master/docs/)
* [Dependencies](#dependencies)
* [Setup](#setup)
* [Test Data](#test-data)
* [Testing](#testing)
* [Docker](#docker)

# Dependencies

According to [Go 1.5 Vendor experiment](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo), all dependencies are stored in the vendor directory. This approach is called `vendoring` and is the best practice for Go projects to lock versions of dependencies in order to achieve reproducible builds.

To update dependencies during development:

```
make update-deps
```

To install dependencies:

```
make install-deps
```

# Setup

If you are developing on OSX, install `etcd`, `Postgres`:

```
brew install etcd
brew install postgres
```

You might want to create a `Postgres` database:

```
createuser --createdb pinglist
createdb -U pinglist pinglist
```

Load a development configuration into `etcd`:

```
curl -L http://localhost:2379/v2/keys/config/pinglist.json -XPUT -d value='{
	"Database": {
		"Type": "postgres",
		"Host": "localhost",
		"Port": 5432,
		"User": "pinglist",
		"Password": "",
		"DatabaseName": "pinglist",
		"MaxIdleConns": 5,
		"MaxOpenConns": 5
	},
	"Oauth": {
		"AccessTokenLifetime": 3600,
		"RefreshTokenLifetime": 1209600,
		"AuthCodeLifetime": 3600
	},
	"Session": {
		"Secret": "test_secret",
		"Path": "/",
		"MaxAge": 604800,
		"HTTPOnly": true
	},
	"AWS": {
		"Region": "us-west-2",
		"APNSPlatformApplicationARN": "apns_platform_application_arn",
		"GCMPlatformApplicationARN":  "gcm_platform_application_arn"
	},
	"Facebook": {
		"AppID": "facebook_app_id",
		"AppSecret": "facebook_app_secret"
	},
	"Stripe": {
		"SecretKey": "stripe_secret_key",
		"PublishableKey": "stripe_publishable_key"
	},
	"Sendgrid": {
		"APIKey": "sendgrid_api_key"
	},
	"Web": {
		"Scheme": "http",
		"Host": "localhost:8080"
	},
	"IsDevelopment": true
}'
```

Run migrations:

```
go run main.go migrate
```

And finally, run the app:

```
go run main.go runserver
```

When deploying, you can set `ETCD_HOST` and `ETCD_PORT` environment variables.

# Test Data

You might want to insert some test data if you are testing locally using `curl` examples from this README:

```
go run main.go loaddata \
	oauth/fixtures/scopes.yml \
	oauth/fixtures/test_clients.yml \
	oauth/fixtures/test_users.yml \
	accounts/fixtures/roles.yml \
	accounts/fixtures/test_accounts.yml \
	accounts/fixtures/test_users.yml \
	subscriptions/fixtures/plans.yml \
	alarms/fixtures/regions.yml \
	alarms/fixtures/alarm_states.yml \
	alarms/fixtures/incident_types.yml
```

# Testing

I have used a mix of unit and functional tests so you need to have `sqlite` installed in order for the tests to run successfully as the suite creates an in-memory database.

Set the `STRIPE_KEY` environment variable to match your test private key, then run `make test`:

```
STRIPE_KEY=YOUR_API_KEY make test
```

# Docker

Build a Docker image and run the app in a container:

```
docker build -t pinglist-api .
docker run -e ETCD_HOST=localhost -e ETCD_PORT=2379 -p 6060:8080 pinglist-api
```

You can load fixtures with `docker exec` command:

```
docker exec <container_id> /go/bin/pinglist-api loaddata \
	oauth/fixtures/scopes.yml \
	accounts/fixtures/roles.yml \
	subscriptions/fixtures/plans.yml \
	alarms/fixtures/regions.yml \
	alarms/fixtures/alarm_states.yml \
	alarms/fixtures/incident_types.yml \
	oauth/fixtures/test_clients.yml \
	accounts/fixtures/test_accounts.yml
```

You can also execute interactive commands by passing `-i` flag:

```
docker exec -i <container_id> /go/bin/pinglist-api createaccount
docker exec -i <container_id> /go/bin/pinglist-api createsuperuser
```
