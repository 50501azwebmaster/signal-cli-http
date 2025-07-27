package conf

/* This file contains the Config object and its methods, which handle reading
   from a config file and matching requests to the whitelist. */

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

/* Object to handle what is in a JSON config */
type Config struct {
	configData map[string][]string;
}

/* Default Config object */
var GlobalConfig * Config;

/* Opens and reads a file at the path */
func NewConfig(filePath string) (newConfig *Config, err error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {return}
	defer file.Close()
	
	// Create configuration
	newConfigData := make(map[string][]string);
	
	// Read lines into newConfigData
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), " ", 2);
		if len(parts) != 2 {err = errors.New("Bad config file!"); return;}
		newConfigData[parts[0]] = append(newConfigData[parts[0]], parts[1]);
	}
	
	// Create Config object and copy a reference to newConfigData into it
	return &Config{configData: newConfigData}, nil;
}

/* Gets a reference copy to the config data */
func (config * Config) GetConfigData() map[string][]string {
	return config.configData;
}

/* Returns if a bearer key is authorized for the path in this Config object
   @return false for any situation that isn't a valid match */
func (config * Config) ValidateBearerKey(bearerKey string, request string) bool {
	paths, exists := config.configData[bearerKey];
	if !exists {return false}
	
	for _, matchTo := range paths {
		if match(request, matchTo) {
			return true;
		}
	}
	
	return false;
}