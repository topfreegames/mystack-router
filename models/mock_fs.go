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

//MockFS implements FileSystem interface
type MockFS struct {
	AppFS afero.Fs
}

//NewMockFS constructs a new mock
func NewMockFS() *MockFS {
	return &MockFS{
		AppFS: afero.NewMemMapFs(),
	}
}

//MkdirAll creates a mock directory
func (m *MockFS) MkdirAll(path string, perm os.FileMode) error {
	return m.AppFS.MkdirAll(path, perm)
}

//Create creates a mock file
func (m *MockFS) Create(name string) (afero.File, error) {
	return m.AppFS.Create(name)
}
