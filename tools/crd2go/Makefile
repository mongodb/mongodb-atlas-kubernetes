BINARY_DIR  ?= bin
BINARY_NAME := crd2go
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)

build: clean $(BINARY_PATH) ## Build the binary

$(BINARY_PATH): $(GO_FILES) ## File-based build target.
	@echo "==> Building $(BINARY_PATH)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_PATH) cmd/crd2go/main.go

clean:
	rm -f $(BINARY_PATH)
