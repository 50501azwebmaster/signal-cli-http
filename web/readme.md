# Web - Signal-CLI HTTP

This module handles the HTTP requests. The path used for requests genuinely does not matter. However, every request must have a JSON object body, and an `Authorization: bearer <token>` header. The response will also always be a JSON.

Possible response codes:

* 400 Bad Request: Bad JSON, missing Authorization header, etc.
* 401 Unauthorized: bearer token not allowed for presented request.
* 500 Internal Server Error: Self explanitory.
* 501 Not Implemented: Will be returned for message receiving.
* 200 OK: Request was forwarded to JSONRPC subprocess and that process returned something, including if the "error" key is present in the returned JSON.

Any non-200 response code will come with a JSON with the following format:

```json
{"error":"Error message string"}
```