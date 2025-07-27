package conf

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
	line := scanner.Text()
		parts := strings.SplitN(line, " ", 2);
		if len(parts) != 2 {err = errors.New("Bad config file!"); return;}
		newConfigData[parts[0]] = append(newConfigData[parts[0]], parts[1]);
	}
	
	// Create Config object and copy a reference to newConfigData into it
	return &Config{configData: newConfigData}, nil;
}

func (config Config) GetConfigData() map[string][]string {return config.configData;}