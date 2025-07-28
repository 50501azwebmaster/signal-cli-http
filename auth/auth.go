package auth

/* This file contains the AuthAuthConfig object and its methods, which handle reading
   from a config file and matching requests to the whitelist. */

import (
	"errors"
	"os"
)

/* Stores a map between a string (bearer token) and a list of unmarshaled JSONS */
var authConfig any;
var authConfigSetup bool = false;

/* Opens and reads a file at the path */
func SetupAuthConfig(filePath string) (err error) {
	if authConfigSetup {return errors.New("Auth configuration already set up!")}
	
	// Open and read file contents
	fileContents, err := os.ReadFile(filePath);
	if err != nil {return}
	
	// Unmarshal
	authConfig = unmarshalJSON(fileContents);
	if authConfig == nil {return errors.New("Invalid JSON config!");}
	
	print(match(authConfig, authConfig), "\n")
	
	// Finish setup
	authConfigSetup = true;
	return nil;
}

/* Gets a reference copy to the config data */
func GetAuthConfigData() (any, bool) {
	return authConfig, authConfigSetup;
}