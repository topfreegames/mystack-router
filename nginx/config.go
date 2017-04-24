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
http {
	server_names_hash_bucket_size {{.ServerNamesHashBucketSize}};
	server_names_hash_max_size {{.ServerNamesHashMaxSize}};
	{{range .AppConfigs}}{{$name := .AppName}}{{$namespace := .AppNamespace}}{{$domain := .Domain}}{{$customDomains := .CustomDomains}}{{range .Ports}}
	server {
		listen 80;
		server_name {{$domain}};
		{{range $domain := $customDomains}}server_name {{$domain}}{{end}};
		location / {
			proxy_pass http://{{$name}}.{{$namespace}}:{{.}};
		}
	}
	{{end}}{{end}}	
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
