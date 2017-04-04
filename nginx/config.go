// mystack api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
	"github.com/Masterminds/sprig"
	"github.com/topfreegames/mystack-router/models"
	"os"
	"text/template"
)

const configTemplate = `
worker_processes {{.WorkerProcesses}};

events {
	worker_connections {{.MaxWorkerConnections}};
}

http {
	server_names_hash_bucket_size {{.ServerNamesHashBucketSize}};
	server_names_hash_max_size {{.ServerNamesHashMaxSize}};

	{{range .AppConfigs}}
	server {
		listen 8080;
		server_name {{.Domain}};

		location / {
			proxy_pass http://{{.ServiceIP}}:80;
		}
	}
	{{end}}	

  server {
    listen 8080 default_server;
    server_name _;
    return 404;
  }
}
`

//WriteConfig writes a new nginx file config
func WriteConfig(routerConfig *models.RouterConfig, filePath string) error {
	tmpl, err := template.New("nginx").Funcs(sprig.TxtFuncMap()).Parse(configTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(os.Stdout, routerConfig)
	err = tmpl.Execute(file, routerConfig)

	return err
}
