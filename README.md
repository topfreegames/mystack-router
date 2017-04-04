mystack-router
==============

The router for mystack.

## About
Discovers new services on Kubernetes cluster and creates routes on [Nginx](http://nginx.org) for your specific domain.

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
Supposing the following router.yaml file
```yaml
apiVersion: v1
kind: Service
metadata:
  name: mystack-router
spec:
  selector:
    app: mystack-router
  ports:
    - protocol: TCP
      port: 8080
  type: LoadBalancer
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: mystack-router
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: mystack-router
    spec:
      containers:
        - name: mystack-router
          image: dockerhub-user/mystack-router:v1
          ports:
            - containerPort: 8080
```
Run it with:
```shell
  kubectl create -f router.yaml
```

#### Access your services
Now, if there are services running on Kubernetes as ClusterIP, they are accessable through mystack-router.
For example, given that:
* There is a service running on namespace `mystack-user`
* The service name is `hello-world`
* The labels `mystack/routable: "true"` and `mystack/owner: user` are defined
* The domain sufix is `example.com` (defined on config/local.yaml)
* The Kubernetes IP is `k8s_ip`

Then this service is reachable through mystack-router with:
```shell
  curl -v -H 'Host: hello-world.mystack-user.example.com' http://k8s_ip:8080
```

If you have a domain with prefix `example.com` on the internet, you can point `*.example.com` to your mystack-router external-ip and access your service with:
```shell
  curl -v hello-world.mystack-user.example.com:8080
```
