// mystack-router
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//CustomDomains implements CustomDomainsInterface interface
type CustomDomains struct{}

func (*CustomDomains) GetCustomDomains(controllerDomain, clusterName string) (DomainsPerApp, error) {
	url := fmt.Sprintf("%s/domains/%s", controllerDomain, clusterName)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d when requesting custom domains at %s", res.StatusCode, clusterName)
	}

	customDomains := make(DomainsPerApp)
	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bts, &customDomains)
	if err != nil {
		return nil, err
	}

	return customDomains, nil
}
