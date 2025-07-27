package args

/* This file manages declaration of, parsing of, and returning copies of
   command line arguments. */

import "flag"

/* Method to trigger flag parsing. Can be called multiple times safely. */
func Parse() {
	if (flagsParsed) {return}
	flag.Parse();
	flagsParsed = true;
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