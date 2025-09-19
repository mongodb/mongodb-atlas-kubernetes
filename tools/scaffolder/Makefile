
build:
	go build -o bin/ako-controller-scaffolder cmd/main.go

generate:
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/api/..."
