# fix for some Linux distros (i.e. WSL)
SHELL := /usr/bin/env bash

# CONTAINER ENGINE: docker | podman
CONTAINER_ENGINE?=docker

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= $(shell git describe --tags --dirty --broken | cut -c 2-)

# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "preview,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=preview,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="preview,fast,stable")
CHANNELS ?= beta
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# Used by the olm-deploy if you running on Mac and deploy to K8S/Openshift
ifndef TARGET_ARCH
TARGET_ARCH := $(shell go env GOARCH)
endif
ifndef TARGET_OS
TARGET_OS := $(shell go env GOOS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
DEFAULT_CHANNEL ?= beta
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Base registry for the operator, bundle, catalog images
REGISTRY ?= quay.io/mongodb
# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= $(REGISTRY)/mongodb-atlas-kubernetes-operator-prerelease-bundle:$(VERSION)

#IMG ?= mongodb-atlas-kubernetes-operator:latest
#BUNDLE_REGISTRY ?= $(REGISTRY)/mongodb-atlas-operator-bundle
OPERATOR_REGISTRY ?= $(REGISTRY)/mongodb-atlas-kubernetes-operator-prerelease
CATALOG_REGISTRY ?= $(REGISTRY)/mongodb-atlas-kubernetes-operator-prerelease-catalog
OPERATOR_IMAGE ?= ${OPERATOR_REGISTRY}:${VERSION}
CATALOG_IMAGE ?= ${CATALOG_REGISTRY}:${VERSION}
TARGET_NAMESPACE ?= mongodb-atlas-operator-system-test

# Image URL to use all building/pushing image targets
ifndef IMG
IMG := ${OPERATOR_REGISTRY}:${VERSION}
endif

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

.PHONY: unit-test
unit-test:
	go test -race -cover ./pkg/...

.PHONY: int-test
int-test: ENVTEST_ASSETS_DIR = $(shell pwd)/testbin
int-test: export ATLAS_ORG_ID=$(shell grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)
int-test: export ATLAS_PUBLIC_KEY=$(shell grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
int-test: export ATLAS_PRIVATE_KEY=$(shell grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
# magical env that if specified makes the test output 0 on successful runs
# https://github.com/onsi/ginkgo/blob/master/ginkgo/run_command.go#L130
int-test: export GINKGO_EDITOR_INTEGRATION="true"
int-test: generate manifests ## Run integration tests. Sample with labels: `make int-test label=AtlasProject` or `make int-test label='AtlasDeployment && !slow'`
	mkdir -p $(ENVTEST_ASSETS_DIR)
	test -f $(ENVTEST_ASSETS_DIR)/setup-envtest.sh || curl -sSLo $(ENVTEST_ASSETS_DIR)/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.0/hack/setup-envtest.sh
	source $(ENVTEST_ASSETS_DIR)/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); ./scripts/int_local.sh $(label)

.PHONY: e2e
e2e: run-kind ## Run e2e test. Command `make e2e label=cluster-ns` run cluster-ns test
	./scripts/e2e_local.sh $(label) $(build)

.PHONY: e2e-openshift-upgrade
e2e-openshift-upgrade:
	cd scripts && ./openshift-upgrade-test.sh

.PHONY: manager
manager: generate fmt vet ## Build manager binary
	@echo "Building operator with version $(VERSION)"
	 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o bin/manager -ldflags="-X github.com/mongodb/mongodb-atlas-kubernetes/pkg/version.Version=$(VERSION)" cmd/manager/main.go

.PHONY: run
run: generate fmt vet manifests ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./cmd/manager/main.go

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

.PHONY: deploy
deploy: generate manifests run-kind ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	@./scripts/deploy.sh

.PHONY: manifests
# Produce CRDs that work back to Kubernetes 1.16 (so 'apiVersion: apiextensions.k8s.io/v1')
manifests: CRD_OPTIONS ?= "crd:crdVersions=v1"
manifests: controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	@./scripts/split_roles_yaml.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...
	find . -name "*.go" -not -path "./vendor/*" -exec gofmt -w "{}" \;
	find . -name "*.go" -not -path "./vendor/*" -exec goimports -l -w "{}" \;

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: generate
generate: controller-gen ## Generate code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: controller-gen
CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.2)


.PHONY: kustomize
KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary
ifeq ("$(wildcard $(KUSTOMIZE))", "")
	rm -f ./kustomize
	wget "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" -O kinstall.sh
	chmod +x ./kinstall.sh && bash -c ./kinstall.sh && mv ./kustomize $(GOBIN)/kustomize
	rm -f ./kinstall.sh
endif

# go-get-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

.PHONY: bundle
bundle: manifests kustomize ## Generate bundle manifests and metadata, then validate generated files.
	@echo "Building bundle $(VERSION)"
	operator-sdk generate kustomize manifests -q --apis-dir=pkg/api
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

.PHONY: image
image: ## Build the operator image
	$(CONTAINER_ENGINE) build --build-arg VERSION=$(VERSION) -t $(OPERATOR_IMAGE) .
	$(CONTAINER_ENGINE) push $(OPERATOR_IMAGE)

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	$(CONTAINER_ENGINE) build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: bundle-push
bundle-push:
	$(CONTAINER_ENGINE) push $(BUNDLE_IMG)

.PHONY: catalog-build
CATALOG_DIR ?= ./scripts/openshift/atlas-catalog
catalog-build: image
	CATALOG_DIR=$(CATALOG_DIR) \
	CHANNEL=$(DEFAULT_CHANNEL) \
	CATALOG_IMAGE=$(CATALOG_IMAGE) \
	BUNDLE_IMAGE=$(BUNDLE_IMG) \
	VERSION=$(VERSION) \
	./scripts/build_catalog.sh

.PHONY: catalog-push
catalog-push:
	$(CONTAINER_ENGINE) push $(CATALOG_IMAGE)

.PHONY: build-subscription
build-subscription:
	CATALOG_DIR=$(shell dirname "$(CATALOG_DIR)") \
	CHANNEL=$(DEFAULT_CHANNEL) \
	TARGET_NAMESPACE=$(TARGET_NAMESPACE) \
	./scripts/build_subscription.sh

.PHONY: build-catalogsource
build-catalogsource:
	CATALOG_DIR=$(shell dirname "$(CATALOG_DIR)") \
	CATALOG_IMAGE=$(CATALOG_IMAGE) \
	CATALOG_DISPLAY_NAME="MongoDB Atlas operator $(VERSION)" \
	./scripts/build_catalogsource.sh

.PHONY: deploy-olm
# Deploy atlas operator to the running openshift cluster with OLM
deploy-olm: bundle-build bundle-push catalog-build catalog-push build-catalogsource build-subscription
	oc -n openshift-marketplace delete catalogsource mongodb-atlas-kubernetes-local --ignore-not-found
	oc delete namespace $(TARGET_NAMESPACE) --ignore-not-found
	oc create namespace $(TARGET_NAMESPACE)
	oc -n openshift-marketplace apply -f ./scripts/openshift/catalogsource.yaml
	oc -n $(TARGET_NAMESPACE) apply -f ./scripts/openshift/operatorgroup.yaml
	oc -n $(TARGET_NAMESPACE) apply -f ./scripts/openshift/subscription.yaml

## Disabled for now
## .PHONY: docker-login-olm
## docker-login-olm:
## docker login -u $(shell oc whoami) -p $(shell oc whoami -t) $(REGISTRY)

.PHONY: image-push
image-push: ## Push the docker image
	$(CONTAINER_ENGINE) push ${IMG}

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

.PHONY: clear-atlas
clear-atlas: export INPUT_ATLAS_PUBLIC_KEY=$(shell grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
clear-atlas: export INPUT_ATLAS_PRIVATE_KEY=$(shell grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
clear-atlas: ## Clear Atlas organization
	bash .github/actions/cleanup/entrypoint.sh

.PHONY: post-install-hook
post-install-hook:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o bin/helm-post-install cmd/post-install/main.go
	chmod +x bin/helm-post-install

.PHONY: x509-cert
x509-cert: ## Create X.509 cert at path tmp/x509/ (see docs/x509-user.md)
	go run scripts/create_x509.go
