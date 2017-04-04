// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx_test

import (
	"github.com/topfreegames/mystack-router/models"
	. "github.com/topfreegames/mystack-router/nginx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("WriteConfig", func() {
		var routerConfig *models.RouterConfig
		var err error
		var fs models.FileSystem

		BeforeEach(func() {
			fs = models.NewMockFS()
			routerConfig = models.NewRouterConfig()
		})

		It("should write file from RouterConfig", func() {
			err = WriteConfig(routerConfig, fs, "/etc/nginx/nginx.conf")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
