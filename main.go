package main

/* This is the main function of this program. It handles setting everything up */

import (
	"signal-cli-http/args"
	"signal-cli-http/conf"
	"signal-cli-http/http"
	
	"log"
)

func main() {
	// Read arguments
	args.Parse();
	configLocation, confLocationSet := args.GetConfLocation();
	if !confLocationSet {log.Default().Print("No config value!"); return}
	log.Default().Print("Reading config value from ", configLocation);
	
	// Set up config 
	conf.GlobalConfig, _ = conf.NewConfig(configLocation);
	if conf.GlobalConfig == nil {log.Default().Print("Error reading config"); return}
	
	port, portSet := args.GetHTTPPort();
	if !portSet {log.Default().Print("No port value!"); return;}
	http.StartWebserver(port)
}