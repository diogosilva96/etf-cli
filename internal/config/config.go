package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// The default config.
	// This is used mostly for fallback purposes.
	DefaultConfig Config
)

type Config struct {
	Symbols []string
}

// Contains checks whether the Config contains the named symbol.
func (config *Config) Contains(symbol string) bool {
	for _, s := range config.Symbols {
		if strings.EqualFold(s, symbol) {
			return true
		}
	}
	return false
}

// Adds the named symbol to the config.
func (config *Config) Add(symbol string) error {
	if len(strings.TrimSpace(symbol)) == 0 {
		return errors.New("Cannot add empty symbol.")
	}

	config.Symbols = append(config.Symbols, symbol)
	return nil
}

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

func init() {
	DefaultConfig = Config{Symbols: []string{"VWCE.DE"}}
}
