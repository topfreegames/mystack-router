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

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/pkg/labels"
)

//CustomDomains implements CustomDomainsInterface interface
type CustomDomains struct{}

func (c *CustomDomains) GetCustomDomains(controllerDomain, clusterName string) (DomainsPerApp, error) {
	if len(controllerDomain) == 0 || len(clusterName) == 0 {
		return make(DomainsPerApp), nil
	}

	url := fmt.Sprintf("http://%s/cluster-configs/%s/domains", controllerDomain, clusterName)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	fmt.Println(url)

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

func (c *CustomDomains) GetControllerServiceName(clientset kubernetes.Interface) (string, error) {
	labelMap := labels.Set{
		"mystack/controller": "true",
	}
	listOptions := v1.ListOptions{
		LabelSelector: labelMap.AsSelector().String(),
		FieldSelector: fields.Everything().String(),
	}
	services, err := clientset.CoreV1().Services(api.NamespaceAll).List(listOptions)
	if err != nil {
		return "", err
	}

	if len(services.Items) == 0 {
		return "", nil
	}

	name := services.Items[0].GetName()
	port := services.Items[0].Spec.Ports[0].Port
	uri := fmt.Sprintf("%s:%d", name, port)

	return uri, err
}
