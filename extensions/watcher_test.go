// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	ext "github.com/topfreegames/mystack-router/extensions"
	"github.com/topfreegames/mystack-router/models"
	"github.com/topfreegames/mystack-router/nginx"
	mystackTest "github.com/topfreegames/mystack-router/testing"
)

var _ = Describe("Watcher", func() {
	var err error

	Describe("GetMyStackServices", func() {
		It("should return empty list of services", func() {
			services, err := watcher.GetMyStackServices()

			Expect(err).NotTo(HaveOccurred())
			Expect(services.Items).To(BeEmpty())
		})

		It("should return list with one element after create service", func() {
			_, err = mystackTest.CreateService(clientset)
			Expect(err).NotTo(HaveOccurred())

			services, err := watcher.GetMyStackServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(services.Items).To(HaveLen(1))
		})
	})

	Describe("Build", func() {
		It("should create RouterConfig with empty AppConfigs", func() {
			routerConfig, err := watcher.Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(routerConfig.AppConfigs).To(BeEmpty())
		})

		It("should have list with one element after create service", func() {
			_, err = mystackTest.CreateService(clientset)
			Expect(err).NotTo(HaveOccurred())

			routerConfig, err := watcher.Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(routerConfig.AppConfigs).To(HaveLen(1))
		})
	})

	Describe("CreateConfigFile", func() {
		It("should create a config file", func() {
			fs := models.NewMockFS(nil)
			err := watcher.CreateConfigFile(fs)
			Expect(err).NotTo(HaveOccurred())

			exists, err := afero.DirExists(fs.AppFS, ext.NginxConfigDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			exists, err = afero.Exists(fs.AppFS, ext.NginxConfigFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should return error", func() {
			fs := models.NewMockFS(fmt.Errorf("error"))
			err := watcher.CreateConfigFile(fs)
			Expect(err).To(HaveOccurred())

			exists, err := afero.DirExists(fs.AppFS, ext.NginxConfigDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())

			exists, err = afero.Exists(fs.AppFS, ext.NginxConfigFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})

	Describe("Start", func() {
		It("should start nginx and config file", func() {
			ng := &nginx.Mock{}
			fs := models.NewMockFS(nil)
			watcher, err := ext.NewWatcher(config, clientset)
			Expect(err).NotTo(HaveOccurred())

			timeout := time.After(1 * time.Second)
			go watcher.Start(ng, fs)

			select {
			case <-timeout:
				exists, err := afero.DirExists(fs.AppFS, ext.NginxConfigDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())

				exists, err = afero.Exists(fs.AppFS, ext.NginxConfigFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())

				exists, err = afero.Exists(fs.AppFS, fmt.Sprintf("%s/conf.d", ext.NginxConfigDir))
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())
			}
		})

		It("should return error from nginx", func() {
			ng := &nginx.Mock{Err: fmt.Errorf("error")}
			fs := models.NewMockFS(nil)
			watcher, err := ext.NewWatcher(config, clientset)
			Expect(err).NotTo(HaveOccurred())

			timeout := time.After(1 * time.Second)
			go watcher.Start(ng, fs)

			select {
			case <-timeout:
				exists, err := afero.DirExists(fs.AppFS, ext.NginxConfigDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())

				exists, err = afero.Exists(fs.AppFS, ext.NginxConfigFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())

				exists, err = afero.Exists(fs.AppFS, fmt.Sprintf("%s/conf.d", ext.NginxConfigDir))
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())
			}
		})

		It("should return error from filesystem", func() {
			ng := &nginx.Mock{}
			fs := models.NewMockFS(fmt.Errorf("error"))
			watcher, err := ext.NewWatcher(config, clientset)
			Expect(err).NotTo(HaveOccurred())

			timeout := time.After(1 * time.Second)
			go watcher.Start(ng, fs)

			select {
			case <-timeout:
				Expect(err).NotTo(HaveOccurred())

				exists, err := afero.DirExists(fs.AppFS, ext.NginxConfigDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())

				exists, err = afero.Exists(fs.AppFS, ext.NginxConfigFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())

				exists, err = afero.Exists(fs.AppFS, fmt.Sprintf("%s/conf.d", ext.NginxConfigDir))
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())
			}
		})
	})
})
