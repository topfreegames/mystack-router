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

//RealFS implements FileSystem interface
type RealFS struct{}

//NewRealFS constructs a new mock
func NewRealFS() *RealFS {
	return &RealFS{}
}

//MkdirAll creates a mock directory
func (m *RealFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

//Create creates a mock file
func (m *RealFS) Create(name string) (afero.File, error) {
	return os.Create(name)
}
