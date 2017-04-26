// mystack-router
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
)

// RouterConfig is the primary type used to encapsulate all router configuration.
type RouterConfig struct {
	WorkerProcesses           string `key:"workerProcesses" constraint:"^(auto|[1-9]\\d*)$"`
	MaxWorkerConnections      string `key:"maxWorkerConnections" constraint:"^[1-9]\\d*$"`
	ServerNamesHashMaxSize    string `key:"serverNamesHashMaxSize" constraint:"^[1-9]\\d*$"`
	ServerNamesHashBucketSize string `key:"serverNamesHashBucketSize" constraint:"^[1-9]\\d*$"`
	AppConfigs                []*AppConfig
	ControllerDomain          string
}

// AppConfig encapsulates the configuration for all routes to a single back end.
type AppConfig struct {
	Domain       string
	AppName      string
	AppNamespace string
	ClusterName  string
	Ports        []int
}

//NewRouterConfig builds new router config with default values
//If this error occurs "could not build the server_names_hash", either increase ServerNamesHashMaxSize
//to a number close to the number of servers (users * services he/she uses), or increase ServerNamesHashBucketSize (in this case, the server_name is becoming too long)
func NewRouterConfig(kubeDomainSuffix string) *RouterConfig {
	return &RouterConfig{
		WorkerProcesses:           "1",
		MaxWorkerConnections:      "1024",
		ServerNamesHashMaxSize:    "512",
		ServerNamesHashBucketSize: "128",
		ControllerDomain:          fmt.Sprintf("controller.%s", kubeDomainSuffix),
	}
}

//BuildAppConfig builds AppConfig from service
func BuildAppConfig(
	service *v1.Service,
	kubeDomainSuffix string,
) *AppConfig {
	appConfig := &AppConfig{}

	switch service.GetName() {
	case "mystack-controller":
		appConfig.Domain = fmt.Sprintf("controller.%s", kubeDomainSuffix)
		break
	case "mystack-logger":
		appConfig.Domain = fmt.Sprintf("logger.%s", kubeDomainSuffix)
		break
	default:
		appConfig.Domain = fmt.Sprintf("%s.%s.%s", service.Name, service.Namespace, kubeDomainSuffix)
	}

	appConfig.AppName = service.ObjectMeta.GetName()
	appConfig.AppNamespace = service.ObjectMeta.GetNamespace()
	appConfig.ClusterName = service.ObjectMeta.Labels["mystack/cluster"]

	ports := make([]int, len(service.Spec.Ports))
	for i, port := range service.Spec.Ports {
		ports[i] = int(port.Port)
	}
	appConfig.Ports = ports

	return appConfig
}
