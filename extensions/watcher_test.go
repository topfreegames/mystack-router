// mystack
// +build unit
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	. "github.com/topfreegames/mystack-router/extensions"
	mystackTest "github.com/topfreegames/mystack-router/testing"

	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Watcher", func() {

	const (
		nginxConfigDir      = "/etc/nginx"
		nginxConfigFilePath = nginxConfigDir + "/nginx.conf"
	)

	BeforeEach(func() {
		appFS := afero.NewOsFs()
		appFS.MkdirAll(nginxConfigDir, 0777)
		appFS.Create(nginxConfigFilePath)
	})

	Describe("NewWatcher", func() {
		It("should construct a new watcher", func() {
			fakeClientset := fake.NewSimpleClientset()
			_, err := NewWatcher(config, fakeClientset)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetMyStackServices", func() {
		var fakeClientset *fake.Clientset
		var watcher *Watcher
		var err error

		BeforeEach(func() {
			fakeClientset = fake.NewSimpleClientset()
			watcher, err = NewWatcher(config, fakeClientset)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return empty list of services", func() {
			services, err := watcher.GetMyStackServices()

			Expect(err).NotTo(HaveOccurred())
			Expect(services.Items).To(BeEmpty())
		})

		It("should return list with one element after create service", func() {
			_, err = mystackTest.CreateService(fakeClientset)
			Expect(err).NotTo(HaveOccurred())

			services, err := watcher.GetMyStackServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(services.Items).To(HaveLen(1))
		})
	})
})
