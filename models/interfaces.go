// mystack-router
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"database/sql"
	"os"

	"github.com/spf13/afero"
)

//DomainsPerApp holds the custom domains for each app or service
type DomainsPerApp map[string][]string

//FileSystem interface
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (afero.File, error)
}

//DB is the mystack-controller db interface
type DB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}
