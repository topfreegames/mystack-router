build:
	@mkdir -p bin && go build -o ./bin/mystack-router main.go

build-docker: cross-build-linux-amd64
	@docker build -t mystack-router .

cross-build-linux-amd64:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/mystack-router-linux-amd64
	@chmod a+x ./bin/mystack-router-linux-amd64

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

gather-integration-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-integration.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-integration.out; done'

integration int: clear-coverage-profiles integration-run gather-integration-profiles

integration-run:
	@ginkgo -tags integration -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

merge-profiles:
	@mkdir -p _build
	@gocovmerge _build/*.out > _build/coverage-all.out

run-dev:
	@env MYSTACK_KUBERNETES_IN_CLUSTER=false MYSTACK_NGINX_DIR=nginx go run main.go start --config ./config/local.yaml

setup-ci:
	@go get -u github.com/golang/dep/...
	@go get github.com/onsi/ginkgo/ginkgo
	@go get github.com/wadey/gocovmerge
	@dep init
	@dep ensure

test: unit integration test-coverage-func

test-coverage-func coverage-func: merge-profiles
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions NOT COVERED by Tests"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@go tool cover -func=_build/coverage-all.out | egrep -v "100.0[%]"

unit: clear-coverage-profiles unit-run gather-unit-profiles

unit-run:
	@ginkgo -tags unit -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}
