# Define the output binary location
BINARY_DIR  ?= bin
BINARY_NAME := openapi2crd
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
# Define the generated CRD file
CRD_FILE ?= crds.yaml
# Go source files (excluding vendor)
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
# Go packages for testing
PACKAGES := $(shell go list ./...)

# GO TOOLS
GCI := go tool -modfile=../toolbox/go.mod gci

crds: build ## Generate CRDs from config file
	@echo "==> Generating CRDs..."
	$(BINARY_PATH) --config config.yaml --output $(CRD_FILE)

crds-force: ## Generate CRDs from config file
	@echo "==> Generating CRDs..."
	@go run main.go --config config.yaml --force --output $(CRD_FILE)

build: clean $(BINARY_PATH) ## Build the binary

$(BINARY_PATH): $(GO_FILES) ## File-based build target.
	@echo "==> Building $(BINARY_PATH)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) main.go

fmt: ## Format all Go code
	@echo "==> Formatting code..."
	$(GCI) write -s standard -s default -s localmodule .

unit-test: ## Run unit tests with race detection and coverage
	@echo "==> Running unit tests..."
	@go test -race -cover $(GO_TEST_FLAGS) $(PACKAGES)

gen-mock: ## Generate mocks for interfaces
	@echo "==> Generating mocks..."
	@mockery --config .mockery.yaml

all: gen-mock fmt unit-test build

clean: ## Clean up built artifacts
	@echo "==> Cleaning..."
	@rm -f $(BINARY_PATH)
	@rm -f $(CRD_FILE)

.PHONY: ci
ci: all ## Standard CI tests
