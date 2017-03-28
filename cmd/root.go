// mystack api
// https://github.com/topfreegames/mystack
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config *viper.Viper

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mystack-watcher",
	Short: "mystack watcher for kubernetes services that updates the ingress controller",
	Long:  `mystack watcher for kubernetes services that updates the ingress controller`,
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/local.yaml", "config file (default is config/local.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config = viper.New()
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	}

	config.SetConfigName("")
	config.SetEnvPrefix("MYSTACK")
	config.AutomaticEnv()
	fmt.Printf("olha env %s\n", config.GetString("ola"))

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", config.ConfigFileUsed())
	}
}
