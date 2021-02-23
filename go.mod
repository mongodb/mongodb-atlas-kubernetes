module github.com/mongodb/mongodb-atlas-kubernetes

go 1.15

require (
	github.com/fatih/structtag v1.2.0
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-logr/zapr v0.4.0
	github.com/google/go-cmp v0.5.4
	github.com/mongodb-forks/digest v1.0.2
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/atlas v0.7.2
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0 // indirect
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)
