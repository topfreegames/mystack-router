// mystack-router
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/spf13/afero"
	"os"
)

//DomainsPerApp holds the custom domains for each app or service
type DomainsPerApp map[string][]string

//FileSystem interface
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (afero.File, error)
}

//CustomDomains interface
type CustomDomainsInterface interface {
	GetCustomDomains(controllerDomain, clusterName string) (DomainsPerApp, error)
}
