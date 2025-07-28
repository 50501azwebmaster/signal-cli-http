package http

/* This file handles listening to HTTP requests */

import (
	"signal-cli-http/conf"
	"signal-cli-http/subprocess"
	
	"fmt"
	"io"
	"log"
	"net/http"
)

func StartWebserver(port int) {
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":"+fmt.Sprint(port), nil)
	fmt.Println(err)
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
	
	// Check that the request is allowed for the path
	if !conf.GlobalConfig.ValidateBearerKey(bearer, r.URL.Path) {
		w.WriteHeader(403);
		w.Write([]byte("Bearer key not whitelisted for this path\n"))
		return;
	}
	
	// OK authentication wise
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500);
		w.Write([]byte("Error reading body\n"))
		return;
	}
	
	// Call subprocess
	status, bodyContent, err := subprocess.Run(r.URL.Path, body)
	
	// Error
	if err != nil {
		w.WriteHeader(500);
		w.Write([]byte("Internal server error: " + err.Error() + "\n"));
		return
	}
	
	// Respond to client with status
	w.WriteHeader(status);
	w.Write(bodyContent);
	
	// Log the request
	log.Default().Print("HTTP Request: ", bearer, " " , r.URL.Path, " ", status)
}