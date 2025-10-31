BINARY_DIR  ?= bin
BINARY_NAME := scaffolder
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)

build: clean $(BINARY_PATH) ## Build the binary

$(BINARY_PATH): $(GO_FILES) ## File-based build target.
	@echo "==> Building $(BINARY_PATH)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) cmd/main.go

clean:
	rm -f $(BINARY_PATH)

.PHONY: list
list: build
	./bin/scaffolder --input ../openapi2crd/crds.yaml --list

.PHONY: generate
generate: build
	./bin/scaffolder --input ../openapi2crd/crds.yaml --crd $(CRD_NAME) \
		--indexer-out ../../internal/indexer \
		--controller-out ../../internal/controller \

.PHONY: generate-all
generate-all: build
	./bin/scaffolder --input ../openapi2crd/crds.yaml --all \
		--indexer-out ../../internal/generated/indexer \
		--controller-out ../../internal/generated/controller \
