# Signal-CLI HTTP

Very simple HTTP frontend to [signal-cli](https://github.com/AsamK/signal-cli) JSON RPC.

Please see the JSONRPC documentation for `signal-cli`: [https://github.com/AsamK/signal-cli/blob/master/man/signal-cli-jsonrpc.5.adoc](https://github.com/AsamK/signal-cli/blob/master/man/signal-cli-jsonrpc.5.adoc)

Please also read the following README files for the individual modules to understand how to configure and interact with this program:

* [args](args/readme.md) handles command line arguments. Go here to learn how to run the program.
* [auth](auth/readme.md) handles the authentication JSON and checking requests. Go here to learn how to secure the program, and whitelist authentication keys.
* [subprocess](subprocess/readme.md) manages the underlying `signal-cli` JSONRPC process, along with caching incoming messages. Go here to understand how this program relays and returns the JSON objects.
* [web](web/readme.md) - handles the HTTP requests to this program, including the necessary edge cases. Go here to understand how to send and understand the responses to the program's HTTP endpoint.