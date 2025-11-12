# Define the output binary location
BINARY_DIR  ?= bin
BINARY_NAME := scaffolder
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)

# Define the generated CRD file
CRD_FILE ?= $(realpath ../..)/config/generated/crd/bases/crds.yaml
# Define the output folder for the CRDs
CONTROLLER_OUT ?= $(realpath ../..)/internal/generated/controller
INDEXER_OUT ?= $(realpath ../..)/internal/generated/indexer

build: clean $(BINARY_PATH) ## Build the binary

$(BINARY_PATH): $(GO_FILES) ## File-based build target.
	@echo "==> Building $(BINARY_PATH)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) cmd/scaffolder/main.go

clean:
	@rm -f $(BINARY_PATH)

.PHONY: list
list: build
	$(BINARY_PATH) --input $(CRD_FILE) --list

.PHONY: generate
generate: build
	$(BINARY_PATH) --input $(CRD_FILE) --crd $(CRD_NAME) \
		--indexer-out $(INDEXER_OUT) \
		--controller-out $(CONTROLLER_OUT) \

.PHONY: generate-all
generate-all: build
	$(BINARY_PATH) --input $(CRD_FILE) --all \
		--indexer-out $(INDEXER_OUT) \
		--controller-out $(CONTROLLER_OUT) \

.PHONY: generate-all-override
generate-all-override:
	$(BINARY_PATH) --input $(CRD_FILE) --all --override \
		--indexer-out $(INDEXER_OUT) \
		--controller-out $(CONTROLLER_OUT)

.PHONY: unit-test
unit-test:
	go test -race -cover ./... 

.PHONY: ci
ci: build unit-test ## Standard CI tests
