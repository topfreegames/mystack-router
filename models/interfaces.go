// mystack-router
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"os"

	"github.com/spf13/afero"
)

//FileSystem interface
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (afero.File, error)
	RemoveAll(path string) error
}
