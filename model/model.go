// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package model

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// RouterConfig is the primary type used to encapsulate all router configuration.
type RouterConfig struct {
	WorkerProcesses           string `key:"workerProcesses" constraint:"^(auto|[1-9]\\d*)$"`
	MaxWorkerConnections      string `key:"maxWorkerConnections" constraint:"^[1-9]\\d*$"`
	ServerNamesHashMaxSize    string `key:"serverNamesHashMaxSize" constraint:"^[1-9]\\d*$"`
	ServerNamesHashBucketSize string `key:"serverNamesHashBucketSize" constraint:"^[1-9]\\d*$"`
	AppConfigs                []*AppConfig
}

// AppConfig encapsulates the configuration for all routes to a single back end.
type AppConfig struct {
	Domain    string
	ServiceIP string
	Available bool
}

//NewRouterConfig builds new router config with default values
//If this error occurs "could not build the server_names_hash", either increase ServerNamesHashMaxSize
//to a number close to the number of servers (users * services he/she uses), or increase ServerNamesHashBucketSize (in this case, the server_name is becoming too long)
func NewRouterConfig() *RouterConfig {
	return &RouterConfig{
		WorkerProcesses:           "1",
		MaxWorkerConnections:      "1024",
		ServerNamesHashMaxSize:    "512",
		ServerNamesHashBucketSize: "128",
	}
}

//BuildAppConfig builds AppConfig from kubeclient
func BuildAppConfig(kubeClient *kubernetes.Clientset, service v1.Service, routerConfig *RouterConfig) (*AppConfig, error) {
	appConfig := &AppConfig{}

	appConfig.Domain = fmt.Sprintf("%s.%s.tfgapps.com", service.Name, service.Namespace)

	appConfig.ServiceIP = service.Spec.ClusterIP
	endpointsClient := kubeClient.Endpoints(service.Namespace)
	endpoints, err := endpointsClient.Get(service.Name)
	if err != nil {
		return nil, err
	}
	appConfig.Available = len(endpoints.Subsets) > 0 && len(endpoints.Subsets[0].Addresses) > 0
	return appConfig, nil
}
