package subprocess

/* This file manages calling signal-cli */

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/creack/pty"
)

var cmd *exec.Cmd;
var cmdStarted = false;
var f *os.File;
var fLock sync.RWMutex;
var reader *bufio.Scanner;

// This is here to ignore lines written to STDIN echoed back through STDOUT
var ignoreEcho map[string]bool = make(map[string]bool);
var ignoreEchoMutex sync.RWMutex;

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
	
	// Make sure cmd is killed when this program is
	go catchKillSignal();
	
	// No problem
	return nil;
}

/* Catch signals from this program to kill child */
func catchKillSignal() {
	signalCatcher := make(chan os.Signal, 1)
	signal.Notify(signalCatcher, syscall.SIGINT, syscall.SIGTERM, syscall.SIGILL)
	<- signalCatcher;
	cmd.Process.Kill();
	os.Exit(0);  // Without this the process never exits
}

/* Continuously reads the next line up to 10MiB and forwards it to response */
func readCMD() {
	/* Holds temporary stdin data */
	buf := make([]byte, 16384);
	/* Holds the read data */
	var outBuffer string;
	
	for  {
		// Read up to 16384 bytes
		n, err := f.Read(buf)
		if err != nil {break;}
		if n == 0 {continue;}
		
		// Append read bytes to buffer
		outBuffer = outBuffer + string(buf[:n]);
		
		// Check for line ending
		newLineIndex := strings.Index(outBuffer, "\n");
		if newLineIndex < 0 {continue}
		
		// Extract line from buffer
		line := outBuffer[:newLineIndex];
		// Adjust for "\r\n" line ending
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		// And remove line from buffer
		outBuffer = outBuffer[newLineIndex+1:];
		
		// Check for echo
		ignoreEchoMutex.Lock();
		_, exists := ignoreEcho[line];
		if exists {delete(ignoreEcho, line)}
		ignoreEchoMutex.Unlock();
		if exists {continue}
		
		// Unmarshal the JSON
		var unmarshaledJSON any;
		if err := json.Unmarshal([]byte(line), &unmarshaledJSON); err != nil {continue}
		
		// Make sure it's a JSON map
		unmarshaledJSONMap, ok := unmarshaledJSON.(map[string]any)
		if !ok {continue}
		
		// Get method
		method, ok := unmarshaledJSONMap["method"];
		if ok && method == "receive" {
			handleIncoming(line, unmarshaledJSONMap);
			continue;
		}
	
		handleResponse(line, unmarshaledJSONMap);
	}
}

/* Write a line into the subprocess */
func writeCMD(line string) (ok bool) {
	// Write into ignoreEcho map so reader can skip the echoed line
	ignoreEchoMutex.Lock();
	ignoreEcho[line] = true;
	ignoreEchoMutex.Unlock();

	fLock.Lock();
	if line[len(line)-1] != '\n' {line += "\n"}
	f.WriteString(line);
	fLock.Unlock();
	
	return true;
}