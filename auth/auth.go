package auth

/* This file contains the AuthAuthConfig object and its methods, which handle
   reading from a config file and matching requests to the whitelist. */

import (
	"errors"
	"os"
)

/* Stores a map between a string (bearer token) and a list of unmarshaled JSONS */
var authConfig map[string][]any = make(map[string][]any);
var authConfigSetup bool = false;

/* Opens, reads, and parses a file at the path */
func SetupAuthConfig(filePath string) (err error) {
	if authConfigSetup {return errors.New("Auth configuration already set up!")}
	
	// Open and read file contents
	fileContents, err := os.ReadFile(filePath);
	if err != nil {return}
	
	// Unmarshal
	unmarshaled := UnmarshalJSON(fileContents);
	if unmarshaled == nil {return errors.New("Invalid JSON object in config file!");}
	
	// Check type assertion for base JSON object
	if _, ok :=  unmarshaled.(map[string]any); !ok {
		return errors.New("JSON is incorrect format");
	}
	
	// Loop through each bearer key
	for key, val := range unmarshaled.(map[string]any) {
		// Check type assertion
		if _, ok :=  val.([]any); !ok {
			return errors.New("JSON is incorrect format for key " + key);
		}
		
		// Copy over array
		authConfig[key] = val.([]any);
	}
	
	// Finish setup
	authConfigSetup = true;
	return nil;
}

/* Gets a copy to the config data */
func GetAuthConfigData() (map[string][]any, bool) {
	return authConfig, authConfigSetup;
}

/* Returns true iff bearer is authorized for this request JSON */
func Authenticate(bearer string, requestJSON []byte) bool {
	// Check if bearer token exists at all
	if _, ok := authConfig[bearer]; !ok {return false;}
	
	// Unmarshal JSON
	unmarshaledRequest := UnmarshalJSON(requestJSON);
	
	// Check for any object
	for _, jsonObject := range authConfig[bearer] {
		if match(unmarshaledRequest, jsonObject) {return true}
	}
	
	return false;
}