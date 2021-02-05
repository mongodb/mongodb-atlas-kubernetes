module github.com/mongodb/mongodb-atlas-kubernetes

go 1.15

require (
	github.com/go-logr/zapr v0.1.0
	github.com/google/go-cmp v0.5.4
	github.com/mongodb-forks/digest v1.0.2
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.5
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/atlas v0.7.1
	go.uber.org/zap v1.16.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)
