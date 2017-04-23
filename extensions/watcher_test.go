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
			_, err = mystackTest.CreateService(watcher.KubeClientSet)
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
			_, err = mystackTest.CreateService(watcher.KubeClientSet)
			Expect(err).NotTo(HaveOccurred())

			routerConfig, err := watcher.Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(routerConfig.AppConfigs).To(HaveLen(1))
		})
	})
})
