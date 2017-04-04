// mystack
// +build unit
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/topfreegames/mystack-router/models"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockFs", func() {
	Describe("NewMockFS", func() {
		It("should create a new mock_fs", func() {
			fs := NewMockFS()
			Expect(fs.AppFS).NotTo(BeNil())
		})
	})

	Describe("MkdirAll", func() {
		It("should create directory and its subdirectories", func() {
			fs := NewMockFS()
			err := fs.MkdirAll("/new/long/path", os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Create", func() {
		var fs *MockFS

		BeforeEach(func() {
			fs = NewMockFS()
		})

		It("should create a file if directory exists", func() {
			name := "/etc/file"
			file, err := fs.Create(name)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
			Expect(file.Name()).To(Equal(name))
		})

		It("should create a file if directory doesn't exist", func() {
			name := "/etc/new/sub/dir/file"
			file, err := fs.Create(name)
			Expect(file).NotTo(BeNil())
			Expect(file.Name()).To(Equal(name))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
