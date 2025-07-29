package web

/* This file handles listening to HTTP requests */

import (
	"signal-cli-http/auth"
	"signal-cli-http/subprocess"
	
	"fmt"
	"io"
	"log"
	"net/http"
)

func StartWebserver(port int) error {
	http.HandleFunc("/", getRoot);
	return http.ListenAndServe(":"+fmt.Sprint(port), nil);
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	// Check that Authentication header exists
	authArr, ok := r.Header["Authentication"]
	if (!ok) || (len(authArr) == 0) {
		w.WriteHeader(400);
		w.Write([]byte("Authentication header missing\n"))
		return;
	}
	bearer := authArr[0];
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500);
		w.Write([]byte("Error reading body\n"));
		return;
	}
	// Check Authentication header
	if !auth.Authenticate(bearer, body) {
		w.WriteHeader(403);
		w.Write([]byte("Bearer key not whitelisted for this request type\n"));
		return;
	}
	// Attempt to unmarshal JSON
	bodyUnmarshaled := auth.UnmarshalJSON(body);
	if bodyUnmarshaled == nil {
		w.WriteHeader(400);
		w.Write([]byte("Body content is not a valid JSON"));
		return;
	}
	// Type assertion
	b, ok := bodyUnmarshaled.(map[string]any)
	if !ok {
		w.WriteHeader(400);
		w.Write([]byte("Body content is not of the write format"));
		return;
	}
	// Run request
	bodyContent, err := subprocess.Request(b)
	if err != nil {
		w.WriteHeader(500);
		w.Write([]byte("Internal server error: " + err.Error() + "\n"));
		return
	}
	
	w.WriteHeader(200);
	w.Write([]byte(bodyContent));
	
	// Log the request
	log.Default().Print("HTTP Request: ", bearer, " " , 200, " ", string(body))
}