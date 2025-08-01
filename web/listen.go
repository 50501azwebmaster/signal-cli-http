package web

/* This file handles listening to HTTP requests */

import (
	"encoding/json"
	"signal-cli-http/auth"
	"signal-cli-http/subprocess"

	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func StartWebserver(port int) error {
	http.HandleFunc("/", getRoot);
	return http.ListenAndServe(":"+fmt.Sprint(port), nil);
}

func writeLog(method string, status int, start time.Time) {
	duration := time.Now().Sub(start);
	log.Default().Printf("%s %d %s", method, status, duration.String())
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	
	// Check that Authentication header exists
	authArr, ok := r.Header["Authentication"]
	if (!ok) || (len(authArr) == 0) {
		w.WriteHeader(400);
		w.Write([]byte("Authentication header missing\n"))
		writeLog(r.Method, 400, startTime)
		return;
	}
	bearer := authArr[0];
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500);
		w.Write([]byte("Error reading body\n"));
		writeLog(r.Method, 500, startTime)
		return;
	}
	
	// Check Authentication header
	if !auth.Authenticate(bearer, body) {
		w.WriteHeader(403);
		w.Write([]byte("Bearer key not whitelisted for this request type\n"));
		writeLog(r.Method, 403, startTime)
		return;
	}
	
	// Attempt to unmarshal JSON
	var bodyUnmarshaled any;
	if err := json.Unmarshal(body, &bodyUnmarshaled); err != nil {
		w.WriteHeader(400);
		w.Write([]byte("Body content is not a valid JSON"));
		writeLog(r.Method, 400, startTime)
		return;
	}
	
	// Type assertion
	b, ok := bodyUnmarshaled.(map[string]any);
	if !ok {
		w.WriteHeader(400);
		w.Write([]byte("Body content is not of the write format"));
		writeLog(r.Method, 400, startTime)
		return;
	}
	
	// Handle incoming
	method, ok := b["method"];
	
	if method == "receive" {
		incoming := subprocess.GetIncoming(b)
		w.WriteHeader(200);
		w.Write([]byte(incoming));
		writeLog(r.Method, 200, startTime)
		return
	}
	
	// Run request
	bodyContent, err := subprocess.Request(b)
	if err != nil {
		w.WriteHeader(500);
		w.Write([]byte("Internal server error: " + err.Error() + "\n"));
		return
	}
	
	// Request returned something
	w.WriteHeader(200);
	w.Write([]byte(bodyContent));
	
	// Log the request
	log.Default().Print("HTTP Request: ", bearer, " " , 200, " ", string(body))
}