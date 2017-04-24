// mystack-router
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//MockCustomDomains implements CustomDomainsInterface interface
type MockCustomDomains struct {
	CustomDomains DomainsPerApp
	Err           error
}

func (m *MockCustomDomains) GetCustomDomains(controllerDomain, clusterName string) (DomainsPerApp, error) {
	return m.CustomDomains, m.Err
}
