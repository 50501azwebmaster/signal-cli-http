package auth

/* This file contains some JSON helper functions */

import (
	"encoding/json"
	"reflect"
)

/* Unmarshals a JSON into a recursive map. Returns nil on error */
func UnmarshalJSON(marshaledJSON []byte) (unmarshaled any) {
	json.Unmarshal(marshaledJSON, &unmarshaled);
	return;
}

/* Meat and bones of determining if a request is allowed by a filter */
func match(request any, filter any) bool {
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
				if !match(request.(map[string]any)[key], filter.(map[string]any)[key]) {
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
					if match(filter.([]any)[i], request.([]any)[j]) {
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
					if match(filter.([]any)[i], request.([]any)[j]) {
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