package http

/* This file handles listening to HTTP requests */

import (
	"signal-cli-http/conf"
	
	"fmt"
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
	if (!ok) || (len(authArr) == 0) {w.WriteHeader(400); return}
	bearer := authArr[0];
	
	// Check that the request is allowed for the path
	if !conf.GlobalConfig.ValidateBearerKey(bearer, r.URL.Path) {
		w.WriteHeader(403);
		return;
	}
	
	log.Default().Print("HTTP Request: ", bearer, " " , r.URL.Path)
	
	// OK authentication wise
	w.WriteHeader(200);
}