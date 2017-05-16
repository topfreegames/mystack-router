// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	. "github.com/topfreegames/mystack-router/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockFs", func() {
	Describe("NewMockFS", func() {
		It("should create a new mock_fs", func() {
			fs := NewMockFS(nil)
			Expect(fs.AppFS).NotTo(BeNil())
		})
	})

	Describe("MkdirAll", func() {
		It("should create directory and its subdirectories", func() {
			fs := NewMockFS(nil)
			err := fs.MkdirAll("/new/long/path", os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error", func() {
			fs := NewMockFS(fmt.Errorf("error"))
			err := fs.MkdirAll("/new/long/path", os.ModePerm)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Create", func() {
		var fs *MockFS

		BeforeEach(func() {
			fs = NewMockFS(nil)
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

			exists, err := afero.Exists(fs.AppFS, name)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should return error", func() {
			fs := NewMockFS(fmt.Errorf("error"))
			name := "/etc/new/sub/dir/file"
			_, err := fs.Create(name)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("RemoveAll", func() {
		var fs *MockFS

		BeforeEach(func() {
			fs = NewMockFS(nil)
		})

		It("should remove file if exists", func() {
			name := "/etc/file"
			_, err := fs.Create(name)
			Expect(err).NotTo(HaveOccurred())

			exists, err := afero.Exists(fs.AppFS, name)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			err = fs.RemoveAll(name)
			Expect(err).NotTo(HaveOccurred())

			exists, err = afero.Exists(fs.AppFS, name)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should returne error", func() {
			fs = NewMockFS(fmt.Errorf("error"))
			name := "/etc/file"
			err := fs.RemoveAll(name)
			Expect(err).To(HaveOccurred())
		})
	})
})
