// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"fmt"
	"strings"
	"testing"
)

var config *viper.Viper

func TestExtensions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extensions Suite")
}

var _ = BeforeSuite(func() {
	configFile := "../config/test.yaml"
	config = viper.New()
	config.SetConfigFile(configFile)
	config.SetConfigType("yaml")
	config.SetEnvPrefix("MYSTACK")
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	// If a Config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Config file %s failed to load: %s.\n", configFile, err.Error())
		panic("Failed to load Config file")
	}
})