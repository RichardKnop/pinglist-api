# Devices

* [Register Device](#register-device)

## Register Device

Platforms:
- `iOS`
- `Android`

Example request:

```
curl --compressed -v localhost:8080/v1/devices \
	-H "Content-Type: application/json" \
  -H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
    "platform": "iOS",
    "token": "AAAA1111BBBB2222AAAA1111BBBB2222AAAA1111BBBB2222AAAA1111BBBB2222",
  }'
```

Returns `204` empty response on success.
