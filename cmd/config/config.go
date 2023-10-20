package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigType = "yaml"
	ConfigName = ".etf-cli-config"
	etfsKey    = "etfs"
)

// InitConfig initializes the configuration.
func InitConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)
	viper.SetDefault(etfsKey, []string{"VWCE.DE", "VWCE.MI"})
	viper.SafeWriteConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("No config file could be found.\n Falling back to default config file '%s'", viper.ConfigFileUsed())
	}
}

// AddEtf adds an etf symbol to the configuration.
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

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}
