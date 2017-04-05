mystack-router
==============
[![Build Status](https://travis-ci.org/topfreegames/mystack-router.svg?branch=master)](https://travis-ci.org/topfreegames/mystack-router)

The router for mystack.

## About
This is the mystack router component, it will discover new services of apps deployed by mystack on Kubernetes cluster and creates routes on [Nginx](http://nginx.org) for your specific domain.
The routes are filtered by namespace (one for each user) and service. 

## Dependencies
* Go 1.7
* Docker

## Building
#### Build a linux binary
```shell
  make cross-build-linux-amd64
```


## Running
This router must run inside Kubernetes cluster. So you need to create a docker image, push it to Dockerhub and run a service using this image. 
Here is an example of how to do it.

#### Build a docker image
On project root, run (mind the dot):
```shell
  docker build -t dockerhub-user/mystack-router:v1 .
```

#### Push it to Dockerhub
```shell
  docker push dockerhub-user/mystack-router:v1
```

#### Create a yaml file for the router
```
kubectl create -f ./manifests
```

#### Configure your domain
If you have the domain `yourdomain.com` registered, you can point `*.yourdomain.com` to your mystack-router loadbalancer external-ip and access your service with:
```shell
curl -v {{appname}}.{{user}}.yourdomain.com
```

#### Access your services
Given that you've pointed `*.yourdomain.com` to the router's LB address, access a service with:
```
example:
app: testapp
user: test-user

curl testapp.test-user.yourdomain.com
```
