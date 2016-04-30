# Accounts

* [Contact](#contact)

## Contact

Sends an email to the admin.

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/contact \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret \
	-d '{
		"name": "John Reese",
    "email": "john@reese.com",
    "subject": "Test Subject",
    "message": "Test Message"
	}'
```

Returns `204` empty response on success.
