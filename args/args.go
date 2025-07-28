package args

/* This file manages declaration of, parsing of, and returning copies of
   command line arguments. */

import (
	"flag"
	"os"
)

/* Method to trigger flag parsing. Can be called multiple times safely. */
func Parse() {
	if (flagsParsed) {return}
	flag.Parse();
	flagsParsed = true;
	
	// Process DBUS socket location
	if socketLocation == nil || *socketLocation == ""  {
		*socketLocation = os.Getenv("DBUS_SYSTEM_BUS_ADDRESS");
	}
}
/* module-specific variable to avoid re-parsing flags */
var flagsParsed bool = false;

/* what JSON file to read config values from */
var confLocation = flag.String("conf", "./config.txt", "Config file to read from")
/* @return set boolean will be true if argument is not nil */
func GetConfLocation() (location string, set bool) {
	if confLocation == nil {return "", false}
	return *confLocation, true;
}

/* TCP port to bind to */
var httpPort = flag.Int("port", 11938, "Port number to bind to")
/* @return set boolean will be true if argument is not nil */
func GetHTTPPort() (port int, set bool) {
	if httpPort == nil {return -1, false}
	return *httpPort, true;
}

/* Listen on a UNIX socket */
var socketLocation = flag.String("socket", "", "Location of UNIX socket to listen on. Setting will disable TCP.")
/* @return set boolean will be true if argument is not nil */
func GetSocketLocation() (port string, set bool) {
	if socketLocation == nil {return "", false}
	return *socketLocation, true;
}

/* Where the signal-cli binary is */
var binaryLocation = flag.String("binary", "/usr/local/bin/signal-cli", "Location of the signal-cli binary.")
/* @return set boolean will be true if argument is not nil */
func GetBinaryLocation() (port string, set bool) {
	if binaryLocation == nil {return "", false}
	return *binaryLocation, true;
}