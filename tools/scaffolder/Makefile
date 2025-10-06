build:
	go build -o bin/scaffolder cmd/main.go

list: build
	./bin/scaffolder --input ../openapi2crd/result.yaml --list

generate: build
	./bin/scaffolder --input ../openapi2crd/result.yaml --crd $(CRD_NAME) \
		--indexer-out ../../internal/indexer \
		--controller-out ../../internal/controller \
		--translation-out ../../internal/translation


