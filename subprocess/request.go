package subprocess

/* This file manages verifying/sanatizing the request and response and forwarding it to the subprocess */

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

/* Stores a map between job ID and a string pointer. Nil for job not done yet */
var waitingJobs map[string]*string = make (map[string]*string);
var waitingJobsMutex sync.RWMutex;

/* Method which module http calls to forward the request to the subprocess*/
func Request(body map[string]any) (responseJSON string, err error) {
	if _, ok := body["id"]; ok {
		err = errors.New("Request cannot contain id!");
		return
	}
	
	// Generate ID and store it
	var id string = genID();
	if len(id) == 0 {
		err = errors.New("Error generating ID!");
		return
	}
	body["id"] = id;
	// Also set JSONRPC-2.0
	body["jsonrpc"] = "2.0";
	
	// Marshal JSON to bytes
	contents, err := json.Marshal(body)
	if err != nil {return "", err}
	
	// Lock job when enqueueing
	waitingJobsMutex.Lock();
	writeCMD(string(contents));
	waitingJobs[id] = nil;
	waitingJobsMutex.Unlock();
	
	// Wait for request to finish
	for waitingJobs[id] == nil {
		time.Sleep(time.Millisecond)
	}
	
	// Lock job when dequeueing
	waitingJobsMutex.Lock();
	responseJSON = *waitingJobs[id];
	delete(waitingJobs, id)
	waitingJobsMutex.Unlock();
	
	err = nil;
	return;
}

/* Handles putting a response into waitingJobs. Returns false on error */
func Response(body string) {
	var unmarshaledJSON any;
	json.Unmarshal([]byte(body), &unmarshaledJSON);
	unmarshaledJSONMap, ok := unmarshaledJSON.(map[string]any)
	if !ok {return}
	
	val, ok := unmarshaledJSONMap["id"];
	if !ok {return}
	id, ok := val.(string)
	if !ok {return}
	
	
	// Write response into 
	waitingJobsMutex.Lock();
	waitingJobs[id] = &body;
	waitingJobsMutex.Unlock();
	
	
}

func genID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {return ""}
	return base64.RawURLEncoding.EncodeToString(b)
}