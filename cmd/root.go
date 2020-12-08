package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/netgear-gs308e-exporter/config.yaml", "Config file path")
}

func initConfig() {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}
}

var rootCmd = &cobra.Command{
	Use:   "netgear-gs308e-exporter",
	Short: "netgear-gs308e-exporter",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
