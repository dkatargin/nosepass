package encryption

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Keydir string
}

func Config() (Configuration, error) {
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
