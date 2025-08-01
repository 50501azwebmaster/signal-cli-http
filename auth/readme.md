# Auth - Signal-CLI HTTP

This module handles the reading and parsing of the auth JSON file. It also acts as a verifier in relation to that information. The file is a JSON object. It acts as a whitelist for which bearer token can do what action. It is passed to the HTTP endpoint via the `Authorization: <bearerToken>` header. Nore that this is not `Authorization: Bearer <token>`

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

When an HTTP request comes in, this software will do the following (sending error responses when appropriate):

1. Check that there's an `Authorization` header
2. Get the authorization header's value (bearer token)
3. Read the JSON array corresponding to the bearer token.
4. See if any JSON object in that array "matches" the request JSON
5. Forward the request to the subprocess, and return the result.

The rules for matching a request JSON to a filter JSON is a recursive process. At each step it goes through the following checks:
1. The types (map JSON, array JSON, value literal, etc) must be the same
2. For the map JSON type, each key inside the filter json must be present inside the request JSON. Each key-value pair in the filter JSON must also match (recursively).
3. For the array JSON type, each key inside the filter json must be present inside the request JSON **and vice-versa**. Each key-value pair must also match (recursively).
4. For anything else, it matches the object directly. This is invoked when checking equality of a value literal.

Here's some examples for each case:

1. the request `{"method":"send","params":{"recipient":["+16028675309"],"message":"message"},"id":"SomeID"},` would not match the filter `["+5555555555"]` because one is a JSON map and the other a JSON array.
2. the request `{"method":"something","params":{"recipient":["+16028675309"],"message":"message"},"id":"SomeID"},` would not match the filter `{"method":"send","params":{"recipient":["+16028675309"],"message":"message"}}` because the "method" differs. This would also fail to match if the `method` key was missing in the request JSON.
3. `{"method":"send","params":{"recipient":["+16028675309","someBadNumber"]}}` would not match the filter `{"method":"send","params":{"recipient":["+16028675309",]}}`  because of the `someBadNumber` number in the request. This rule exists so that a malicious request cant send a message to both a room/concact that it's whitelisted for, and one that it isn't.
4. `"+16028675309"` would not match the filter `"+15555555555"` because their values differ.
Here's what each filter JSON object in the above sample JSON does: 

`{"method":"send","params":{"recipient":["+16028675309"]}}` allows sending to `+16028675309` (any message, timestamp, etc.)
`{"method":"send","params":{"groupId":["67a13c3e-8d29-2539-ce8e-41129c349d6d"]}}`: allows sending to group `67a13c3e-8d29-2539-ce8e-41129c349d6d` (any message, timestamp, etc.)
`{"method":"receive","params":{"envelope":{"source":"67a13c3e-8d29-2539-ce8e-41129c349d6d"}}}` allows receiving from group `67a13c3e-8d29-2539-ce8e-41129c349d6d`