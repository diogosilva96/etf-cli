package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Save writes the config data to the named file path, creating it if necessary.
func (config *Config) Save(filePath string) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return errors.New(fmt.Sprintf("An error occurred while marshalling the config. Details:%s", err))
	}

	err = os.WriteFile(filePath, data, os.ModeAppend)

	if err != nil {
		return errors.New(fmt.Sprintf("An error occurred while saving the config file. Details:%s", err))
	}
	return nil
}
