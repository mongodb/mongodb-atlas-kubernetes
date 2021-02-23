# fix for some Linux distros (i.e. WSL)
SHELL := /usr/bin/env bash

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= 0.3.0

# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "preview,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=preview,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="preview,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= controller-bundle:$(VERSION)

# Image URL to use all building/pushing image targets
IMG ?= controller:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen
# adapted from https://github.com/operator-framework/operator-sdk
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
all: manager ## Build all binaries

.PHONY: int-test
int-test: ENVTEST_ASSETS_DIR = $(shell pwd)/testbin
int-test: export ATLAS_ORG_ID=$(shell grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)
int-test: export ATLAS_PUBLIC_KEY=$(shell grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
int-test: export ATLAS_PRIVATE_KEY=$(shell grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
# magical env that if specified makes the test output 0 on successful runs
# https://github.com/onsi/ginkgo/blob/master/ginkgo/run_command.go#L130
int-test: export GINKGO_EDITOR_INTEGRATION="true"
int-test: generate manifests ## Run integration tests
	mkdir -p $(ENVTEST_ASSETS_DIR)
	test -f $(ENVTEST_ASSETS_DIR)/setup-envtest.sh || curl -sSLo $(ENVTEST_ASSETS_DIR)/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.0/hack/setup-envtest.sh
	source $(ENVTEST_ASSETS_DIR)/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); ginkgo -v -p -nodes=4 ./test/int -coverprofile cover.out

.PHONY: e2e
e2e: run-kind
	./scripts/e2e_local.sh

.PHONY: manager
manager: generate fmt vet ## Build manager binary
	go build -o bin/manager main.go

.PHONY: run
run: generate fmt vet manifests ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

.PHONY: deploy
deploy: run-kind ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	@ act -j deploy -s KUBE_CONFIG_DATA='$(shell kind get kubeconfig | yq r - -j)'

.PHONY: manifests
# Produce CRDs that work back to Kubernetes 1.16 (so 'apiVersion: apiextensions.k8s.io/v1')
manifests: CRD_OPTIONS ?= "crd:crdVersions=v1"
manifests: controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	@./scripts/split_roles_yaml.sh


.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: generate
generate: controller-gen ## Generate code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: controller-gen
# Download controller-gen locally if necessary
CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen:
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)


.PHONY: kustomize
# Download kustomize locally if necessary
KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize:
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v4@v4.0.1)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests kustomize
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

# Additional make goals
.PHONY: run-kind
run-kind: ## Create a local kind cluster
		bash ./scripts/create_kind_cluster.sh;

.PHONY: stop-kind
stop-kind: ## Stop the local kind cluster
	kind delete cluster

.PHONY: log
log: ## View manager logs
	kubectl logs deploy/mongodb-atlas-operator manager -n mongodb-atlas-system -f
