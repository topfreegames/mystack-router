// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package model

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// RouterConfig is the primary type used to encapsulate all router configuration.
type RouterConfig struct {
	WorkerProcesses      string `key:"workerProcesses" constraint:"^(auto|[1-9]\\d*)$"`
	MaxWorkerConnections string `key:"maxWorkerConnections" constraint:"^[1-9]\\d*$"`
	AppConfigs           []*AppConfig
}

// AppConfig encapsulates the configuration for all routes to a single back end.
type AppConfig struct {
	Name      string
	Domains   []string `key:"domains" constraint:"(?i)^((([a-z0-9]+(-*[a-z0-9]+)*)|((\\*\\.)?[a-z0-9]+(-*[a-z0-9]+)*\\.)+[a-z0-9]+(-*[a-z0-9]+)+)(\\s*,\\s*)?)+$"`
	ServiceIP string
	Available bool
}

//NewRouterConfig builds new router config with default values
func NewRouterConfig() *RouterConfig {
	return &RouterConfig{
		WorkerProcesses:      "1",
		MaxWorkerConnections: "1024",
	}
}

//BuildAppConfig builds AppConfig from kubeclient
func BuildAppConfig(kubeClient *kubernetes.Clientset, service v1.Service, routerConfig *RouterConfig) (*AppConfig, error) {
	appConfig := &AppConfig{}

	appConfig.Name = service.Labels["app"]
	if appConfig.Name == "" {
		appConfig.Name = service.Name
	}
	if appConfig.Name != service.Namespace {
		appConfig.Name = service.Namespace + "/" + appConfig.Name
	}

	domain := service.Annotations["router.mystack/domains"]
	appConfig.Domains = []string{domain}

	// If no domains are found, we don't have the information we need to build routes
	// to this application.  Abort.
	if len(appConfig.Domains) == 0 {
		return nil, nil
	}

	appConfig.ServiceIP = service.Spec.ClusterIP
	endpointsClient := kubeClient.Endpoints(service.Namespace)
	endpoints, err := endpointsClient.Get(service.Name)
	if err != nil {
		return nil, err
	}
	appConfig.Available = len(endpoints.Subsets) > 0 && len(endpoints.Subsets[0].Addresses) > 0
	return appConfig, nil
}
