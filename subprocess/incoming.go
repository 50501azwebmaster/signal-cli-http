package subprocess

/* This file manages incoming messages, marked with "method":"receive" in the JSON */

import (
	"signal-cli-http/auth"
	"sort"
	"sync"
	"time"
)

/* Stores an incoming message */
type IncomingMessage struct {
	/* Stores the time the message came in */
	receivedTime time.Time;
	/* Stores the message body */
	body string;
	/* Stores the unmarshaledJSON */
	unmarshaledJSONMap map[string]any;
}

func newIncomingMessage() *IncomingMessage {return &IncomingMessage{}}

/* Stores unlimited incoming messages up to 5 minutes old.
   This is intentionally an array of pointers so cache-clearing and appending
   don't require a copy of unmarshaledJSONMap. */
var incomingMessageCache []*IncomingMessage;
/* Locks incomingMessageCache */
var incomingMessageCacheLock sync.RWMutex;

/* Here exclusively for testing purposes */
func GetIMC() []*IncomingMessage {return incomingMessageCache}

/* Handler for incoming JSON objects which have "method":"receive" */
func handleIncoming(body string, unmarshaledJSONMap map[string]any) (ok bool) {
	// Check that the message's method is "receive"
	if val, ok := unmarshaledJSONMap["method"]; !ok || val != "receive" {return false}
	
	// Create new message structure
	var newMessage *IncomingMessage = newIncomingMessage();
	// Using time.Now() to ensure that pre/post-dated messages don't have issue
	newMessage.receivedTime = time.Now();
	newMessage.body = body;
	newMessage.unmarshaledJSONMap = unmarshaledJSONMap;
	
	// Add message into cache
	incomingMessageCacheLock.Lock();
	incomingMessageCache = append(incomingMessageCache, newMessage);
	incomingMessageCacheLock.Unlock();
	return true;
}

/* Handles clearing space in incomingMessageCache */
func StartCacheClear() {go cacheClear()}
/* Runs in an infinite loop to try and clear the cache when needed */
func cacheClear() {
	for {
		// More than reasonable delay
		time.Sleep(time.Millisecond * 25);
		// Only attempt to clear when it's 1000 items or more
		if len(incomingMessageCache) < 1000 {continue}
		
		incomingMessageCacheLock.Lock();
		
		fiveMinutesAgo := time.Now().Add(time.Minute*(-5));
		
		// Find first index in incomingMessageCache that is at most 15 minutes old
		i := sort.Search(len(incomingMessageCache), func(i int) bool {
			return incomingMessageCache[i].receivedTime.After(fiveMinutesAgo)
		})
		incomingMessageCache = incomingMessageCache[i:]
		
		incomingMessageCacheLock.Unlock();
	}
}

/* Returns a list of encoded JSON strings from incomingMessageCache that match
   the filter from
   @return always valid JSON array object. Can be empty. */
func GetIncoming(filter map[string]any) string {
	// Create copy of incomingMessageCache as the following loop can be slow
	incomingMessageCacheLock.RLock();
	incomingMessageCacheCopy := incomingMessageCache;
	incomingMessageCacheLock.RUnlock();
	
	// Create list of messages that match the filter
	var list []string;
	for _, message := range incomingMessageCacheCopy {
		if !auth.Match(message.unmarshaledJSONMap, filter) {continue}
		list = append(list, message.body)
	}
	
	// Constructs the JSON string without considering the JSON object
	var encoded string = "["
	for index, object := range list {
		encoded += object
		if index == len(list) - 1 {continue}
		encoded += ","
	}
	encoded += "]\n"
	
	return encoded;
}