package subprocess

/* This file manages calling signal-cli */

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

var cmd *exec.Cmd;
var cmdStarted = false;
var f *os.File;
var fLock sync.RWMutex;
var reader *bufio.Scanner;

func SetupCMD(binaryLocation string) error {
	// Avoid double set-up
	if cmdStarted {return errors.New("cmd already started")};
	
	// Create cmd object
	cmd = exec.Command(binaryLocation, "jsonRpc");
	if cmd == nil {return errors.New("Creating process failed!")}
	
	// Start it
	file, err := pty.Start(cmd);
	f = file;
	if err != nil {return err}
	
	// Set up reader object and loop
	reader = bufio.NewScanner(f);
	go readCMD();
	
	// No problem
	return nil;
}

/* Continuously reads the next line up to 40960 bytes and forwards it to response */
func readCMD() {
	var maxCapacity int = 40960;
	buf := make([]byte, maxCapacity);
	reader.Buffer(buf, maxCapacity);
	
	for reader.Scan() {
		// Read the line
		line := reader.Text();
		
		// Unmarshal the JSON
		var unmarshaledJSON any;
		if err := json.Unmarshal([]byte(line), &unmarshaledJSON); err != nil {continue}
		
		// Make sure it's a JSON map
		unmarshaledJSONMap, ok := unmarshaledJSON.(map[string]any)
		if !ok {continue}
		
		// Get method
		method, ok := unmarshaledJSONMap["method"];
		if !ok {continue}
		
		// Redirect to handlers based off method
		if method == "receive" {
			handleIncoming(line, unmarshaledJSONMap);
		} else {
			handleResponse(line, unmarshaledJSONMap);
		}
	}
}

/* Write a line into the subprocess */
func writeCMD(line string) (ok bool) {
	fLock.Lock();
	if line[len(line)-1] != '\n' {line += "\n"}
	f.WriteString(line);
	fLock.Unlock();
	return true;
}