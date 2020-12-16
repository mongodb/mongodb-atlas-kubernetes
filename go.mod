module github.com/mongodb/mongodb-atlas-kubernetes

go 1.15

require (
	github.com/go-logr/zapr v0.1.0
	github.com/jinzhu/copier v0.1.0
	github.com/mongodb-forks/digest v1.0.1
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/stretchr/testify v1.6.1
	go.mongodb.org/atlas v0.5.1-0.20201208094933-0e2a93147ccd
	go.uber.org/zap v1.16.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)
