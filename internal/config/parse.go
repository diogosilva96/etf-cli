package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Parse reads the file in the specified path and returns its parsed Config.
func Parse(path string) (*Config, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New(fmt.Sprintf("The file '%s' does not exist.", path))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be read. Details: %s", path, err))
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be parsed. Details: %s", path, err))
	}

	return &config, nil
}
