package common

import (
	"encoding/json"
	"os"
	"strings"
)

type Configuration struct {
	DbDir  string
	KeyDir string
}

func Config() (Configuration, error) {
	// read config and represent to struct
	configuration := Configuration{}
	file, err := os.Open("settings.conf")
	if err != nil {
		return configuration, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, err
	}
	homeDir, _ := os.UserHomeDir()
	configuration.DbDir = strings.Replace(configuration.DbDir, "~", homeDir, -1)
	configuration.KeyDir = strings.Replace(configuration.KeyDir, "~", homeDir, -1)
	return configuration, nil
}
