package main

/* This is the main function of this program. It handles setting everything up */

import (
	"signal-cli-http/args"
	"signal-cli-http/conf"

	"fmt"
	"log"
)

func main() {
	// Read arguments
	args.Parse();
	configLocation, confLocationSet := args.GetConfLocation();
	if !confLocationSet {log.Default().Print("No config value!"); return;}
	log.Default().Print("Reading config value from ", configLocation);
	
	// Set up config 
	config, err := conf.NewConfig(configLocation);
	if err != nil {log.Default().Print("Error reading config: ", err); return;}
	fmt.Println(config.GetConfigData())
}