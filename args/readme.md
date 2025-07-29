# Args - Signal-CLI HTTP

This module handles command line arguments. A list follows. You can also pass through -h for this list.

```
-auth string
	  Authorization file to read from (default "./auth.json")
-binary string
	  Location of the signal-cli binary. (default "/usr/local/bin/signal-cli")
-port int
	  Port number to bind to (default 11938)
```

Note: the `dummy,py` python file echoes back a valid JSON with just the "id" key in it, for every JSON sent to it with the "id" key. It's useful for testing things that don't need signal-cli specifically.