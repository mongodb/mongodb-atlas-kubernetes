.PHONY: generate
generate:
	hack/generate.sh

.PHONY: crds
crds:
	go run main.go --config config.yaml --output crds.yaml

build:
	go build -o bin/openapi2crd main.go
