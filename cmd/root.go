// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config *viper.Viper

// ConfigFile is the configuration file used for running a command
var ConfigFile string

// Execute runs RootCmd to initialize mystack CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mystack-watcher",
	Short: "mystack watcher for kubernetes services that updates the ingress controller",
	Long:  `mystack watcher for kubernetes services that updates the ingress controller`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&ConfigFile, "config", "c", "./config/local.yaml",
		"config file (default is ./config/local.yaml)",
	)
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	config = viper.New()
	if ConfigFile != "" { // enable ability to specify config file via flag
		config.SetConfigFile(ConfigFile)
	}
	config.SetConfigType("yaml")
	config.SetEnvPrefix("MYSTACK")
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	config.AutomaticEnv()

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Config file %s failed to load: %s.\n", ConfigFile, err.Error())
		panic("Failed to load config file")
	}
}
