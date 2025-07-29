package args

/* This file manages declaration of, parsing of, and returning copies of
   command line arguments. */

import (
	"flag"
)

/* Method to trigger flag parsing. Can be called multiple times safely. */
func Parse() {
	if (flagsParsed) {return}
	flag.Parse();
	flagsParsed = true;
}
/* module-specific variable to avoid re-parsing flags */
var flagsParsed bool = false;

/* what file to read the configuration */
var authJson = flag.String("auth", "./auth.json", "Authorization file to read from")
/* @return set boolean will be true if argument is not nil */
func GetAuthJson() (location string, set bool) {
	if authJson == nil {return "", false}
	return *authJson, true;
}

/* TCP port to bind to */
var httpPort = flag.Int("port", 11938, "Port number to bind to")
/* @return set boolean will be true if argument is not nil */
func GetHTTPPort() (port int, set bool) {
	if httpPort == nil {return -1, false}
	return *httpPort, true;
}

/* Where the signal-cli binary is */
var binaryLocation = flag.String("binary", "/usr/local/bin/signal-cli", "Location of the signal-cli binary.")
/* @return set boolean will be true if argument is not nil */
func GetBinaryLocation() (binary string, set bool) {
	if binaryLocation == nil {return "", false}
	return *binaryLocation, true;
}