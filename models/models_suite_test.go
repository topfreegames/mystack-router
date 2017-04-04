package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

var config *viper.Viper

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Suite")
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
