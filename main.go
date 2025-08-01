package main

/* This is the main function of this program. It handles setting everything up */

import (
	"fmt"
	"signal-cli-http/args"
	"signal-cli-http/auth"
	"signal-cli-http/subprocess"
	"signal-cli-http/web"
	"time"

	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup; wg.Add(1);
	
	// Read arguments
	args.Parse();
	authLocation, ok := args.GetAuthJson();
	if !ok {
		log.Default().Print("No auth config value!");
		return;
	}
	log.Default().Print("Reading auth config value from ", authLocation);
	
	// Set up authentication
	err := auth.SetupAuthConfig(authLocation);
	if err != nil {log.Default().Print("Error reading config: ", err); return;}
	log.Default().Print("Read auth config data");
	
	// Setup Subprocess
	binary, binarySet := args.GetBinaryLocation();
	if !binarySet {log.Default().Print("Read auth config data"); return;};
	err = subprocess.SetupCMD(binary);
	if err != nil {log.Default().Print("Error running subprocess at ", binary, ": ", err); return;};
	log.Default().Print("Started subprocess at ", binary);
	
	// HTTP Listen
	port, portSet := args.GetHTTPPort();
	if !portSet {log.Default().Print("No port value!"); return;}
	go func() {
		defer wg.Done();
		log.Default().Print("Error with web server: ", web.StartWebserver(port));
	}()
	log.Default().Print("Listening on port ", port);
	
	log.Default().Print("Startup tasks complete!");
	
	time.Sleep(time.Millisecond * 500);
	fmt.Println(subprocess.GetIMC())
	
	wg.Wait();
}