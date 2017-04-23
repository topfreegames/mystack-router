// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions_test

import (
	"database/sql"

	"k8s.io/client-go/kubernetes/fake"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	. "github.com/topfreegames/mystack-router/extensions"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"fmt"
	"strings"
	"testing"
)

var config *viper.Viper
var watcher *Watcher
var db *sql.DB
var mock sqlmock.Sqlmock
var clusterName string = "MyCustomApps"

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

var _ = BeforeEach(func() {
	var err error
	db, mock, err = sqlmock.New()
	Expect(err).NotTo(HaveOccurred())
	watcher = &Watcher{
		DB:               sqlx.NewDb(db, "postgres"),
		KubeDomainSuffix: "mystack.com",
		KubeClientSet:    fake.NewSimpleClientset(),
	}
})

var _ = AfterEach(func() {
	defer db.Close()
	err := mock.ExpectationsWereMet()
	Expect(err).NotTo(HaveOccurred())
})
