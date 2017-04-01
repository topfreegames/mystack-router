// mystack api
// https://github.com/topfreegames/mystack
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
	"github.com/Masterminds/sprig"
	"github.com/topfreegames/mystack-router/model"
	"os"
	"text/template"
)

const (
	confTemplate = `{{ $routerConfig := . }}daemon off;
pid /tmp/nginx.pid;
worker_processes {{ $routerConfig.WorkerProcesses }};
events {
	worker_connections {{ $routerConfig.MaxWorkerConnections }};
	# multi_accept on;
}
http {
	# basic settings
	sendfile on;
	tcp_nopush on;
	tcp_nodelay on;
	# The timeout value must be greater than the front facing load balancers timeout value.
	# Default is the deis recommended timeout value for ELB - 1200 seconds + 100s extra.
	keepalive_timeout {{ $routerConfig.DefaultTimeout }};
	types_hash_max_size 2048;
	server_names_hash_max_size {{ $routerConfig.ServerNameHashMaxSize }};
	server_names_hash_bucket_size {{ $routerConfig.ServerNameHashBucketSize }};
	client_max_body_size {{ $routerConfig.BodySize }};
	real_ip_recursive on;
	log_format upstreaminfo '{{ $routerConfig.LogFormat }}';
	access_log /tmp/logpipe upstreaminfo;
	error_log  /tmp/logpipe {{ $routerConfig.ErrorLogLevel }};
	map $http_upgrade $connection_upgrade {
		default upgrade;
		'' close;
	}
	# Determine the forwarded port:
	# 1. First map the unprivileged ports that Nginx (as a non-root user) actually listen on to the
	# familiar, equivalent privileged ports. (These would be the ports the k8s service listens on.)
	map $server_port $standard_server_port {
		default $server_port;
		8080 80;
		6443 443;
	}
	# 2. If the X-Forwarded-Port header has been set already (e.g. by a load balancer), use its
	# value, otherwise, the port we're forwarding for is the $standard_server_port we determined
	# above.
	map $http_x_forwarded_port $forwarded_port {
		default $http_x_forwarded_port;
		'' $standard_server_port;
	}
	# Default server handles requests for unmapped hostnames, including healthchecks
	server {
		listen 8080 default_server reuseport;
		set $app_name "router-default-vhost";
		server_name _;
		location ~ ^/healthz/?$ {
			access_log off;
			default_type 'text/plain';
			return 200;
		}
		location / {
			return 404;
		}
	}
	# Healthcheck on 9090 -- never uses proxy_protocol
	server {
		listen 9090 default_server;
		server_name _;
		set $app_name "router-healthz";
		location ~ ^/healthz/?$ {
			access_log off;
			default_type 'text/plain';
			return 200;
		}
		location / {
			return 404;
		}
	}
	{{range $appConfig := $routerConfig.AppConfigs}}{{range $domain := $appConfig.Domains}}server {
		listen 8080;
		server_name {{ if contains "." $domain }}{{ $domain }}{{ else if ne $routerConfig.PlatformDomain "" }}{{ $domain }}.{{ $routerConfig.PlatformDomain }}{{ else }}~^{{ $domain }}\.(?<domain>.+)${{ end }};
		server_name_in_redirect off;
		port_in_redirect off;
		set $app_name "{{ $appConfig.Name }}";
		vhost_traffic_status_filter_by_set_key {{ $appConfig.Name }} application::*;
		location / {
			proxy_set_header Host $host;
			proxy_set_header X-Forwarded-For $remote_addr;
			proxy_set_header X-Forwarded-Proto $access_scheme;
			proxy_set_header X-Forwarded-Port $forwarded_port;
			proxy_redirect off;
			proxy_http_version 1.1;
			proxy_set_header Upgrade $http_upgrade;
			proxy_set_header Connection $connection_upgrade;
			proxy_pass http://{{$appConfig.ServiceIP}}:80;{{ else }}return 503;{{ end }}
		}
	}
	{{end}}
}
`
)

//WriteConfig writes a new nginx file config
func WriteConfig(routerConfig *model.RouterConfig, filePath string) error {
	tmpl, err := template.New("nginx").Funcs(sprig.TxtFuncMap()).Parse(confTemplate)
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
