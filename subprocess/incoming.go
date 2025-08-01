package subprocess

/* This file manages incoming messages*/

import (
	"fmt"
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

/* Stores incoming messages in a room up to at least 15 minutes old 10,000
   This is intentionally an array of pointers so cache-clearing is faster */
var incomingMessageCache []*IncomingMessage;
var incomingMessageCacheLock sync.RWMutex;

func GetIMC() []*IncomingMessage {return incomingMessageCache}

/* Handler for incoming JSON objects which have "method":"receive" */
func handleIncoming(body string, unmarshaledJSONMap map[string]any) (ok bool) {
	if val, ok := unmarshaledJSONMap["method"]; !ok || val != "receive" {return false}
	fmt.Println(body)
	
	var newMessage *IncomingMessage = newIncomingMessage();
	// Using time.Now to ensure that pre/post-dated messages don't have issue
	newMessage.receivedTime = time.Now();
	newMessage.body = body;
	newMessage.unmarshaledJSONMap = unmarshaledJSONMap;
	
	// Obtain read-write lock
	incomingMessageCacheLock.Lock();
	incomingMessageCache = append(incomingMessageCache, newMessage);
	incomingMessageCacheLock.Unlock();
	return true;
}

/* Handles clearing space in incomingMessageCache */
func main() {go cacheClear()}

/* Runs in an infinite loop to try and clear the cache when needed */
func cacheClear() {
	for {
		// More than reasonable delay
		time.Sleep(time.Millisecond);
		// Only clear when it's 1000 items over
		if len(incomingMessageCache) < 1000 {continue}
		
		incomingMessageCacheLock.Lock();
		
		// Don't clear anything after this time
		fifteenMinutesAgo := time.Now().Add(-15 * time.Minute);
		
		// Find index in incomingMessageCache that is closest above 15 minutes ago
		i := sort.Search(len(incomingMessageCache), func(i int) bool {
			return incomingMessageCache[i].receivedTime.After(fifteenMinutesAgo)
		})
		
		incomingMessageCache = incomingMessageCache[i:]
		
		incomingMessageCacheLock.Unlock();
	}
}

/* Returns a list of encoded JSON strings from incomingMessageCache that match
   the filter from */
func GetIncoming(filter map[string]any) string {
	var list []string;
	// Create copy of incomingMessageCache for efficency
	incomingMessageCacheLock.RLock();
	incomingMessageCacheCopy := incomingMessageCache;
	incomingMessageCacheLock.RUnlock();
	
	// Create list of messages that match the filter
	for _, message := range incomingMessageCacheCopy {
		if !auth.Match(message.unmarshaledJSONMap, filter) {continue}
		fmt.Println(message.body)
		list = append(list, message.body)
	}
	
	var encoded string = "["
	for index, object := range list {
		encoded += object
		if index == len(list) - 1 {continue}
		encoded += ","
	}
	encoded += "]"
	
	return encoded;
}