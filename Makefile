# fix for some Linux distros (i.e. WSL)
SHELL := /usr/bin/env bash

# CONTAINER ENGINE: docker | podman
CONTAINER_ENGINE?=docker

DOCKER_SBOM_PLUGIN_VERSION=0.6.1

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= $(shell git describe --tags --dirty --broken | cut -c 2-)

MAJOR_VERSION = $(shell cat major-version)

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

GO_UNIT_TEST_FOLDERS=$(shell go list ./... |grep -v 'test/int\|test/e2e')

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

#BUNDLE_REGISTRY ?= $(REGISTRY)/mongodb-atlas-operator-bundle
RELEASED_OPERATOR_IMAGE ?= mongodb/mongodb-atlas-kubernetes-operator
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

# Track changes to sources to avoid repeating operations on sourced when soruces did not change
TMPDIR ?= /tmp
TIMESTAMPS_DIR := $(TMPDIR)/mongodb-atlas-kubernetes
GO_SOURCES = $(shell find . -type f -name '*.go' -not -path './vendor/*')

# Defaults for make run
OPERATOR_POD_NAME = mongodb-atlas-operator
OPERATOR_NAMESPACE = mongodb-atlas-system
ATLAS_DOMAIN = https://cloud-qa.mongodb.com/
ATLAS_KEY_SECRET_NAME = mongodb-atlas-operator-api-key

BASE_GO_PACKAGE = github.com/mongodb/mongodb-atlas-kubernetes/v2
GO_LICENSES = go-licenses
DISALLOWED_LICENSES = restricted,reciprocal

# golangci-lint
GOLANGCI_LINT_VERSION := v1.54.2

REPORT_TYPE = flakiness
SLACK_WEBHOOK ?= https://hooks.slack.com/services/...

# Signature definitions
SIGNATURE_REPO ?= OPERATOR_REGISTRY
AKO_SIGN_PUBKEY = https://cosign.mongodb.com/atlas-kubernetes-operator.pem

# Licenses status
GOMOD_SHA := $(shell git ls-files -s go.mod | awk '{print $$1" "$$2" "$$4}')
LICENSES_GOMOD_SHA_FILE := .licenses-gomod.sha256
GOMOD_LICENSES_SHA := $(shell cat $(LICENSES_GOMOD_SHA_FILE))

OPERATOR_NAMESPACE=atlas-operator
OPERATOR_POD_NAME=mongodb-atlas-operator
RUN_YAML= # Set to the YAML to run when calling make run

LOCAL_IMAGE=mongodb-atlas-kubernetes-operator:compiled
CONTAINER_SPEC=.spec.template.spec.containers[0]

SILK_ASSET_GROUP="atlas-kubernetes-operator"

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

go-licenses:
	go install github.com/google/go-licenses@latest

licenses.csv: go-licenses go.mod ## Track licenses in a CSV file
	@echo "Tracking licenses into file $@"
	@echo "========================================"
	GOOS=linux GOARCH=amd64 go mod download
	# https://github.com/google/go-licenses/issues/244
	GOTOOLCHAIN=local GOOS=linux GOARCH=amd64 $(GO_LICENSES) csv --include_tests $(BASE_GO_PACKAGE)/... > $@
	echo $(GOMOD_SHA) > $(LICENSES_GOMOD_SHA_FILE)

recompute-licenses: ## Recompute the licenses.csv only if needed (gomod was changed)
	@[ "$(GOMOD_SHA)" == "$(GOMOD_LICENSES_SHA)" ] || $(MAKE) licenses.csv

licenses-up-to-date:
	@if [ "$(GOMOD_SHA)" != "$(GOMOD_LICENSES_SHA)" ]; then \
	echo "licenses.csv needs to be recalculated: git rebase AND run 'make licenses.csv'"; exit 1; \
	else echo "licenses.csv is OK! (up to date)"; fi

.PHONY: check-licenses
check-licenses: go-licenses licenses-up-to-date ## Check licenses are compliant with our restrictions
	@echo "Checking licenses not to be: $(DISALLOWED_LICENSES)"
	@echo "============================================"
	# https://github.com/google/go-licenses/issues/244
	GOTOOLCHAIN=local GOOS=linux GOARCH=amd64 $(GO_LICENSES) check --include_tests \
	--disallowed_types $(DISALLOWED_LICENSES) $(BASE_GO_PACKAGE)/...
	@echo "--------------------"
	@echo "Licenses check: PASS"

.PHONY: unit-test
unit-test:
	go test -race -cover $(GO_UNIT_TEST_FOLDERS)

.PHONY: int-test
int-test: ENVTEST_ASSETS_DIR = $(shell pwd)/testbin
int-test: ENVTEST_K8S_VERSION = 1.26.1
int-test: export ATLAS_ORG_ID=$(shell grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)
int-test: export ATLAS_PUBLIC_KEY=$(shell grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
int-test: export ATLAS_PRIVATE_KEY=$(shell grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
# magical env that if specified makes the test output 0 on successful runs
# https://github.com/onsi/ginkgo/blob/master/ginkgo/run_command.go#L130
int-test: export GINKGO_EDITOR_INTEGRATION="true"
int-test: generate manifests ## Run integration tests. Sample with labels: `make int-test label=AtlasProject` or `make int-test label='AtlasDeployment && !slow'`
	mkdir -p $(ENVTEST_ASSETS_DIR)
	test -f $(ENVTEST_ASSETS_DIR)/setup-envtest.sh || curl -sSLo $(ENVTEST_ASSETS_DIR)/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.0/hack/setup-envtest.sh
	export ENVTEST_K8S_VERSION=$(ENVTEST_K8S_VERSION) && source $(ENVTEST_ASSETS_DIR)/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); ./scripts/int_local.sh $(label)

.PHONY: e2e
e2e: run-kind ## Run e2e test. Command `make e2e label=cluster-ns` run cluster-ns test
	./scripts/e2e_local.sh $(label) $(build)

.PHONY: e2e-openshift-upgrade
e2e-openshift-upgrade:
	cd scripts && ./openshift-upgrade-test.sh

bin/$(TARGET_OS)/$(TARGET_ARCH):
	mkdir -p $@

bin/$(TARGET_OS)/$(TARGET_ARCH)/manager: $(GO_SOURCES) bin/$(TARGET_OS)/$(TARGET_ARCH)
	@echo "Building operator with version $(VERSION); $(TARGET_OS) - $(TARGET_ARCH)"
	CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $@ -ldflags="-X github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version.Version=$(VERSION)" cmd/manager/main.go
	@touch $@

bin/manager: bin/$(TARGET_OS)/$(TARGET_ARCH)/manager
	cp bin/$(TARGET_OS)/$(TARGET_ARCH)/manager $@

.PHONY: manager
manager: generate fmt vet bin/manager recompute-licenses ## Build manager binary

.PHONY: install
install: manifests kustomize ## Install CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

.PHONY: deploy
deploy: generate manifests run-kind ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	@./scripts/deploy.sh
 
.PHONY: manifests
# Produce CRDs that work back to Kubernetes 1.16 (so 'apiVersion: apiextensions.k8s.io/v1')
manifests: CRD_OPTIONS ?= "crd:crdVersions=v1,ignoreUnexportedFields=true"
manifests: fmt controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	controller-gen $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./pkg/api/v1/..." output:crd:artifacts:config=config/crd/bases
	@./scripts/split_roles_yaml.sh

golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: golangci-lint
	golangci-lint run --timeout 10m

$(TIMESTAMPS_DIR)/fmt: $(GO_SOURCES)
	@echo "goimports -local github.com/mongodb/mongodb-atlas-kubernetes/v2 -l -w \$$(GO_SOURCES)"
	@goimports -local github.com/mongodb/mongodb-atlas-kubernetes/v2 -l -w $(GO_SOURCES)
	@mkdir -p $(TIMESTAMPS_DIR) && touch $@

.PHONY: fmt
fmt: $(TIMESTAMPS_DIR)/fmt ## Run go fmt against code

fix-lint:
	find . -name "*.go" -not -path "./vendor/*" -exec gofmt -w "{}" \;
	goimports -local github.com/mongodb/mongodb-atlas-kubernetes/v2 -w ./pkg
	goimports -local github.com/mongodb/mongodb-atlas-kubernetes/v2 -w ./test
	golangci-lint run --fix

$(TIMESTAMPS_DIR)/vet: $(GO_SOURCES)
	go vet ./...
	@mkdir -p $(TIMESTAMPS_DIR) && touch $@

.PHONY: vet
vet: $(TIMESTAMPS_DIR)/vet ## Run go vet against code

.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.14.0

.PHONY: generate
generate: controller-gen ${GO_SOURCES} ## Generate code
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/api/..."
	$(MAKE) fmt

FILES = $(shell git ls-files -o -m --directory --exclude-standard --no-empty-directory)
TEMP := $(shell echo "$(FILES)" | awk 'NF')
LINES := $(shell echo "$(FILES)" | awk -F=' ' "{ print NF }")
check-missing-files:  ## Check autogenerated runtime objects and manifest files are missing
	@echo "Checking autogenerated runtime objects and manifests"
	@echo "============================================"
ifneq ("$(LINES)","0")
	@echo "Detected files that need to be committed:"
	@echo "$(FILES)" | sed -e "s/^/  /"
	@echo ""
	@echo "Try running: make generate manifests"
	@exit 1
endif
	@echo "--------------------"
	@echo "Check: PASS"

.PHONY: validate-manifests
validate-manifests: generate manifests check-missing-files

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

clean: ## Clean built binaries
	rm -rf bin/*

.PHONY: all-platforms
all-platforms:
	$(MAKE) bin/linux/amd64/manager TARGET_OS=linux TARGET_ARCH=amd64 VERSION=$(VERSION)
	$(MAKE) bin/linux/arm64/manager TARGET_OS=linux TARGET_ARCH=arm64 VERSION=$(VERSION)

.PHONY: all-platforms-docker
all-platforms-docker: all-platforms
	docker build --build-arg BINARY_PATH=bin/linux/amd64 -f fast.Dockerfile -t manager-amd64 .
	docker build --build-arg BINARY_PATH=bin/linux/arm64 -f fast.Dockerfile -t manager-arm64 .

actions.txt: .github/workflows/ ## List GitHub Action dependencies
	./scripts/list-actions.sh > $@

.PHONY: check-major-version
check-major-version: ## Check that VERSION starts with MAJOR_VERSION
	@[[ $(VERSION) == $(MAJOR_VERSION).* ]] && echo "Version OK" || \
	(echo "Bad major version for $(VERSION) expected $(MAJOR_VERSION)"; exit 1)

tools/makejwt/makejwt: tools/makejwt/*.go
	cd tools/makejwt && go test . && go build .

.PHONY: check-version
check-version: ## Check the version is correct & releasable (vX.Y.Z and not "*-dirty" or "unknown")
	VERSION=$(VERSION) BINARY=bin/$(TARGET_OS)/$(TARGET_ARCH)/manager ./scripts/version-check.sh

.PHONY: release-helm
release-helm: tools/makejwt/makejwt ## kick the operator helm chart release
ifndef JWT_RSA_PEM_KEY_BASE64
	@echo "Must set JWT_RSA_PEM_KEY_BASE64 and JWT_APP_ID to $@"
	@exit 1
endif
	@APP_ID=$(JWT_APP_ID) RSA_PEM_KEY_BASE64=$(JWT_RSA_PEM_KEY_BASE64) \
	VERSION=$(VERSION) ./scripts/release-helm.sh

.PHONY: github-token
github-token: tools/makejwt/makejwt ## github-token gets a GitHub token from a key in env vars
ifndef JWT_RSA_PEM_KEY_BASE64
	@echo "Must set JWT_RSA_PEM_KEY_BASE64 and JWT_APP_ID to get $@"
	@exit 1
endif
	@REPO=mongodb/mongodb-atlas-kubernetes APP_ID=$(JWT_APP_ID) \
	RSA_PEM_KEY_BASE64=$(JWT_RSA_PEM_KEY_BASE64) ./scripts/gh-access-token.sh

tools/metrics/metrics: tools/metrics/*.go
	cd tools/metrics && go test . && go build .

.PHONY: slack-report
slack-report: tools/metrics/metrics ## slack a report
ifndef GITHUB_TOKEN
	echo "Getting GitHub token..."
	$(eval GITHUB_TOKEN := $(shell $(MAKE) -s github-token))
endif
	@echo "Computing and sending $(REPORT_TYPE) report to Slack..."
	@GITHUB_TOKEN=$(GITHUB_TOKEN) FORMAT=summary ./tools/metrics/metrics $(REPORT_TYPE) | \
	./scripts/slackit.sh $(SLACK_WEBHOOK)

.PHONY: test-clean
test-clean:
	cd tools/clean && go test ./...

.PHONY: test-makejwt
test-makejwt:
	cd tools/makejwt && go test ./...

.PHONY: test-metrics
test-metrics:
	cd tools/metrics && go test ./...

.PHONY: test-tools ## Test all tools
test-tools: test-clean test-makejwt test-metrics

.PHONY: sign 
sign: ## Sign an AKO multi-architecture image
	@echo "Signing multi-architecture image $(IMG)..."
	IMG=$(IMG) SIGNATURE_REPO=$(SIGNATURE_REPO) ./scripts/sign-multiarch.sh

cosign:
	@which cosign || go install github.com/sigstore/cosign/v2/cmd/cosign@latest

./ako.pem:
	curl $(AKO_SIGN_PUBKEY) > $@

.PHONY: verify 
verify: cosign ./ako.pem ## Verify an AKO multi-architecture image's signature
	@echo "Verifying multi-architecture image signature $(IMG)..."
	IMG=$(IMG) SIGNATURE_REPO=$(SIGNATURE_REPO) \
	./scripts/sign-multiarch.sh verify && echo "VERIFIED OK"

govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: vulncheck
vulncheck: govulncheck ## Run govulncheck to find vulnerabilities in code
	@./scripts/vulncheck.sh ./vuln-ignore

envsubst:
	@which envsubst || go install github.com/drone/envsubst/cmd/envsubst@latest

docker-sbom:
	@docker sbom --help > /dev/null || \
	echo "You might need to install the SBOM plugin for docker, check out docs/dev/release.md#tools"

.PHONY: generate-sboms
generate-sboms: docker-sbom ## Generate a released version SBOMs
	@mkdir -p docs/releases/v$(VERSION) && \
	./scripts/generate_upload_sbom.sh -i $(RELEASED_OPERATOR_IMAGE):$(VERSION) -o docs/releases/v$(VERSION) && \
	ls -l docs/releases/v$(VERSION)

.PHONY: gen-sdlc-checklist
gen-sdlc-checklist: envsubst ## Generate the SDLC checklist
	@VERSION="$(VERSION)" AUTHORS="$(AUTHORS)" ./scripts/gen-sdlc-checklist.sh

# TODO: avoid leaving leftovers in the first place
.PHONY: clear-e2e-leftovers
clear-e2e-leftovers: ## Clear the e2e test leftovers quickly
	git restore bundle* config deploy
	cd helm-charts && git restore .
	git submodule update helm-charts

.PHONY: install-crds
install-crds: ## Install CRDs in Kubernetes
	kubectl apply -k config/crd

.PHONY: set-namespace
set-namespace:
	kubectl create namespace $(OPERATOR_NAMESPACE) || echo "Namespace already in place"

.PHONY: install-credentials
install-credentials: set-namespace ## Install the Atlas credentials for the Operator
	kubectl create secret generic $(ATLAS_KEY_SECRET_NAME) \
	--from-literal="orgId=$(MCLI_ORG_ID)" \
	--from-literal="publicApiKey=$(MCLI_PUBLIC_API_KEY)" \
	--from-literal="privateApiKey=$(MCLI_PRIVATE_API_KEY)" \
	-n "$(OPERATOR_NAMESPACE)" || echo "Secret already in place"
	kubectl label secret -n "${OPERATOR_NAMESPACE}" \
	$(ATLAS_KEY_SECRET_NAME) atlas.mongodb.com/type=credentials

.PHONY: prepare-run
prepare-run: generate vet manifests manager run-kind install-crds install-credentials

.PHONY: run
run: prepare-run ## Run a freshly compiled manager against kind
ifdef RUN_YAML
	kubectl apply -f $(RUN_YAML)
endif
	OPERATOR_POD_NAME=$(OPERATOR_POD_NAME) \
	OPERATOR_NAMESPACE=$(OPERATOR_NAMESPACE) \
	bin/manager --object-deletion-protection=false --log-level=debug \
	--atlas-domain=$(ATLAS_DOMAIN) \
	--global-api-secret-name=$(ATLAS_KEY_SECRET_NAME)

.PHONY: local-docker-build
local-docker-build:
	docker build -t $(LOCAL_IMAGE) .

.PHONY: yq
yq:
	which yq || go install github.com/mikefarah/yq/v4@latest

.PHONY: prepare-all-in-one
prepare-all-in-one: yq local-docker-build run-kind
	kubectl create namespace mongodb-atlas-system || echo "Namespace already in place"
	kind load docker-image $(LOCAL_IMAGE)

.PHONY: test-all-in-one
test-all-in-one: prepare-all-in-one install-credentials ## Test the deploy/all-in-one.yaml definition
	# Test all in one with a local image and at $(ATLAS_DOMAIN) (cloud-qa)
	kubectl apply -f deploy/all-in-one.yaml
	yq deploy/all-in-one.yaml \
	| yq 'select(.kind == "Deployment") | $(CONTAINER_SPEC).imagePullPolicy="IfNotPresent"' \
	| yq 'select(.kind == "Deployment") | $(CONTAINER_SPEC).image="$(LOCAL_IMAGE)"' \
	| yq 'select(.kind == "Deployment") | $(CONTAINER_SPEC).args[0]="--atlas-domain=$(ATLAS_DOMAIN)"' \
	| kubectl apply -f -
.PHONY: upload-sbom-to-silk
upload-sbom-to-silk: ## Upload a give SBOM (lite) to Silk
	@ARTIFACTORY_USERNAME=$(ARTIFACTORY_USERNAME) ARTIFACTORY_PASSWORD=$(ARTIFACTORY_PASSWORD) \
	SILK_CLIENT_ID=$(SILK_CLIENT_ID) SILK_CLIENT_SECRET=$(SILK_CLIENT_SECRET) \
	SILK_ASSET_GROUP=$(SILK_ASSET_GROUP) ./scripts/upload-to-silk.sh $(SBOM_JSON_FILE)

.PHONY: contract-tests
contract-tests: ## Run contract tests
	go clean -testcache
	AKO_CONTRACT_TEST=1 go test -v -race -cover ./test/contract/...
