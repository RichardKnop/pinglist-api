# Accounts

* [Create User](#create-user)
* [Confirm Email](#confirm-email)
* [Get Me](#get-me)
* [Get User](#get-user)
* [Update User](#update-user)
* [Change Password](#change-password)
* [Create Invitation](#create-invitation)
* [Confirm Invitation](#confirm-invitation)
* [Create Password Reset](#create-password-reset)
* [Confirm Password Reset](#confirm-password-reset)

## Create User

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/users \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret \
	-d '{
		"email": "test@user",
		"password": "test_password"
	}'
```

Example response:

```json
{
    "_links": {
        "self": {
            "href": "/v1/accounts/users/1"
        }
    },
    "id": 1,
    "email": "test@user",
    "first_name": "",
    "last_name": "",
    "role": "user",
    "confirmed": false,
    "slack_api_key": "",
    "slack_channel": "",
    "created_at": "2015-12-17T06:17:54Z",
    "updated_at": "2015-12-17T06:17:54Z"
}
```

## Confirm Email

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/confirmations/bd89167b-0e82-46a0-8e30-cff25436f40f \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret
```

Returns `204` empty response on success.

## Get Me

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/me \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
    "_links": {
        "self": {
            "href": "/v1/accounts/users/1"
        }
    },
    "id": 1,
    "email": "test@user",
    "first_name": "",
    "last_name": "",
    "role": "user",
    "confirmed": true,
    "slack_api_key": "",
    "slack_channel": "",
    "created_at": "2015-12-17T06:17:54Z",
    "updated_at": "2015-12-17T06:17:54Z"
}
```

## Get User

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/users/1 \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
    "_links": {
        "self": {
            "href": "/v1/accounts/users/1"
        }
    },
    "id": 1,
    "email": "test@user",
    "first_name": "",
    "last_name": "",
    "role": "user",
    "confirmed": true,
    "slack_api_key": "",
    "slack_channel": "",
    "created_at": "2015-12-17T06:17:54Z",
    "updated_at": "2015-12-17T06:17:54Z"
}
```

## Update User

Example request:

```
curl -XPUT --compressed -v localhost:8080/v1/accounts/users/1 \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"first_name": "test_first_name_updated",
		"last_name": "test_last_name_updated",
		"slack_api_key": "slack_api_key",
		"slack_channel": "slack_channel"
	}'
```

Example response:

```json
{
    "_links": {
        "self": {
            "href": "/v1/accounts/users/1"
        }
    },
    "id": 1,
    "email": "test@user",
    "first_name": "test_first_name_updated",
    "last_name": "test_last_name_updated",
    "role": "user",
    "confirmed": true,
    "slack_api_key": "slack_api_key",
    "slack_channel": "slack_channel",
    "created_at": "2015-12-17T06:17:54Z",
    "updated_at": "2015-12-18T07:09:15Z"
}
```

## Change Password

Example request:

```
curl -XPUT --compressed -v localhost:8080/v1/accounts/users/1 \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"password": "test_password",
		"new_password": "some_new_password",
	}'
```

Example response:

```json
{
    "_links": {
        "self": {
            "href": "/v1/accounts/users/1"
        }
    },
    "id": 1,
    "email": "test@user",
    "first_name": "test_first_name",
    "last_name": "test_last_name",
    "role": "user",
    "confirmed": true,
    "slack_api_key": "slack_api_key",
    "slack_channel": "slack_channel",
    "created_at": "2015-12-17T06:17:54Z",
    "updated_at": "2015-12-18T07:09:15Z"
}
```

## Create Invitation

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/invitations \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"email": "john@reese",
		"role": "user"
	}'
```

Example response:

```json
{
    "_links": {
        "self": {
            "href": "/v1/accounts/users/invitations/1"
        }
    },
    "id": 1,
    "reference": "57040678-e910-4de0-a3e6-e30c3851289b",
    "invited_user_id": 2,
    "invited_by_user_id": 1,
    "created_at": "2015-12-11T04:42:19Z",
    "updated_at": "2015-12-11T04:42:19Z"
}
```

The invited user should receive an email with a link to a web page where he/she can set a password and therefor activate the account:

![Set Password Screenshot][1]

## Confirm Invitation

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/invitations/bd89167b-0e82-46a0-8e30-cff25436f40f \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"password": "new_password"
	}'
```

Returns `204` empty response on success.

## Create Password Reset

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/password-reset \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret \
	-d '{
		"email": "test@user"
	}'
```

Returns `204` empty response on success.

## Confirm Password Reset

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/password-reset/bd89167b-0e82-46a0-8e30-cff25436f40f \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"password": "new_password"
	}'
```

Returns `204` empty response on success.
