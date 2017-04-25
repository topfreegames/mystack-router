// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	db     *sql.DB
	sqlxDB *sqlx.DB
	mock   sqlmock.Sqlmock
	err    error
	config *viper.Viper
)

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

var _ = BeforeEach(func() {
	db, mock, err = sqlmock.New()
	Expect(err).NotTo(HaveOccurred())
	sqlxDB = sqlx.NewDb(db, "postgres")
})

var _ = AfterEach(func() {
	defer db.Close()
	err = mock.ExpectationsWereMet()
	Expect(err).NotTo(HaveOccurred())
})
