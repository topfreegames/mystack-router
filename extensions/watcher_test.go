// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/topfreegames/mystack-router/extensions"
	"github.com/topfreegames/mystack-router/models"
	mystackTest "github.com/topfreegames/mystack-router/testing"

	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Watcher", func() {
	var fakeClientset *fake.Clientset
	var watcher *Watcher
	var err error

	BeforeEach(func() {
		fakeClientset = fake.NewSimpleClientset()
		watcher, err = NewWatcher(config, fakeClientset)

		Expect(err).NotTo(HaveOccurred())
	})

	Describe("NewWatcher", func() {
		It("should construct a new watcher", func() {
			_, err := NewWatcher(config, fakeClientset)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetMyStackServices", func() {
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

	Describe("Build", func() {
		It("should create RouterConfig with empty AppConfigs", func() {
			customDomains := &models.MockCustomDomains{
				ControllerServiceName: "mystack-controller",
			}
			routerConfig, err := watcher.Build(customDomains)
			Expect(err).NotTo(HaveOccurred())
			Expect(routerConfig.AppConfigs).To(BeEmpty())
		})

		It("should have list with one element after create service", func() {
			_, err = mystackTest.CreateService(fakeClientset)
			Expect(err).NotTo(HaveOccurred())

			customDomains := &models.MockCustomDomains{
				ControllerServiceName: "mystack-controller",
			}

			routerConfig, err := watcher.Build(customDomains)
			Expect(err).NotTo(HaveOccurred())
			Expect(routerConfig.AppConfigs).To(HaveLen(1))
		})
	})
})
