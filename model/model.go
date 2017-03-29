// mystack-router api
// https://github.com/topfreegames/mystack/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
//"k8s.io/client-go/kubernetes"
)

// RouterConfig is the primary type used to encapsulate all router configuration.
type RouterConfig struct {
	WorkerProcesses      string `key:"workerProcesses" constraint:"^(auto|[1-9]\\d*)$"`
	MaxWorkerConnections string `key:"maxWorkerConnections" constraint:"^[1-9]\\d*$"`
	DefaultTimeout       string `key:"defaultTimeout" constraint:"^[1-9]\\d*(ms|[smhdwMy])?$"`
	ErrorLogLevel        string `key:"errorLogLevel" constraint:"^(debug|info|notice|warn|error|crit|alert|emerg)$"`
}
