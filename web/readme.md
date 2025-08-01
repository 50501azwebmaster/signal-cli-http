# Web - Signal-CLI HTTP

This module handles the HTTP requests. The path used for requests genuinely does not matter. However, every request must have a JSON object body, and an `Authorization: bearer <token>` header.

Possible response codes:

* 400 Bad Request: Bad JSON, missing Authorization header, **sent a request with a defined id**, etc.
* 401 Unauthorized: bearer token not allowed for presented request.
* 500 Internal Server Error: Self explanitory.
* 501 Not Implemented: Will be returned for message receiving.
* 200 OK: Request was forwarded to JSONRPC subprocess and that process returned something, including if the "error" key is present in the returned JSON.

There is no body content with non-200 response codes. With 200 the response is a valid JSON map or array.

This program simply relays requests to the signal-cli program. **It will not prevent you from breaking anything, outside of not whitelisting certain requests. This program does not understand what requests mean.** Each request comes formatted as a JSON object outlined in the [JSON-RPC documentation](https://github.com/AsamK/signal-cli/blob/master/man/signal-cli-jsonrpc.5.adoc).

The program will ensure that the request object is a JSON map, and that the `request` key is present. For any request type that is not `receive`, the program will generate an ID for your request (do not put on in the request, it will return an error) and return the program's response.

For `receive` it will return a JSON list of JSON maps in the incoming message cache that match to your request. There is no limit to how many messages it returns so be careful.