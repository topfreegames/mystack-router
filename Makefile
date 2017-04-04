build:
	@mkdir -p bin && go build -o ./bin/mystack-router main.go

cross-build-linux-amd64:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/mystack-router-linux-amd64
	@chmod a+x ./bin/mystack-router-linux-amd64

build-docker: cross-build-linux-amd64
	@docker build -t mystack-router .

run:
	@go run main.go start --config ./config/local.yaml

test: unit integration test-coverage-func

unit: clear-coverage-profiles unit-run gather-unit-profiles

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

unit-run:
	@ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

integration int: clear-coverage-profiles integration-run gather-integration-profiles

integration-run:
	@ginkgo -tags integration -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-integration-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-integration.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-integration.out; done'
