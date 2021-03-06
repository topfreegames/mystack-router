// mystack api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/topfreegames/mystack-router/models"
)

const configTemplate = `
worker_processes {{.WorkerProcesses}};
events {
  worker_connections {{.MaxWorkerConnections}};
}
stream {
	server {
		listen 28000;
		proxy_pass mystack-controller:28000;
	}
	
  {{- range .AppConfigs}}
	{{- $name := .AppName}}
	{{- $namespace := .AppNamespace}}
	{{- $ports := .Ports}}
	{{- range $i, $s := .SocketPorts}}
	server {
		listen {{$s}}
		proxy_pass {{$name}}.{{$namespace}}.svc.cluster.local.:{{index $ports $i}};
	}
	{{- end}}
	{{- end}}
}
http {
  server_names_hash_bucket_size {{.ServerNamesHashBucketSize}};
  server_names_hash_max_size {{.ServerNamesHashMaxSize}};
  keepalive_timeout 1300s;
  {{- $controllerDomain := .ControllerDomain}}
  {{- $loggerDomain := .LoggerDomain}}
  {{- range .AppConfigs}}
	{{- $name := .AppName}}
	{{- $namespace := .AppNamespace}}
	{{- $domain := .Domain}}
	{{- $isSocket := .IsSocket}}
	{{- range .Ports}}
  {{- if eq $domain $controllerDomain}}
  server {
    listen 80;
    server_name login;
    location / {
      proxy_connect_timeout 60s;
      proxy_send_timeout 1300s;
      proxy_read_timeout 1300s;
      proxy_pass http://{{$name}}.{{$namespace}}:{{.}};
    }
  }
  server {
    listen 80;
    server_name {{$domain}};
    location / {
      proxy_connect_timeout 60s;
      proxy_send_timeout 1300s;
      proxy_read_timeout 1300s;
      proxy_pass http://{{$name}}.{{$namespace}}:{{.}};
    }
  }
  {{- else if eq $domain $loggerDomain}}
  server {
    listen 80;
    server_name {{$domain}};
    location / {
      proxy_connect_timeout 60s;
      proxy_send_timeout 1300s;
      proxy_read_timeout 1300s;
      proxy_pass http://{{$name}}.{{$namespace}}:{{.}};
    }
  }
  {{- else if not $isSocket}}
  server {
    listen 80;
    server_name {{$domain}};
    location / {
      proxy_pass http://{{$name}}.{{$namespace}}:{{.}};
    }
  }
  {{- end}}
	{{- end}}
	{{- end}}  
  server {
    listen 80 default_server;
    server_name _;
    return 404;
  }
}
`

//WriteConfig writes a new nginx file config
func WriteConfig(routerConfig *models.RouterConfig, fs models.FileSystem, filePath string) error {
	tmpl, err := template.New("nginx").Funcs(sprig.TxtFuncMap()).Parse(configTemplate)
	if err != nil {
		return err
	}

	file, err := fs.Create(filePath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, routerConfig)

	return err
}
