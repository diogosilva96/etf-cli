package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	configType = "json"
	configName = ".etf-cli-config"
	etfsKey    = "etfs"
)

// InitConfig initializes the configuration.
func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(home)
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.SetDefault(etfsKey, []string{"VWCE.DE"})
	_ = viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		return errors.New("could not find config file")
	}

	return nil
}

// AddEtf adds an etf to the configuration.
func AddEtf(etf string) error {
	if len(strings.TrimSpace(etf)) == 0 {
		return errors.New(fmt.Sprintf("the input should not be empty"))
	}

	etfs := viper.GetStringSlice(etfsKey)
	if contains(etfs, etf) {
		return errors.New(fmt.Sprintf("the etf '%s' already exists in the configuration", etf))
	}

	etfs = append(etfs, etf)
	viper.Set(etfsKey, etfs)

	err := viper.WriteConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("an error occurred while adding etf '%s' to the configuration, details: %s", etf, err))
	}
	return nil
}

// RemoveEtf removes an etf from the configuration.
func RemoveEtf(etf string) error {
	if len(strings.TrimSpace(etf)) == 0 {
		return errors.New(fmt.Sprintf("the input should not be empty"))
	}

	etfs := viper.GetStringSlice(etfsKey)
	for i, e := range etfs {
		if strings.EqualFold(e, etf) {
			etfs = append(etfs[:i], etfs[i+1:]...)
			viper.Set(etfsKey, etfs)
			err := viper.WriteConfig()
			if err != nil {
				return errors.New(fmt.Sprintf("an error occurred while removing etf '%s' from the configuration, details: %s", etf, err))
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("the etf '%s' could not be found in the configuration", etf))
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
