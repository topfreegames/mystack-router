// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package testing

import (
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	deployYaml = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: test
  namespace: mystack-user
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
        - name: hello-world
          image: hello-world
          ports:
            - containerPort: 5000
`
	serviceYaml = `
apiVersion: v1
kind: Service
metadata:
  name: test
  namespace: mystack-user
  labels:
    mystack/routable: "true"
    mystack/owner: user
    app: test
spec:
  selector:
    app: test
  ports:
    - protocol: TCP
      port: 80
      targetPort: 5000
  type: ClusterIP
`
	namespace = "mystack-user"
)

//CreateService creates a mock service on kubernetes for testing purposes
func CreateService(clientset *fake.Clientset) (*v1.Service, error) {
	d := api.Codecs.UniversalDecoder()
	obj, _, err := d.Decode([]byte(serviceYaml), nil, nil)
	if err != nil {
		return nil, err
	}

	src := obj.(*api.Service)
	dst := &v1.Service{}

	err = api.Scheme.Convert(src, dst, 0)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Services(namespace).Create(dst)
}

//CreateDeployment creates a mock deployment on kubernetes for testing purposes
func CreateDeployment(clientset *fake.Clientset) (*v1beta1.Deployment, error) {
	d := api.Codecs.UniversalDecoder()
	obj, _, err := d.Decode([]byte(deployYaml), nil, nil)
	if err != nil {
		return nil, err
	}

	src := obj.(*extensions.Deployment)
	dst := &v1beta1.Deployment{}

	err = api.Scheme.Convert(src, dst, 0)
	if err != nil {
		return nil, err
	}

	return clientset.ExtensionsV1beta1().Deployments(namespace).Create(dst)
}
