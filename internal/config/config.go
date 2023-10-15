package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Represents a config.
type Config struct {
	Symbols []string `json:"symbols"`
}

// Represents a config option.
type ConfigOption func(c *Config) error

// HasSymbol checks whether the Config has the named symbol.
func (config *Config) HasSymbol(symbol string) bool {
	for _, s := range config.Symbols {
		if strings.EqualFold(s, symbol) {
			return true
		}
	}
	return false
}

// AddSymbol adds the named symbol to the config.
func (config *Config) AddSymbol(symbol string) error {
	if len(strings.TrimSpace(symbol)) == 0 {
		return errors.New("Cannot add empty symbol.")
	}

	if config.HasSymbol(symbol) {
		return nil
	}

	config.Symbols = append(config.Symbols, symbol)
	return nil
}

// RemoveSymbol removes the named symbol from the config.
func (config *Config) RemoveSymbol(symbol string) error {
	if len(strings.TrimSpace(symbol)) == 0 {
		return nil
	}

	for i, s := range config.Symbols {
		if strings.EqualFold(s, symbol) {
			config.Symbols = append(config.Symbols[:i], config.Symbols[i+1:]...)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("The symbol '%s' does not exist in the tracked list.", symbol))
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

// NewConfig creates a new Config based on the specified options.
func NewConfig(options ...ConfigOption) (*Config, error) {
	c := &Config{}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithSymbols adds the specified symbols to the Config.
func WithSymbols(symbols ...string) ConfigOption {
	return func(c *Config) error {
		for _, s := range symbols {
			err := c.AddSymbol(s)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
