package encryption

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Keydir string
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
	return configuration, nil
}
