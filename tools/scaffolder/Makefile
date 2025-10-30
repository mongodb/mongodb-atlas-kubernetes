build: clean
	go build -o ../../bin/scaffolder ./cmd/main.go

clean:
	rm -f ../../bin/scaffolder

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
