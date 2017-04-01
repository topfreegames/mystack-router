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
	WorkerProcesses          string `key:"workerProcesses" constraint:"^(auto|[1-9]\\d*)$"`
	MaxWorkerConnections     string `key:"maxWorkerConnections" constraint:"^[1-9]\\d*$"`
	DefaultTimeout           string `key:"defaultTimeout" constraint:"^[1-9]\\d*(ms|[smhdwMy])?$"`
	ServerNameHashMaxSize    string `key:"serverNameHashMaxSize" constraint:"^[1-9]\\d*[kKmM]?$"`
	ServerNameHashBucketSize string `key:"serverNameHashBucketSize" constraint:"^[1-9]\\d*[kKmM]?$"`
	BodySize                 string `key:"bodySize" constraint:"^[0-9]\\d*[kKmM]?$"`
	LogFormat                string `key:"logFormat"`
	ErrorLogLevel            string `key:"errorLogLevel" constraint:"^(debug|info|notice|warn|error|crit|alert|emerg)$"`
	PlatformDomain           string `key:"platformDomain" constraint:"(?i)^([a-z0-9]+(-[a-z0-9]+)*\\.)+[a-z0-9]+(-*[a-z0-9]+)+$"`
	AppConfigs               []*AppConfig
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
		WorkerProcesses:          "auto",
		MaxWorkerConnections:     "768",
		DefaultTimeout:           "1300s",
		ServerNameHashMaxSize:    "512",
		ServerNameHashBucketSize: "64",
		BodySize:                 "1m",
		LogFormat:                `[$time_iso8601] - $app_name - $remote_addr - $remote_user - $status - "$request" - $bytes_sent - "$http_referer" - "$http_user_agent" - "$server_name" - $upstream_addr - $http_host - $upstream_response_time - $request_time`,
		ErrorLogLevel:            "error",
	}
}

//BuildAppConfig builds AppConfig from kubeclient
func BuildAppConfig(kubeClient *kubernetes.Clientset, service v1.Service, routerConfig *RouterConfig) (*AppConfig, error) {
	appConfig := &AppConfig{}

	appConfig.Name = service.Labels["app"]
	// If we didn't get the app name from the app label, fall back to inferring the app name from
	// the service's own name.
	if appConfig.Name == "" {
		appConfig.Name = service.Name
	}
	// if app name and Namespace are not same then combine the two as it
	// makes deis services (as an example) clearer, such as deis/controller
	if appConfig.Name != service.Namespace {
		appConfig.Name = service.Namespace + "/" + appConfig.Name
	}
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
