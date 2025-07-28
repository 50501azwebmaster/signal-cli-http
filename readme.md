# Signal-CLI HTTP

**Very** early in development.

Very simple HTTP frontend to [signal-cli](https://github.com/AsamK/signal-cli) JSON RPC.

Please see the JSONRPC documentation for `signal-cli`: [https://github.com/AsamK/signal-cli/blob/master/man/signal-cli-jsonrpc.5.adoc](https://github.com/AsamK/signal-cli/blob/master/man/signal-cli-jsonrpc.5.adoc)

Please also read the following README files for the individual modules to understand how to configure and interact with this program:

* [args](args/readme.md) handles command line arguments.
* [auth](auth/readme.md) handles the authentication JSON and checking requests.
* [subprocess](subprocess/readme.md) manages the underlying `signal-cli` JSONRPC process, along with caching incoming messages.
* [web](web/readme.md) - handles the HTTP requests to this program, including the necessary edge cases.