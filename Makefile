SHELL=/bin/bash

# Project variables
BINARY_NAME ?= vmclarity
VERSION ?= $(shell git rev-parse HEAD)
DOCKER_REGISTRY ?= ghcr.io/openclarity
DOCKER_IMAGE ?= $(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_TAG ?= ${VERSION}
VMCLARITY_TOOLS_BASE ?=

# Dependency versions
GOLANGCI_VERSION = 1.52.2
LICENSEI_VERSION = 0.5.0

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: build
build: ui build-all-go ## Build All

.PHONY: build-all-go
build-all-go: bin/vmclarity-apiserver cli ## Build All GO executables

.PHONY: ui
ui: ## Build UI
	@(echo "Building UI ..." )
	@(cd ui; npm i ; npm run build; )
	@ls -l ui/build

bin/vmclarity-apiserver: $(shell find api/models) $(shell find api/server) $(shell find api/client) $(shell find cmd/vmclarity-apiserver) $(shell find pkg) go.mod go.sum ## Build apiserver
	@(echo "Building API server..." )
	go build -race -o bin/vmclarity-apiserver cmd/vmclarity-apiserver/main.go

.PHONY: cli
cli: ## Build CLI
	@(echo "Building CLI ..." )
	go build -race -ldflags="-X 'github.com/openclarity/vmclarity/pkg/cli.GitRevision=${VERSION}'" -o bin/vmclarity-cli cmd/vmclarity-cli/main.go

.PHONY: docker
docker: docker-apiserver docker-cli ## Build All Docker images

.PHONY: push-docker
push-docker: push-docker-apiserver push-docker-cli ## Build and Push All Docker images

ifneq ($(strip $(VMCLARITY_TOOLS_BASE)),)
VMCLARITY_TOOLS_CLI_DOCKER_ARG=--build-arg VMCLARITY_TOOLS_BASE=${VMCLARITY_TOOLS_BASE}
endif

.PHONY: docker-cli
docker-cli: ## Build CLI Docker image
	@(echo "Building cli docker image ..." )
	docker build \
		--file ./Dockerfile.cli \
		--build-arg COMMIT_HASH=$(shell git rev-parse HEAD) \
		${VMCLARITY_TOOLS_CLI_DOCKER_ARG} \
		-t ${DOCKER_IMAGE}-cli:${DOCKER_TAG} \
		.

.PHONY: push-docker-cli
push-docker-cli: docker-cli ## Build and Push CLI Docker image
	@echo "Publishing cli docker image ..."
	docker push $(DOCKER_IMAGE)-cli:$(DOCKER_TAG)

.PHONY: docker-apiserver
docker-apiserver: ## Build Backend Docker image
	@(echo "Building API server docker image ..." )
	docker build --file ./Dockerfile.apiserver --build-arg VERSION=${VERSION} \
		--build-arg BUILD_TIMESTAMP=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg COMMIT_HASH=$(shell git rev-parse HEAD) \
		-t ${DOCKER_IMAGE}-apiserver:${DOCKER_TAG} .

.PHONY: push-docker-apiserver
push-docker-apiserver: docker-apiserver ## Build and Push Backend Docker image
	@echo "Publishing API server docker image ..."
	docker push ${DOCKER_IMAGE}-apiserver:${DOCKER_TAG}

.PHONY: test
test: ## Run Unit Tests
	@go test ./...

.PHONY: clean-ui
clean-ui:
	@(rm -rf ui/build ; echo "UI cleanup done" )

.PHONY: clean
clean: clean-ui clean-backend ## Clean all build artifacts

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint-go
lint-go: bin/golangci-lint
	./bin/golangci-lint run

.PHONY: lint-cfn
lint-cfn:
	# Requires cfn-lint to be installed
	# https://github.com/aws-cloudformation/cfn-lint#install
	cfn-lint installation/aws/VmClarity.cfn

.PHONY: lint
lint: lint-go lint-cfn ## Run linters

.PHONY: fix
fix: bin/golangci-lint ## Fix lint violations
	./bin/golangci-lint run --fix

bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@

.PHONY: license-check
license-check: bin/licensei ## Run license check
	./bin/licensei header

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	./bin/licensei cache --config=../.licensei.toml

.PHONY: check
check: lint test ## Run tests and linters

.PHONY: gomod-tidy
gomod-tidy:
	go mod tidy

.PHONY: api
api: api-backend api-ui ## Generating API code

.PHONY: api-backend
api-backend: ## Generating API for backend code
	@(echo "Generating API for backend code ..." )
	@(cd api; go generate)

.PHONY: api-ui
api-ui: ## Generating API for UI backend code
	@(echo "Generating API for UI backend code ..." )
	@(cd pkg/uibackend/api; go generate)
