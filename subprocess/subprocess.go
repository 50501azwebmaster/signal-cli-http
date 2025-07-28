package subprocess

import (
	"signal-cli-http/args"

	"errors"
	"os/exec"
	"strings"
)

/* This file manages calling signal-cli */

/* Runs the command */
func runCommand(arguments []string) (returnStatus int, bodyContent []byte, err error) {
	// Get binary location
	binary, ok := args.GetBinaryLocation();
	if !ok {
		err = errors.New("Binary cannot be found!");
		return;
	}
	
	// Create command using binary location and arguments
	command := exec.Command(binary, strings.Join(arguments, " "));
	
	// Duplicate pointer into command.Stdout
	//command.Stdout = &bodyContent;

	// Run the command
	err = command.Run();
	if err != nil {return}
	
	// Extract exit code if possible
	if exitError, ok := err.(*exec.ExitError); ok {
		returnStatus = exitError.ExitCode();
	} else {
		err = errors.New("Cannot get exit code!");
	}
	
	// Get output
	bodyContent, err = command.Output();
	
	// Named return values allow this
	return;
}