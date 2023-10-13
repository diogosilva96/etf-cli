package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Parse reads the file in the specified file path and returns its parsed Config.
func Parse(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New(fmt.Sprintf("The file '%s' does not exist.", filePath))
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be read. Details: %s", filePath, err))
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be parsed. Details: %s", filePath, err))
	}

	return &config, nil
}
