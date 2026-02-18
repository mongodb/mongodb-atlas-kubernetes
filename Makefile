# fix for some Linux distros (i.e. WSL)
SHELL := /usr/bin/env bash

CONTAINER_ENGINE ?= docker
CONTROLLER_GEN ?= controller-gen
KUSTOMIZE ?= kustomize
OPERATOR_SDK ?= operator-sdk
JQ ?= jq
YQ ?= yq
AWK ?= awk

DOCKER_SBOM_PLUGIN_VERSION=0.6.1

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION_FILE=version.json
CURRENT_VERSION := $(shell $(JQ) -r .current $(VERSION_FILE))
NEXT_VERSION := $(shell $(JQ) -r .next $(VERSION_FILE))
# Check if there are uncommitted changes
GIT_DIRTY := $(shell git status --porcelain 2>/dev/null | grep -q . && echo "-dirty" || echo "")
VERSION ?= $(CURRENT_VERSION)$(GIT_DIRTY)
BUILDTIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITCOMMIT ?= $(shell git rev-parse --short HEAD 2> /dev/null || true)

# Fix for e2e-gov tests not to use a bad semver version instead
ifdef USE_NEXT_VERSION
VERSION=$(NEXT_VERSION)
endif

# Fix for e2e2 all-in-one test so that it uses an image that already exists in pre-release
ifdef USE_CURRENT_VERSION
VERSION=$(CURRENT_VERSION)
endif

VERSION_PACKAGE = github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version

# LD_FLAGS
LD_FLAGS := -X $(VERSION_PACKAGE).Version=$(VERSION)
LD_FLAGS += -X $(VERSION_PACKAGE).GitCommit=$(GITCOMMIT)
LD_FLAGS += -X $(VERSION_PACKAGE).BuildTime=$(BUILDTIME)
ifdef EXPERIMENTAL
LD_FLAGS += -X $(VERSION_PACKAGE).Experimental=$(EXPERIMENTAL)
endif

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
REGISTRY ?= mongodb
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

# Track changes to sources to avoid repeating operations on sourced when sources did not change
TMPDIR ?= /tmp
TIMESTAMPS_DIR := $(TMPDIR)/mongodb-atlas-kubernetes
GO_SOURCES = $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './tools/*')

# Defaults for make run
OPERATOR_POD_NAME = mongodb-atlas-operator
OPERATOR_NAMESPACE = mongodb-atlas-system
ATLAS_DOMAIN = https://cloud-qa.mongodb.com/
ATLAS_KEY_SECRET_NAME = mongodb-atlas-operator-api-key

# Envtest configuration params
ENVTEST_ASSETS_DIR ?= $(shell pwd)/bin
ENVTEST_K8S_VERSION ?= 1.31.0
KUBEBUILDER_ASSETS ?= $(ENVTEST_ASSETS_DIR)/k8s/$(ENVTEST_K8S_VERSION)-$(TARGET_OS)-$(TARGET_ARCH)

# Ginkgo configuration
GINKGO_NODES ?= 12
GINKGO_EDITOR_INTEGRATION ?= true
GINKGO_OPTS = -vv --randomize-all --output-interceptor-mode=none --trace --timeout 90m --race --nodes=$(GINKGO_NODES) --cover --coverpkg=github.com/mongodb/mongodb-atlas-kubernetes/v2/... --coverprofile=coverprofile.out
GINKGO_FILTER_LABEL ?= $(label)
ifneq ($(GINKGO_FILTER_LABEL),)
GINKGO_FILTER_LABEL_OPT := --label-filter="$(GINKGO_FILTER_LABEL)"
endif
GINKGO=ginkgo run $(GINKGO_OPTS) $(GINKGO_FILTER_LABEL_OPT)

BASE_GO_PACKAGE = github.com/mongodb/mongodb-atlas-kubernetes/v2
DISALLOWED_LICENSES = restricted,reciprocal

SLACK_WEBHOOK ?= https://hooks.slack.com/services/...

# Signature definitions
SIGNATURE_REPO ?= OPERATOR_REGISTRY
AKO_SIGN_PUBKEY = https://cosign.mongodb.com/atlas-kubernetes-operator.pem

OPERATOR_POD_NAME=mongodb-atlas-operator
RUN_YAML= # Set to the YAML to run when calling make run
RUN_LOG_LEVEL ?= debug

LOCAL_IMAGE=mongodb-atlas-kubernetes-operator:compiled
CONTAINER_SPEC=.spec.template.spec.containers[0]

KONDUKTO_REPO="mongodb/mongodb-atlas-kubernetes"
# branch prefix 'main' is always used currently
KONDUKTO_BRANCH_PREFIX=main

HELM_REPO_URL = "https://mongodb.github.io/helm-charts"
HELM_AKO_INSTALL_NAME = local-ako-install
HELM_AKO_NAMESPACE = $(OPERATOR_NAMESPACE)

# Build environment (e.g., dev, prod)
ENV ?= dev

# --- Config Source Directories ---
CONFIG_CRD_BASES := config/crd/bases
CONFIG_MANAGER := config/manager
CONFIG_MANIFESTS_TPL := config/manifests-template
CONFIG_MANIFESTS := config/manifests
CONFIG_CRD := config/crd
CONFIG_RBAC := config/rbac

# ---  Target Directories ---
TARGET_DIR := deploy
CLUSTERWIDE_DIR := $(TARGET_DIR)/clusterwide
NAMESPACED_DIR := $(TARGET_DIR)/namespaced
CRDS_DIR := $(TARGET_DIR)/crds
OPENSHIFT_DIR := $(TARGET_DIR)/openshift
BUNDLE_DIR := bundle
BUNDLE_MANIFESTS_DIR := $(BUNDLE_DIR)/manifests
RELEASE_ALLINONE := config/release/$(ENV)/allinone
RELEASE_CLUSTERWIDE := config/release/$(ENV)/clusterwide
RELEASE_OPENSHIFT := config/release/$(ENV)/openshift
RELEASE_NAMESPACED := config/release/$(ENV)/namespaced
RELEASE_AUTOGENERATED := config/generated

# --- File Targets ---
ALL_IN_ONE_CONFIG := $(TARGET_DIR)/all-in-one.yaml
CLUSTERWIDE_CONFIG := $(CLUSTERWIDE_DIR)/clusterwide-config.yaml
CLUSTERWIDE_CRDS := $(CLUSTERWIDE_DIR)/crds.yaml
NAMESPACED_CONFIG := $(NAMESPACED_DIR)/namespaced-config.yaml
NAMESPACED_CRDS := $(NAMESPACED_DIR)/crds.yaml
OPENSHIFT_CONFIG := $(OPENSHIFT_DIR)/openshift.yaml
OPENSHIFT_CRDS := $(OPENSHIFT_DIR)/crds.yaml
ALL_IN_ONE_AUTOGENERATED_CONFIG := $(TARGET_DIR)/generated/all-in-one.yaml
CSV_FILE := $(BUNDLE_MANIFESTS_DIR)/mongodb-atlas-kubernetes.clusterserviceversion.yaml
BUNDLE_DOCKERFILE := bundle.Dockerfile

GH_RUN_ID=$(shell gh run list -w Test -b main -e schedule -s success --json databaseId | jq '.[0] | .databaseId')

SBOMS_DIR ?= template

SHELLCHECK_OPTIONS ?= -e SC2086

TEST_REGISTRY ?= localhost:5000
DEFAULT_IMAGE_URL := $(TEST_REGISTRY)/mongodb-atlas-kubernetes-operator:$(NEXT_VERSION)-test
export IMAGE_URL

ifndef IMAGE_URL
	IMAGE_URL := $(DEFAULT_IMAGE_URL)
	BUILD_DEPENDENCY := test-docker-image
else
    $(info --- IMAGE_URL is set externally: $(IMAGE_URL))
    BUILD_DEPENDENCY :=
endif

# GO TOOLS
CRD2GO := go tool -modfile=tools/toolbox/go.mod crd2go

# LOCAL TOOLS
OPENAPI2CRD := tools/openapi2crd/bin/openapi2crd
SCAFFOLDER := tools/scaffolder/bin/scaffolder

SCAFFOLDER_FLAGS ?= --all

RH_DRYRUN ?= true
RH_REPOS_DIR ?= $(TMPDIR)/rh-repos
RH_COMMUNITY_OPERATORHUB_REPO_PATH := $(RH_REPOS_DIR)/community-operators
RH_COMMUNITY_OPENSHIFT_REPO_PATH := $(RH_REPOS_DIR)/community-operators-prod
RH_CERTIFIED_OPENSHIFT_REPO_PATH := $(RH_REPOS_DIR)/certified-operators

RELEASES_URL = https://github.com/mongodb/mongodb-atlas-kubernetes/releases/
RELEASE_SBOM_INTEL = $(RELEASES_URL)/download/v$(VERSION)/linux_amd64.sbom.json
RELEASE_SBOM_ARM = $(RELEASES_URL)/download/v$(VERSION)/linux_arm64.sbom.json
RELEASE_SBOM_FILE_INTEL = $(TMPDIR)/linux_amd64.sbom.json
RELEASE_SBOM_FILE_ARM = $(TMPDIR)/linux_arm64.sbom.json

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

.PHONY: compute-labels
compute-labels:
	mkdir -p bin
	go build -o bin/ginkgo-labels tools/compute-test-labels/main.go

.PHONY: build-licenses.csv
build-licenses.csv: go.mod ## Track licenses in a CSV file
	@echo "Tracking licenses into file $@"
	@echo "========================================"
	export GOOS=linux
	export GOARCH=amd64
	GOTOOLCHAIN=local \
	go-licenses csv --include_tests $(BASE_GO_PACKAGE)/... > licenses.csv

.PHONY: check-licenses
check-licenses:  ## Check licenses are compliant with our restrictions
	@echo "Checking licenses not to be: $(DISALLOWED_LICENSES)"
	@echo "============================================"
	export GOOS=linux
	export GOARCH=amd64
	GOTOOLCHAIN=local \
	go-licenses check --include_tests \
	--disallowed_types $(DISALLOWED_LICENSES) $(BASE_GO_PACKAGE)/...
	@echo "--------------------"
	@echo "Licenses check: PASS"

.PHONY: unit-test
unit-test: manifests
	go test -race -cover $(GO_UNIT_TEST_FOLDERS) $(GO_TEST_FLAGS)

## Run integration tests. Sample with labels: `make test/int GINKGO_FILTER_LABEL=AtlasProject`
test/int: envtest manifests
	AKO_INT_TEST=1 KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) $(GINKGO) $(shell pwd)/$@

test/int/clusterwide: envtest manifests
	AKO_INT_TEST=1 KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) $(GINKGO) $(shell pwd)/$@

envtest: envtest-assets
	KUBEBUILDER_ASSETS=$(shell setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(ENVTEST_ASSETS_DIR) -p path)

envtest-assets:
	echo "Env: $(env)"
	mkdir -p $(ENVTEST_ASSETS_DIR)

.PHONY: e2e
e2e: bundle helm-crds manifests run-kind install-crds $(BUILD_DEPENDENCY) ## Run e2e test. Command `make e2e label=cluster-ns` run cluster-ns test
	AKO_E2E_TEST=1 $(GINKGO) $(shell pwd)/test/$@

.PHONY: e2e2
e2e2: bundle helm-crds run-kind manager install-credentials install-crds set-namespace ## Run e2e2 tests. Command `make e2e2 label=integrations-ctlr` run integrations-ctlr e2e2 test
	NO_GORUN=1 \
	AKO_E2E2_TEST=1 \
	OPERATOR_NAMESPACE=$(OPERATOR_NAMESPACE) \
	ginkgo --race --label-filter=$(label) -ldflags="$(LD_FLAGS)" --timeout 120m -vv test/e2e2/

.PHONY: e2e-openshift-upgrade
e2e-openshift-upgrade:
	cd scripts && ./openshift-upgrade-test.sh

bin/$(TARGET_OS)/$(TARGET_ARCH):
	mkdir -p $@

bin/$(TARGET_OS)/$(TARGET_ARCH)/manager: $(GO_SOURCES) bin/$(TARGET_OS)/$(TARGET_ARCH)
	@echo "Building operator with version $(VERSION); $(TARGET_OS) - $(TARGET_ARCH)"
	CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $@ -ldflags="$(LD_FLAGS)" cmd/main.go
	@touch $@

bin/manager: bin/$(TARGET_OS)/$(TARGET_ARCH)/manager
	cp bin/$(TARGET_OS)/$(TARGET_ARCH)/manager $@

.PHONY: manager
manager: generate fmt vet bin/manager ## Build manager binary

.PHONY: build
build: bin/manager

.PHONY: install
install: manifests ## Install CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

.PHONY: manifests
# Produce CRDs that work back to Kubernetes 1.16 (so 'apiVersion: apiextensions.k8s.io/v1')
manifests: CRD_OPTIONS ?= "crd:crdVersions=v1,ignoreUnexportedFields=true"
manifests: ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./api/..." paths="./internal/controller/..." output:crd:artifacts:config=config/crd/bases
	touch config/crd/bases/kustomization.yaml
	sh -c 'cd config/crd/bases; $(KUSTOMIZE) edit add resource *.yaml kustomization.yaml'
	@./scripts/split_roles_yaml.sh
ifdef EXPERIMENTAL
	@if [ -d internal/next-crds ] && find internal/next-crds -maxdepth 1 -name '*.yaml' | grep -q .; then \
	$(CONTROLLER_GEN) crd paths="./internal/nextapi/v1" output:crd:artifacts:config=internal/next-crds; \
	else \
	echo "No experimental CRDs found, skipping apply."; \
	fi
endif

.PHONY: lint
lint: ## Run the lint against the code
	golangci-lint run --timeout 10m

$(TIMESTAMPS_DIR)/fmt: $(GO_SOURCES)
	go tool -modfile=tools/toolbox/go.mod gci write -s standard -s default -s localmodule $(GO_SOURCES)
	@mkdir -p $(TIMESTAMPS_DIR) && touch $@

.PHONY: fmt
fmt: $(TIMESTAMPS_DIR)/fmt ## Run gci (stricter go fmt & imports) against code

$(TIMESTAMPS_DIR)/vet: $(GO_SOURCES)
	go vet ./...
	@mkdir -p $(TIMESTAMPS_DIR) && touch $@

.PHONY: vet
vet: $(TIMESTAMPS_DIR)/vet ## Run go vet against code

.PHONY: generate
generate: ${GO_SOURCES} ## Generate code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/..." paths="./internal/controller/..."
ifdef EXPERIMENTAL
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./internal/nextapi/generated/v1/..."
endif
	go tool -modfile=tools/toolbox/go.mod mockery
	$(MAKE) fmt

.PHONY: check-missing-files
check-missing-files: ## Check autogenerated runtime objects and manifest files are missing
	$(info Checking autogenerated runtime objects and manifests)
	$(info ============================================)
ifneq ($(shell git ls-files -o -m --directory --exclude-standard --no-empty-directory),)
	$(info Detected files that need to be committed:)
	$(info $(shell git ls-files -o -m --directory --exclude-standard --no-empty-directory | sed -e "s/^/  /"))
	$(info )
	$(info Try running: make generate manifests api-docs)
	$(error Check: FAILED)
else
	$(info Check: PASS)
endif

.PHONY: validate-manifests
validate-manifests: generate manifests
	$(MAKE) check-missing-files

.PHONY: helm-crds
helm-crds: bundle
	@cp -r bundle/manifests/atlas.mongodb.com_* helm-charts/atlas-operator-crds/templates/

.PHONY: validate-crds-chart
validate-crds-chart: ## Validate the CRDs in the Helm chart
	@echo "Validating CRDs in the Helm chart"
	@for file in bundle/manifests/atlas.mongodb.com_*.yaml; do \
		helm_file=helm-charts/atlas-operator-crds/templates/$$(basename $$file); \
		if ! cmp -s $$file $$helm_file; then \
			echo "CRD files do not match: $$file and $$helm_file"; \
			exit 1; \
		fi; \
	done
	@echo "All CRD files match"
	@cd helm-charts/atlas-operator-crds && helm template . > /dev/null

.PHONY: image
image: ## Build an operator image for local development
	$(MAKE) bin/linux/amd64/manager TARGET_OS=linux TARGET_ARCH=amd64 VERSION=$(VERSION)
	$(MAKE) bin/linux/arm64/manager TARGET_OS=linux TARGET_ARCH=arm64 VERSION=$(VERSION)
	$(CONTAINER_ENGINE) build -f fast.Dockerfile --build-arg VERSION=$(VERSION) -t $(OPERATOR_IMAGE) .
	$(CONTAINER_ENGINE) push $(OPERATOR_IMAGE)

.PHONY: image-multiarch
image-multiarch: ## Build releaseable multi-architecture operator image
	$(MAKE) bin/linux/amd64/manager TARGET_OS=linux TARGET_ARCH=amd64 VERSION=$(VERSION)
	$(MAKE) bin/linux/arm64/manager TARGET_OS=linux TARGET_ARCH=arm64 VERSION=$(VERSION)
	$(CONTAINER_ENGINE) buildx create --use --name multiarch-builder \
	--driver docker-container --bootstrap || echo "reusing multiarch-builder"
	$(CONTAINER_ENGINE) buildx build --platform linux/amd64,linux/arm64 --push \
	-f fast.Dockerfile --build-arg VERSION=$(VERSION) -t $(OPERATOR_IMAGE) .

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

.PHONY: post-install-hook
post-install-hook:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o bin/helm-post-install cmd/post-install/main.go
	chmod +x bin/helm-post-install

.PHONY: x509-cert
x509-cert: ## Create X.509 cert at path tmp/x509/ (see docs/x509-user.md)
	go run scripts/create_x509.go

.PHONY: clean-gen-crds
clean-gen-crds: ## Clean only generated CRD files
	rm -f config/generated/crd/bases/crds.yaml

clean: clean-gen-crds ## Clean built binaries
	rm -rf bin/*
	rm -rf config/manifests/bases/
	rm -f config/crd/bases/*.yaml
	rm -f helm-charts/atlas-operator-crds/templates/*.yaml
	rm -f config/rbac/clusterwide/role.yaml
	rm -f config/rbac/namespaced/role.yaml
	rm -f config/rbac/role.yaml
	rm -rf deploy/
	rm -rf bundle/
	rm -f bundle.Dockerfile
	rm -rf test/e2e/data/

.PHONY: all-platforms
all-platforms:
	$(MAKE) bin/linux/amd64/manager TARGET_OS=linux TARGET_ARCH=amd64 VERSION=$(VERSION)
	$(MAKE) bin/linux/arm64/manager TARGET_OS=linux TARGET_ARCH=arm64 VERSION=$(VERSION)

# docker-image builds the test image always for linux, even on MacOS.
# This is because the Kubernetes cluster is always run within a Linux VM
.PHONY: docker-image
docker-image: all-platforms
	docker build --build-arg TARGETOS=linux --build-arg TARGETARCH=$(TARGET_ARCH) \
	-f fast.Dockerfile -t $(DEFAULT_IMAGE_URL) .

.PHONY: test-docker-image
test-docker-image: docker-image run-kind
	docker push $(DEFAULT_IMAGE_URL)

.PHONY: check-major-version
check-major-version: ## Check that VERSION starts with MAJOR_VERSION
	@VERSION_MAJOR=$$(echo "$(VERSION)" | cut -d. -f1 | sed 's/v//'); \
	CURRENT_MAJOR=$$(jq -r '.current' $(VERSION_FILE) | cut -d. -f1); \
	[ "$$VERSION_MAJOR" = "$$CURRENT_MAJOR" ] || \
	(echo "Bad major version for $$VERSION expected $$CURRENT_MAJOR"; exit 1)

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

.PHONY: sign
sign: ## Sign an AKO multi-architecture image
	@echo "Signing multi-architecture image $(IMG)..."
	@IMG=$(IMG) SIGNATURE_REPO=$(SIGNATURE_REPO) ./scripts/sign-multiarch.sh

./ako.pem:
	curl $(AKO_SIGN_PUBKEY) > $@

.PHONY: verify
verify: ./ako.pem ## Verify an AKO multi-architecture image's signature
	@echo "Verifying multi-architecture image signature $(IMG)..."
	@IMG=$(IMG) SIGNATURE_REPO=$(SIGNATURE_REPO) \
	./scripts/sign-multiarch.sh verify

.PHONY: push-release-images
push-release-images: ## Push, sign, and verify release images (Phase 6 - point of no return, atomic per target)
	@DEST_PRERELEASE_REPO="$${DEST_PRERELEASE_REPO:-$$DOCKER_PRERELEASE_REPO}" \
	./scripts/push-release-images.sh

.PHONY: prepare-released-branch
prepare-released-branch: ## Checkout released commit and replace CI tooling (Makefile, scripts, devbox files)
	@./scripts/prepare-released-branch.sh $(COMMIT_SHA)

.PHONY: bump-helm-chart-version
bump-helm-chart-version: ## Bump Helm chart version (requires VERSION)
	@VERSION=$(VERSION) ./scripts/bump-helm-chart-version.sh

.PHONY: build-release-pr
build-release-pr: ## Build release artifacts in released-branch directory (requires VERSION, optional: RELEASED_OPERATOR_IMAGE, IMAGE_TAG, AUTHORS, IMAGE_URL, OPERATOR_REGISTRY)
	@if [ ! -d "released-branch" ]; then \
		echo "Error: released-branch directory not found. Run prepare-released-branch first." >&2; \
		exit 1; \
	fi
	@$(MAKE) -C released-branch bundle ENV=prod VERSION=$(VERSION) IMAGE_URL=$(IMAGE_URL) OPERATOR_REGISTRY=$(OPERATOR_REGISTRY)
	@$(MAKE) -C released-branch bump-helm-chart-version VERSION=$(VERSION)
	@$(MAKE) -C released-branch generate-sboms RELEASED_OPERATOR_IMAGE=$(RELEASED_OPERATOR_IMAGE) IMAGE_TAG=$(IMAGE_TAG) VERSION=$(VERSION) SKIP_SIGNATURE_VERIFY=$(SKIP_SIGNATURE_VERIFY)
	@$(MAKE) -C released-branch gen-sdlc-checklist VERSION=$(VERSION) AUTHORS=$(AUTHORS)
	@$(MAKE) -C released-branch build-licenses.csv
	@echo "✓ Release artifacts built in released-branch"

.PHONY: create-release-pr
create-release-pr: ## Create release PR with all updated artifacts (requires GITHUB_TOKEN, creates branches/PRs)
	@./scripts/create-release-pr.sh $(RELEASE_TAG) $(VERSION)

.PHONY: prepare-release-pr
prepare-release-pr: ## Prepare release PR: checkout branch, build artifacts, and create PR (combines prepare-released-branch + build-release-pr + create-release-pr)
	@$(MAKE) prepare-released-branch COMMIT_SHA=$(COMMIT_SHA)
	@$(MAKE) build-release-pr VERSION=$(VERSION)
	@$(MAKE) create-release-pr RELEASE_TAG=$(RELEASE_TAG) VERSION=$(VERSION)

.PHONY: certify-openshift-images
certify-openshift-images: ## Certify OpenShift images using Red Hat preflight
	./scripts/certify-openshift-images.sh

# Helper to set sandbox environment variables for pre-certification (dry-run against prerelease)
# Usage: make pre-cert-sandbox certify-openshift-images
# This sets env vars for testing preflight validation against prerelease image
ifneq ($(filter pre-cert-sandbox,$(MAKECMDGOALS)),)
    PROMOTED_TAG ?= $(error PROMOTED_TAG is not set (e.g., promoted-latest or promoted-<sha>))
    REGISTRY := docker.io
    REPOSITORY := mongodb/mongodb-atlas-kubernetes-operator-prerelease
    VERSION := $(PROMOTED_TAG)
    SUBMIT := false
    # Note: RHCC_PROJECT, RHCC_TOKEN should be set in environment
    # Note: Registry login should be performed before running this target
    export REGISTRY REPOSITORY VERSION SUBMIT
endif

.PHONY: pre-cert-sandbox
pre-cert-sandbox: ## Set sandbox environment variables for pre-certification testing (use with certify-openshift-images)
	@:

# Helper to set sandbox environment variables for certification
# Usage: make cert-sandbox certify-openshift-images
# This sets env vars for testing certification
ifneq ($(filter cert-sandbox,$(MAKECMDGOALS)),)
    VERSION ?= $(error VERSION is not set (e.g., 2.13.1-certified))
    REGISTRY := quay.io
    REPOSITORY := mongodb/mongodb-atlas-kubernetes-operator-prerelease
    SUBMIT := false
    # Note: RHCC_PROJECT, RHCC_TOKEN should be set in environment
    # Note: Registry login should be performed before running this target
    export REGISTRY REPOSITORY VERSION SUBMIT
endif

.PHONY: cert-sandbox
cert-sandbox: ## Set sandbox environment variables for certification testing (use with certify-openshift-images)
	@:

.PHONY: helm-upd-crds
helm-upd-crds:
	HELM_CRDS_PATH=$(HELM_CRDS_PATH) ./scripts/helm-upd-crds.sh

.PHONY: helm-upd-rbac
helm-upd-rbac:
	HELM_RBAC_FILE=$(HELM_RBAC_FILE) ./scripts/helm-upd-rbac.sh

.PHONY: vulncheck
vulncheck: ## Run govulncheck to find vulnerabilities in code
	@./scripts/vulncheck.sh ./vuln-ignore

IMAGE_TAG ?= $(VERSION)
.PHONY: generate-sboms
generate-sboms: ./ako.pem ## Generate a released version SBOMs
	mkdir -p docs/releases/v$(VERSION) && \
	./scripts/generate_upload_sbom.sh -i $(RELEASED_OPERATOR_IMAGE):$(IMAGE_TAG) -o docs/releases/v$(VERSION) && \
	ls -l docs/releases/v$(VERSION)

.PHONY: gen-sdlc-checklist
gen-sdlc-checklist: ## Generate the SDLC checklist
	@VERSION="$(VERSION)" AUTHORS="$(AUTHORS)" ./scripts/gen-sdlc-checklist.sh

.PHONY: install-crds
install-crds: manifests ## Install CRDs in Kubernetes
	kubectl apply -k config/crd
ifdef EXPERIMENTAL
	@if [ -f config/generated/crd/bases/crds.experimental.yaml ]; then \
		echo "Installing experimental CRDs..."; \
		kubectl apply -f config/generated/crd/bases/crds.experimental.yaml; \
	fi
endif

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
prepare-run: generate vet manifests run-kind install-crds install-credentials
	rm -f bin/manager
	$(MAKE) manager VERSION=$(NEXT_VERSION)

.PHONY: run
run: prepare-run ## Run a freshly compiled manager against kind
ifdef RUN_YAML
	kubectl apply -n $(OPERATOR_NAMESPACE) -f $(RUN_YAML)
endif
ifdef BACKGROUND
	@bash -c '(VERSION=$(NEXT_VERSION) \
	OPERATOR_POD_NAME=$(OPERATOR_POD_NAME) \
	OPERATOR_NAMESPACE=$(OPERATOR_NAMESPACE) \
	nohup bin/manager --object-deletion-protection=false --log-level=$(RUN_LOG_LEVEL) \
	--atlas-domain=$(ATLAS_DOMAIN) \
	--global-api-secret-name=$(ATLAS_KEY_SECRET_NAME) > ako.log 2>&1 & echo $$! > ako.pid \
	&& echo "OPERATOR_PID=$$!")'
else
	VERSION=$(NEXT_VERSION) \
	OPERATOR_POD_NAME=$(OPERATOR_POD_NAME) \
	OPERATOR_NAMESPACE=$(OPERATOR_NAMESPACE) \
	bin/manager --object-deletion-protection=false --log-level=$(RUN_LOG_LEVEL) \
	--atlas-domain=$(ATLAS_DOMAIN) \
	--global-api-secret-name=$(ATLAS_KEY_SECRET_NAME)
endif

.PHONY: stop-ako
stop-ako:  
	@kill `cat ako.pid` && rm ako.pid || echo "AKO process not found or already stopped!"  

.PHONY: upload-sbom-to-kondukto
upload-sbom-to-kondukto: ## Upload a given SBOM (lite) file to Kondukto
	@KONDUKTO_REPO=$(KONDUKTO_REPO) KONDUKTO_BRANCH_PREFIX=$(KONDUKTO_BRANCH_PREFIX) \
	./scripts/upload-to-kondukto.sh $(SBOM_JSON_FILE)

.PHONY: augment-sbom
augment-sbom: ## augment the latest SBOM for a given architecture on a given directory
	KONDUKTO_REPO=$(KONDUKTO_REPO) KONDUKTO_BRANCH_PREFIX=$(KONDUKTO_BRANCH_PREFIX) \
	./scripts/augment-sbom.sh $(SBOM_JSON_FILE)

.PHONY: store-augmented-sboms
store-augmented-sboms: ## Augment & Store the latest SBOM for a given version & architecture
	KONDUKTO_BRANCH_PREFIX=$(KONDUKTO_BRANCH_PREFIX) ./scripts/store-sbom-in-s3.sh $(VERSION) $(TARGET_ARCH) $(SBOMS_DIR)

.PHONY: install-ako-helm
install-ako-helm:
	helm repo add mongodb $(HELM_REPO_URL)
	helm upgrade $(HELM_AKO_INSTALL_NAME) mongodb/mongodb-atlas-operator --atomic --install \
	--set-string atlasURI=$(MCLI_OPS_MANAGER_URL) \
	--set objectDeletionProtection=false \
	--set subobjectDeletionProtection=false \
	--namespace=$(HELM_AKO_NAMESPACE) --create-namespace
	kubectl get crds

tools/githubjobs/githubjobs: tools/githubjobs/*.go
	cd tools/githubjobs && go build .

tools/scandeprecation/scandeprecation: tools/scandeprecation/*.go
	cd tools/scandeprecation && go test . && go build .

.PHONY: slack-deprecations
slack-deprecations: tools/scandeprecation/scandeprecation tools/githubjobs/githubjobs
	@echo "Computing and sending deprecation report to Slack..."
	set -o pipefail; GH_RUN_ID=$(GH_RUN_ID) ./tools/githubjobs/githubjobs | grep --text "javaMethod" \
	| ./tools/scandeprecation/scandeprecations | ./scripts/slackit.sh $(SLACK_WEBHOOK)

.PHONY: bump-version-file
bump-version-file:
	@echo "Bumping version in $(VERSION_FILE)..."

	jq '(.next | split(".") | map(tonumber) | .[1] += 1 | .[2] = 0 | map(tostring) | join(".")) as $$new_next | .current = .next | .next = $$new_next' $(VERSION_FILE) > $(VERSION_FILE).tmp && mv $(VERSION_FILE).tmp $(VERSION_FILE)

	@echo "Version updated successfully:"
	@cat $(VERSION_FILE)

.PHONY: api-docs
api-docs: manifests gen-crds
	go tool -modfile=tools/toolbox/go.mod crdoc --resources config/crd/bases --output docs/api-docs.md

.PHONY: validate-api-docs
validate-api-docs: api-docs
	$(MAKE) check-missing-files

.PHONY: addlicenses
addlicense-check:
	addlicense -check -l apache -c "MongoDB Inc" \
	-ignore ".idea/*" \
	-ignore ".vscode/*" \
	-ignore "**/*.md" \
	-ignore "**/*.yaml" \
	-ignore "**/*.yml" \
	-ignore "**/*.nix" \
	-ignore "tools/**" \
	-ignore ".devbox/**" \
	-ignore "tmp/**" \
	-ignore "temp/**" \
	-ignore "**/*Dockerfile" .

.PHONY: shellcheck
shellcheck:
	fd --type f --extension sh --extension bash --extension ksh . | \
	xargs shellcheck --color=always $(SHELLCHECK_OPTIONS)

.PHONY: all-lints
all-lints: fmt lint validate-manifests validate-api-docs check-licenses addlicense-check shellcheck vulncheck
	@echo "✅ CI ALL linting checks PASSED"

.PHONY: ci
ci: unit-test all-lints
	@echo "✅ CI PASSED all checks"

prepare-dirs:
	@echo "Creating directories..."
	@mkdir -p $(BUNDLE_MANIFESTS_DIR)
	@mkdir -p $(CLUSTERWIDE_DIR)
	@mkdir -p $(NAMESPACED_DIR)
	@mkdir -p $(CRDS_DIR)
	@mkdir -p $(OPENSHIFT_DIR)

$(CSV_FILE): prepare-dirs
	@cp $(CONFIG_MANIFESTS_TPL)/bases/mongodb-atlas-kubernetes.clusterserviceversion.yaml $@

update-manager-kustomization:
	@echo "Updating manager kustomization..."
	@cd $(CONFIG_MANAGER); \
		touch kustomization.yaml; \
		$(KUSTOMIZE) edit add resource bases/; \
		$(KUSTOMIZE) edit set image controller="$(OPERATOR_IMAGE)";
	@echo "Manager image set to $(OPERATOR_IMAGE)"

$(ALL_IN_ONE_CONFIG): manifests update-manager-kustomization
	@echo "Building all-in-one manifest..."
	@$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(RELEASE_ALLINONE)" > $@
	@echo "Created $@"

$(CLUSTERWIDE_CONFIG): manifests update-manager-kustomization
	@echo "Building clusterwide config..."
	@$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(RELEASE_CLUSTERWIDE)" > $@
	@echo "Created $@"

$(CLUSTERWIDE_CRDS): manifests update-manager-kustomization
	@echo "Building clusterwide CRDs..."
	@$(KUSTOMIZE) build "$(CONFIG_CRD)" > $@
	@echo "Created $@"

$(NAMESPACED_CONFIG): manifests update-manager-kustomization
	@echo "Building namespaced config..."
	@$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(RELEASE_NAMESPACED)" > $@
	@echo "Created $@"

$(NAMESPACED_CRDS): manifests update-manager-kustomization
	@echo "Building namespaced CRDs..."
	@$(KUSTOMIZE) build "$(CONFIG_CRD)" > $@
	@echo "Created $@"

$(OPENSHIFT_CONFIG): manifests update-manager-kustomization
	@echo "Building OpenShift namespaced config..."
	@$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(RELEASE_OPENSHIFT)" > $@
	@echo "Created $@"

$(OPENSHIFT_CRDS): manifests update-manager-kustomization
	@echo "Building OpenShift namespaced CRDs..."
	@$(KUSTOMIZE) build "$(CONFIG_CRD)" > $@
	@echo "Created $@"

release: $(ALL_IN_ONE_CONFIG) $(CLUSTERWIDE_CONFIG) $(CLUSTERWIDE_CRDS) $(NAMESPACED_CONFIG) $(NAMESPACED_CRDS) $(OPENSHIFT_CONFIG) $(OPENSHIFT_CRDS)
	@find config/crd/bases -type f ! -name 'kustomization.yaml' -exec cp {} "$(CRDS_DIR)" \;
	@echo "Release files prepared successfully."

bundle-validate:
	@echo "Validating bundle..."
	@$(OPERATOR_SDK) bundle validate ./$(BUNDLE_DIR)
	@echo "Bundle validation successful."

bundle: prepare-dirs $(CSV_FILE) release
	@echo "Preparing bundle kustomization..." \
		$(OPERATOR_SDK) generate kustomize manifests --input-dir=$(CONFIG_MANIFESTS_TPL) --interactive=false -q --apis-dir=api; \
		echo "Bundle kustomization generated in $(CONFIG_MANIFESTS).";

	@echo "Building release bundle for version $(VERSION)..."; \
		$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone $(CONFIG_MANIFESTS) | \
			$(OPERATOR_SDK) generate bundle -q --overwrite --version "$(VERSION)" --kustomize-dir="$(CONFIG_MANIFESTS_TPL)" --default-channel="stable" --channels="stable";

	@echo "Patching CSV with replaces: $(CURRENT_VERSION)"; \
		$(AWK) '!/replaces:/' $(CSV_FILE) > $(CSV_FILE).tmp && mv $(CSV_FILE).tmp $(CSV_FILE); \
		echo "  replaces: mongodb-atlas-kubernetes.v$(CURRENT_VERSION)" >> $(CSV_FILE);

	@echo "Patching CSV with WATCH_NAMESPACE..."; \
		value="metadata.annotations['olm.targetNamespaces']" \
		$(YQ) e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].env[2] |= {"name": "WATCH_NAMESPACE", "valueFrom": {"fieldRef": {"fieldPath": env(value)}}}' $(CSV_FILE);

	@echo "Patching CSV with containerImage annotation..."; \
		$(YQ) e -i ".metadata.annotations.containerImage=\"$(OPERATOR_IMAGE)\"" $(CSV_FILE);

	@echo "Patching bundle.Dockerfile with Red Hat labels..."
    	label="LABEL com.redhat.openshift.versions=\"v4.8-v4.18\"\nLABEL com.redhat.delivery.backport=true\nLABEL com.redhat.delivery.operator.bundle=true"; \
    	$(AWK) -v rep="FROM scratch\n\n$${label}" '{sub(/FROM scratch/, rep); print}' $(BUNDLE_DOCKERFILE) > $(BUNDLE_DOCKERFILE).tmp && \
    	mv $(BUNDLE_DOCKERFILE).tmp $(BUNDLE_DOCKERFILE)

bundle-dev: bundle
	@$(AWK) '{gsub(/cloud.mongodb.com/, "cloud-qa.mongodb.com", $$0); print}' $(CSV_FILE) > tmp && mv tmp $(CSV_FILE)
	@echo "✅ Bundle prepared for development."

clean-bundle:
	@echo "🧹 Cleaning up generated files..."
	@rm -rf $(TARGET_DIR) $(BUNDLE_DIR)
	@rm -f $(CONFIG_CRD_BASES)/*.yaml
	@rm -f $(CONFIG_MANAGER)/kustomization.yaml
	@rm -f $(CONFIG_RBAC)/clusterwide/role.yaml
	@rm -f $(CONFIG_RBAC)/namespaced/role.yaml
	@rm -f $(BUNDLE_DOCKERFILE)
	@echo "✅ Cleanup complete."

rbac-autogen:
	$(CONTROLLER_GEN) rbac:roleName=generated-manager-role paths="./internal/generated/controller/..." output:rbac:artifacts:config=config/generated/rbac
	@$(AWK) '/---/{f="xx0"int(++i);} {if(NF!=0)print > f};' config/generated/rbac/role.yaml # Keeping only ClusterRole part while building only all-in-one config
	@rm config/generated/rbac/role.yaml
	@mv xx01 config/generated/rbac/role.yaml

$(ALL_IN_ONE_AUTOGENERATED_CONFIG): manifests update-manager-kustomization rbac-autogen
	@echo "Creating directory..."
	@mkdir -p $(TARGET_DIR)/generated
	@$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone $(RELEASE_AUTOGENERATED) > $@
	@echo "Created $@"

tools/openapi2crd/bin/openapi2crd:
	make -C tools/openapi2crd build

tools/scaffolder/bin/scaffolder:
	make -C tools/scaffolder build

gen-crds: tools/openapi2crd/bin/openapi2crd
	@echo "==> Generating CRDs..."
	$(MAKE) -C tools/openapi2crd build
	$(OPENAPI2CRD) --config config/openapi2crd.yaml \
	--output $(realpath .)/config/generated/crd/bases/crds.yaml
	cp $(realpath .)/config/generated/crd/bases/crds.yaml $(realpath .)/internal/generated/crds/crds.yaml
	$(MAKE) split-generated-crds

.PHONY: split-generated-crds
split-generated-crds: ## Split generated CRDs into production and experimental
	@./scripts/split-generated-crds.sh

.PHONY: regen-crds
regen-crds: clean-gen-crds gen-crds ## Clean and regenerate CRDs

gen-go-types:
	@echo "==> Generating Go models from CRDs..."
	$(CRD2GO) --input $(realpath .)/config/generated/crd/bases/crds.yaml \
	--output $(realpath .)/internal/nextapi/generated/v1

	@echo "==> Generating Go models for scaffolder test CRDs..."
	$(CRD2GO) --input $(realpath .)/test/scaffolder/testdata/crds.yaml \
	--output $(realpath .)/test/scaffolder/generated/types/v1

	@echo "==> Generating Go models for scaffolder test Atlas CRDs..."
	$(CRD2GO) --input $(realpath .)/test/scaffolder/testdata/atlas-crds.yaml \
	--output $(realpath .)/test/scaffolder/generated/types/v1

# In order to override all of the generated versioned handler, use SCAFFOLDER_FLAGS="--all --override" make gen-all
# In order to override a specific generated versioned handler for the Group CRD, use SCAFFOLDER_FLAGS="--kind=Group --override" make gen-all
run-scaffolder: tools/scaffolder/bin/scaffolder
	@echo "==> Generating Go controller scaffolding and indexers..."
	$(MAKE) -C tools/scaffolder build
	$(SCAFFOLDER) --input $(realpath .)/config/generated/crd/bases/crds.yaml \
	$(SCAFFOLDER_FLAGS) \
	--generators indexers,atlas-controllers \
	--indexer-out $(realpath .)/internal/generated/indexers \
	--controller-out $(realpath .)/internal/generated/controller \
	--exporter-out $(realpath .)/pkg/generated/exporter

	@echo "==> Generating scaffolder test exporters for Atlas CRDs..."
	$(SCAFFOLDER) --input $(realpath .)/test/scaffolder/testdata/atlas-crds.yaml \
	$(SCAFFOLDER_FLAGS) \
	--generators atlas-exporters \
	--exporter-out $(realpath .)/test/scaffolder/generated/exporter \
	--types-path github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1

	@echo "==> Generating scaffolder test controller and indexers..."
	$(SCAFFOLDER) --input $(realpath .)/test/scaffolder/testdata/crds.yaml \
	$(SCAFFOLDER_FLAGS) \
	--generators indexers,atlas-controllers \
	--indexer-out $(realpath .)/test/scaffolder/generated/indexers \
	--controller-out $(realpath .)/test/scaffolder/generated/controller \
	--types-path github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1 \
	--indexer-types-path github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1 \
	--indexer-import-path github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/indexers

gen-all: gen-crds gen-go-types run-scaffolder fmt ## Generate all CRDs, Go types, and scaffolding

build-autogen: gen-all $(ALL_IN_ONE_AUTOGENERATED_CONFIG)
	EXPERIMENTAL=1 VERSION=$(shell $(JQ) -r .next $(VERSION_FILE))-EXPERIMENTAL-${GITCOMMIT} $(MAKE) image

.PHONY: create-rh-repos
create-rh-repos: ## Create the Red Hat repos
	./scripts/create-rh-repos.sh $(RH_REPOS_DIR)

.PHONY: reset-rh-creds
reset-rh-creds: ## Reset the Red Hat repository credentials to the value of GH_TOKEN
	GH_TOKEN=$(GH_TOKEN) ./scripts/reset-rh-creds.sh $(RH_REPOS_DIR)

.PHONY: release-rh
release-rh: ## Push the release PRs to the Red Hat repos
	RH_COMMUNITY_OPERATORHUB_REPO_PATH=$(RH_COMMUNITY_OPERATORHUB_REPO_PATH) \
	RH_COMMUNITY_OPENSHIFT_REPO_PATH=$(RH_COMMUNITY_OPENSHIFT_REPO_PATH) \
	RH_CERTIFIED_OPENSHIFT_REPO_PATH=$(RH_CERTIFIED_OPENSHIFT_REPO_PATH) \
	VERSION=$(VERSION) RH_DRYRUN=$(RH_DRYRUN) ./scripts/release-rh.sh

.PHONY: send-sboms
send-sboms: ## Send the SBOMs to Kondukto
	curl -L $(RELEASE_SBOM_INTEL) > $(RELEASE_SBOM_FILE_INTEL)
	curl -L $(RELEASE_SBOM_ARM) > $(RELEASE_SBOM_FILE_ARM)
	make augment-sbom SBOM_JSON_FILE="$(RELEASE_SBOM_FILE_INTEL)"
	make augment-sbom SBOM_JSON_FILE="$(RELEASE_SBOM_FILE_ARM)"

.PHONY: check-kubernetes-versions
check-kubernetes-versions: ## Check the Kubernetes versions are supported
	./scripts/check-kube-versions.sh
