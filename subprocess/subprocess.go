package subprocess

/* This file manages calling signal-cli */

import (
	"bufio"
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

func readCMD() {
	var maxCapacity int = 4096;
	buf := make([]byte, maxCapacity);
	reader.Buffer(buf, maxCapacity);
	
	for reader.Scan() {Response(reader.Text())}
}

func writeCMD(line string) (ok bool) {
	fLock.Lock();
	if line[len(line)-1] != '\n' {line += "\n"}
	f.WriteString(line);
	fLock.Unlock();
	return true;
}