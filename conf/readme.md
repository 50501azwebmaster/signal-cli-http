# Conf - Signal-CLI HTTP

This module handles reading and parsing the config file, and acting as a verifier for the `Authorization` header on the HTTP requests.

The config file is made up of multiple lines. The first token in each line is the `Authorization` bearer token. This cannot have spaces but can be any string. Choose wisely. The remainder of the line contains a path that the `Authorization` header is checked against. It does not matter if you include a leading or trailing slash.
   
Here's a sample config:

```
WGV99fSwgKhdQSa89HQIGxas /+16028675309/room/roomID/*
WGV99fSwgKhdQSa89HQIGxas /+16028675309/direct/username.69/send
ZQR3T6lqsvnXcgcWhpPOWWdv +16028675309/direct/username.69/send/
```

The config file is a **whitelist** for each bearer token to access a specific endpoint (or set of endpoints). The endpoints for this program are granular enough to only allow one action for each endpoint, so this level of whitelisting should™ be okay.

There is a regex-like behavior to these paths using the `*` and `?` characters. For the regex-like behavior to be triggered these characters must be by themselves per path segment (no other characters not separated by a `/` or a start or end of string).

The `*` character matches to any number of path segments. The `?` character matches to only one segment. Here's some examples:

* `HZJWwB0TAjz6pjAHosII5ofR /+16028675309/*` will allow the bearer token to access any endpoint with the phone number `+16028675309`
* `HZJWwB0TAjz6pjAHosII5ofR /+16028675309/direct/?/send` will allow the bearer token to send a direct message to anyone on that phone number.