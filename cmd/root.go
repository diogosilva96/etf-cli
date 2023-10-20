package cmd

import (
	"fmt"
	"os"

	"github.com/diogosilva96/etf-scraper/cmd/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "etf",
	Short: "etf - a simple CLI to generate up to date ETF data reports & manage tracked ETFs.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Printf("")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType(config.ConfigType)
		viper.SetConfigName(config.ConfigName)
		viper.SetDefault("etfs", []string{"VWCE.DE", "VWCE.MI"})
		viper.SafeWriteConfig()
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "No config file could be found.\n Falling back to default config file:", viper.ConfigFileUsed())
	}
}
