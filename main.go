package main

/* This is the main function of this program. It handles setting everything up */

import (
	"signal-cli-http/args"
	"signal-cli-http/auth"
	"signal-cli-http/web"
	
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup; wg.Add(1);
	
	// Read arguments
	args.Parse();
	configLocation, confLocationSet := args.GetAuthJson();
	if !confLocationSet {log.Default().Print("No auth config value!"); return;}
	log.Default().Print("Reading auth config value from ", configLocation);
	
	// Set up config 
	err := auth.SetupAuthConfig(configLocation);
	if err != nil {log.Default().Print("Error reading config: ", err); return;}
	log.Default().Print(auth.GetAuthConfigData());
	
	port, portSet := args.GetHTTPPort();
	if !portSet {log.Default().Print("No port value!"); return;}
	log.Default().Print("Listening on port ", port);
	
	go func() {
		defer wg.Done();
		web.StartWebserver(port);
	}()
	
	log.Default().Print("Startup tasks complete!");
	wg.Wait();
}