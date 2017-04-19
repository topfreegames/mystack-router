// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/topfreegames/mystack-router/models"
	mystackTest "github.com/topfreegames/mystack-router/testing"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Model", func() {
	Describe("NewRouterConfig", func() {
		It("should construct new RouterConfig", func() {
			routerConfig := models.NewRouterConfig()
			Expect(routerConfig.WorkerProcesses).NotTo(BeEmpty())
			Expect(routerConfig.MaxWorkerConnections).NotTo(BeEmpty())
			Expect(routerConfig.ServerNamesHashMaxSize).NotTo(BeEmpty())
			Expect(routerConfig.ServerNamesHashBucketSize).NotTo(BeEmpty())
		})
	})

	Describe("BuildAppConfig", func() {
		var fakeClientset *fake.Clientset

		BeforeEach(func() {
			fakeClientset = fake.NewSimpleClientset()
		})

		It("should create correct app config", func() {
			_, err := mystackTest.CreateDeployment(fakeClientset)
			Expect(err).NotTo(HaveOccurred())

			service, err := mystackTest.CreateService(fakeClientset)
			Expect(err).NotTo(HaveOccurred())
			Expect(service.Namespace).To(Equal("mystack-user"))
			Expect(service.Name).To(Equal("test"))

			appConfig := models.BuildAppConfig(service, "example.com", "controller.mystack.com", "logger.mystack.com")
			Expect(appConfig.Domain).To(Equal("test.mystack-user.example.com"))
		})

		It("should create correct domain for controller", func() {
			controller, err := mystackTest.CreateController(fakeClientset)
			Expect(err).NotTo(HaveOccurred())
			Expect(controller.Namespace).To(Equal("mystack"))
			Expect(controller.Name).To(Equal("controller"))

			appConfig := models.BuildAppConfig(controller, "example.com", "controller.mystack.com", "logger.mystack.com")
			Expect(appConfig.Domain).To(Equal("controller.mystack.com"))
		})

		It("should create correct domain for logger", func() {
			logger, err := mystackTest.CreateLogger(fakeClientset)
			Expect(err).NotTo(HaveOccurred())
			Expect(logger.Namespace).To(Equal("mystack"))
			Expect(logger.Name).To(Equal("logger"))

			appConfig := models.BuildAppConfig(logger, "example.com", "controller.mystack.com", "logger.mystack.com")
			Expect(appConfig.Domain).To(Equal("logger.mystack.com"))
		})
	})
})
