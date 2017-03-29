build:
	@mkdir -p bin && go build -o ./bin/mystack-router main.go

cross-build-linux-amd64:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/mystack-router-linux-amd64
	@chmod a+x ./bin/mystack-router-linux-amd64

build-docker: cross-build-linux-amd64
	@docker build -t mystack-router .

run:
	@go run main.go start --config ./config/local.yaml
