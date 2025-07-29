# Auth - Signal-CLI HTTP

This module handles the reading and parsing of the auth JSON file. It also acts as a verifier in relation to that information. The file is a JSON object. It acts as a whitelist for which bearer token can do what action. It is passed to the HTTP endpoint via the `Authorization: Bearer <bearerToken>` header.

Here's a sample auth JSON:

```json
{
	"WGV99fSwgKhdQSa89HQIGxas": [
		{"method":"send","params":{"recipient":["+16028675309"]}},
		{"method":"send","params":{"groupId":["67a13c3e-8d29-2539-ce8e-41129c349d6d"]}},
	],
	"ZQR3T6lqsvnXcgcWhpPOWWdv": [
		{"method":"receive","params":{"envelope":{"source":"67a13c3e-8d29-2539-ce8e-41129c349d6d"}}}
	]
}
```

When an HTTP request comes in, this software will do the following:

1. Check that there's an `Authorization` header
2. Get the authorization header's value (bearer token)
3. Read the JSON array corresponding to the bearer token.
4. See if any JSON object in that array (called a filter) does not have any data the request JSON doesn't, except for arrays which must match excactly.
5. If the statement in step 4 is true, forward the request into the signal-cli process and return the response.

So for example, the reqest `{"method":"send","params":{"recipient":["+16028675309"],"message":"message"},"id":"SomeID"},` would be allowed by the filter `{"method":"send","params":{"recipient":["+16028675309"]}}` because the filter does not have any data the request does not. But `{"method":"send","params":{"recipient":["+5555555555"],"message":"message"},"id":"SomeID"},` would not because the phone number differs.

Note: items in arrays must "match" exactly, but items in items in arrays follow normal rules. So the request `{"method":"send","params":{"recipient":["+16028675309","someBadNumber"]}}` would NOT match the filter `{"method":"send","params":{"recipient":["+16028675309",]}}` 

These filters can be as granular as you want.

Here's what each filter JSON object in the above sample JSON does: 

`{"method":"send","params":{"recipient":["+16028675309"]}}` allows sending to `+16028675309` (any message, timestamp, etc.)
`{"method":"send","params":{"groupId":["67a13c3e-8d29-2539-ce8e-41129c349d6d"]}}`: allows sending to group `67a13c3e-8d29-2539-ce8e-41129c349d6d` (any message, timestamp, etc.)
`{"method":"receive","params":{"envelope":{"source":"67a13c3e-8d29-2539-ce8e-41129c349d6d"}}}` allows receiving from group `67a13c3e-8d29-2539-ce8e-41129c349d6d`