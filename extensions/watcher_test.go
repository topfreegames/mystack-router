package extensions_test

import (
	"github.com/spf13/afero"
	. "github.com/topfreegames/mystack-router/extensions"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		fakeClientset := fake.NewSimpleClientset()
		Expect(config).To(Equal("qwe"))
		_, err := NewWatcher(config, fakeClientset)

		Expect(err).NotTo(HaveOccurred())
	})
})
