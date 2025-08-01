package auth

/* This file contains the AuthAuthConfig object and its methods, which handle
   reading from a config file and matching requests to the whitelist. */

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
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
	var unmarshaled any;
	if err := json.Unmarshal(fileContents, &unmarshaled); err != nil {
		return errors.New("Invalid JSON object in config file!");
	}
	
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
	var unmarshaledRequest any;
	if err := json.Unmarshal(requestJSON, &unmarshaledRequest); err != nil {
		return false;
	}
	
	// Check for any object
	for _, jsonObject := range authConfig[bearer] {
		if Match(unmarshaledRequest, jsonObject) {return true}
	}
	
	return false;
}

/* Meat and bones of determining if a request is allowed by a filter */
func Match(request any, filter any) bool {
	// Check that the types are the same
	if reflect.TypeOf(request) != reflect.TypeOf(filter) {return false}
	
	// Can safely switch on type of one object at this point since they're equal
	switch filter.(type) {
		
		// Key-value pairs
		case map[string]any:
			// Check for every key that's in the filter
			for key := range filter.(map[string]any) {
				// that it's in the request
				if _, ok := request.(map[string]any)[key]; !ok {return false}
				
				// And recursively check that the value is equal
				if !Match(request.(map[string]any)[key], filter.(map[string]any)[key]) {
					return false;
				}
			}
			return true;
		
		/* Arrays attempt to match every item in the filter to ANY item in the
		   request. Duplicates in the filter are treated as one */
		case []any:
			// Check that for every item in the filter
			for i := 0; i < len(filter.([]any)); i ++ {
				foundMatch := false;
				// That something matches in the request
				for j := 0; j < len(request.([]any)); j ++ {
					if Match(filter.([]any)[i], request.([]any)[j]) {
						foundMatch = true;
						break
					}
				}
				// Cannot find a match for something in the filter
				if !foundMatch {return false}
			}
			// And the other way around
			for i := 0; i < len(filter.([]any)); i ++ {
				foundMatch := false;
				// That something matches in the request
				for j := 0; j < len(request.([]any)); j ++ {
					if Match(filter.([]any)[i], request.([]any)[j]) {
						foundMatch = true;
						break
					}
				}
				// Cannot find a match for something in the filter
				if !foundMatch {return false}
			}
			
			return true;
		
		// Otherwise compare the objects directly using reflect
		default: return reflect.DeepEqual(request, filter);
	}
}