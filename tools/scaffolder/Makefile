build: clean
	go build -o ../../bin/scaffolder ./cmd/main.go

clean:
	rm -f ../../bin/scaffolder

generate:
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/api/..."
