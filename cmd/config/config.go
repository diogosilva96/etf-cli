package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

const (
	configType = "json"
	configName = "config"
	etfsKey    = "etfs"
)

// InitConfig initializes the configuration.
func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.SetDefault(etfsKey, []string{"VWCE.DE", "VWCE.MI"})
	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Could not find config file.\n")
	}
}

// AddEtf adds an etf to the configuration.
func AddEtf(etf string) error {
	if len(strings.TrimSpace(etf)) == 0 {
		return errors.New(fmt.Sprintf("The input should not be empty."))
	}

	etfs := viper.GetStringSlice(etfsKey)
	if contains(etfs, etf) {
		return errors.New(fmt.Sprintf("The etf '%s' already exists in the configuration.", etf))
	}

	etfs = append(etfs, etf)
	viper.Set(etfsKey, etfs)

	err := viper.WriteConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("An error occurred while adding etf '%s' to the configuration. Details: %s", etf, err))
	}
	return nil
}

// RemoveEtf removes an etf from the configuration.
func RemoveEtf(etf string) error {
	if len(strings.TrimSpace(etf)) == 0 {
		return errors.New(fmt.Sprintf("The input should not be empty."))
	}

	etfs := viper.GetStringSlice(etfsKey)
	for i, e := range etfs {
		if strings.EqualFold(e, etf) {
			etfs = append(etfs[:i], etfs[i+1:]...)
			viper.Set(etfsKey, etfs)
			err := viper.WriteConfig()
			if err != nil {
				return errors.New(fmt.Sprintf("An error occurred while removing etf '%s' from the configuration. Details: %s", etf, err))
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("The etf '%s' could not be found in the configuration.", etf))
}

// ListEtfs reads the etfs section from the configuration.
func ListEtfs() []string {
	return viper.GetStringSlice(etfsKey)
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}
