// mystack
// +build integration
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package testing

import (
	"bytes"
	"html/template"

	"github.com/topfreegames/maestro/errors"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
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
	// ServiceYaml ...
	ServiceYaml = `
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
    app: controller
  ports:
    - protocol: TCP
      port: 80
      targetPort: 5000
  type: ClusterIP
`
	controllerDeployYaml = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: mystack-controller
  namespace: mystack
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: mystack-controller
    spec:
      containers:
        - name: mystack-controller
          image: hello-world
          ports:
            - containerPort: 8080
`
	controllerServiceYaml = `
apiVersion: v1
kind: Service
metadata:
  name: mystack-controller
  namespace: mystack
  labels:
    mystack/routable: "true"
spec:
  selector:
    app: test
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
`
	loggerDeployYaml = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: mystack-logger
  namespace: mystack
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: mystack-logger
    spec:
      containers:
        - name: mystack-logger
          image: hello-world
          ports:
            - containerPort: 8080
`
	loggerServiceYaml = `
apiVersion: v1
kind: Service
metadata:
  name: mystack-logger
  namespace: mystack
  labels:
    mystack/routable: "true"
spec:
  selector:
    app: test
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
`
	namespace = "mystack-user"
)

//CreateService creates a mock service on kubernetes for testing purposes
func CreateService(clientset kubernetes.Interface) (*v1.Service, error) {
	return createService(clientset, ServiceYaml, namespace)
}

func createService(clientset kubernetes.Interface, yaml, namespace string) (*v1.Service, error) {
	tmpl, err := template.New("create").Parse(serviceYaml)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, s)
	if err != nil {
		return nil, err
	}

	k8sService := v1.Service{}
	err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader(buf.Bytes()), len(buf.Bytes())).Decode(&k8sService)
	if err != nil {
		return nil, errors.NewKubernetesError("error unmarshaling pod", err)
	}

	service, err := clientset.CoreV1().Services(s.Namespace).Create(&k8sService)
	if err != nil {
		return nil, errors.NewKubernetesError("create pod error", err)
	}

	return service, nil
}

//CreateDeployment creates a mock deployment on kubernetes for testing purposes
func CreateDeployment(clientset kubernetes.Interface) (*v1beta1.Deployment, error) {
	return createDeployment(clientset, deployYaml, namespace)
}

func createDeployment(clientset kubernetes.Interface, yaml, namespace string) (*v1beta1.Deployment, error) {
	d := api.Codecs.UniversalDecoder()
	obj, _, err := d.Decode([]byte(yaml), nil, nil)
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

//CreateController creates a mock service on kubernetes for testing purposes
func CreateController(clientset kubernetes.Interface) (*v1.Service, error) {
	namespaceStr := "mystack"
	_, err := createDeployment(clientset, controllerDeployYaml, namespaceStr)
	if err != nil {
		return nil, err
	}

	service, err := createService(clientset, controllerServiceYaml, namespaceStr)
	return service, err
}

//CreateLogger creates a mock service on kubernetes for testing purposes
func CreateLogger(clientset kubernetes.Interface) (*v1.Service, error) {
	namespaceStr := "mystack"
	_, err := createDeployment(clientset, loggerDeployYaml, namespaceStr)
	if err != nil {
		return nil, err
	}

	service, err := createService(clientset, loggerServiceYaml, namespaceStr)
	return service, err
}
