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
	
	// Wait for request to return
	for {
		time.Sleep(time.Millisecond)
		waitingJobsMutex.RLock();
		if waitingJobs[id] != nil {
			waitingJobsMutex.RUnlock();
			break;
		}
		waitingJobsMutex.RUnlock();
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
func handleResponse(body string, unmarshaledJSONMap map[string]any) {
	val, ok := unmarshaledJSONMap["id"];
	if !ok {return}
	id, ok := val.(string)
	if !ok {return}
	
	// Read-Write lock the mutex
	waitingJobsMutex.Lock();
	defer waitingJobsMutex.Unlock();
	
	// Skip storage if there isn't a request for this ID
	if _, ok := waitingJobs[id]; !ok {return}
	// Store response in waiting Jobs
	waitingJobs[id] = &body;
}

/* Helper function to generate a random ID */
func genID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {return ""}
	return base64.RawURLEncoding.EncodeToString(b)
}