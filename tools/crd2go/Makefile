# Define the output binary location
BINARY_DIR  ?= bin
BINARY_NAME := crd2go
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)

# Define the generated CRD file
CRD_FILE ?= pkg/crd2go/samples/crds.yaml
# Define the output directory for generated Go models
OUTPUT_DIR ?= pkg/crd2go/samples/v1

build: clean $(BINARY_PATH) ## Build the binary

$(BINARY_PATH): $(GO_FILES) ## File-based build target.
	@echo "==> Building $(BINARY_PATH)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) cmd/crd2go/main.go

clean:
	@echo "==> Cleaning..."
	@rm -f $(BINARY_PATH)

generate: build
	@echo "==> Generating Go models from CRDs..."
	@$(BINARY_PATH) --input $(CRD_FILE) --output $(OUTPUT_DIR)