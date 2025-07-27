package conf

/* This file contains regex helper functions for parsing configs  */

import (
	"strings"
)

/* Splits and normalises */
func splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

/* Attempts to match a request path to a set of whitelisted paths
   @return false for anything other than a valid match */
func match(request string, matchTo string) bool {
	return matchSegments(splitPath(request), splitPath(matchTo));
}

/* Returns false for anything other than a valid match */
func matchSegments(request []string, matchTo []string) bool {
	/* This is a recursive function which, at each recursion level, matches the
	   path segments that are in the front of the request and matchTo lists
	   It matches identical strings, anything to &, and splits in two when
	   matching anything to *, to account for consuming or not consuming the *
	   at the current recursion level. */
	
	// Recursion base case for perfect match
	if len(request) == 0 && len(matchTo) == 0 {return true}
	// End of path for one but not the other
	if (len(request) & len(matchTo)) == 0 {return false}
	
	// Grab current path segments
	requestCurrent := request[0];
	matchToCurrent := matchTo[0];
	
	// & character and direct match have the same behavior
	if (matchToCurrent == "&") || (requestCurrent == matchToCurrent) {
		return matchSegments(request[1:], matchTo[1:]);
	}
	
	// * character
	if (matchToCurrent == "*") {
		// These are split for performance
		// Usually the * only refers to 1 or 2 things so putting consumption
		// first is probably a better choice
		if (matchSegments(request[1:], matchTo[1:])) {return true}
		if (matchSegments(request[1:], matchTo)) {return true}
		return false;
	}
	
	// Code will reach here if there's no match for the current segment
	return false;
}