package subprocess

/* This file manages creating the command line arguments to the subprocess */

/* Method which module http calls to create the subprocess */
func Run(path string, body []byte) (status int, bodyContents []byte, err error) {
	arguments := getArguments(path, body);
	
	// Don't know what to do with this request
	if arguments == nil {return 404, []byte("Unknown request\n"), nil;}
	
	// Call subprocess
	status, bodyContents, err = runCommand(arguments);
	
	return
}

/* Converts a request into the correct binary arguments */
func getArguments(path string, body []byte) []string {
	return nil // For now
}