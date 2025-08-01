# Subprocess - Signal-CLI HTTP

This module spawns and handles IO for the signal-cli process.

Do not pass an object with the "id" key into this module's methods. It will reject the request for that reason.

This system works with multiple requests at the same time safely.

This system also caches incoming messages up to at least 5 minutes old for later querying. This process takes in a filter JSON and goes through this list and finds any incoming message JSON objects that match to the filter JSON as outlined in the [auth module](../auth/readme.md).