module github.com/mongodb/mongodb-atlas-kubernetes

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v60.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.19
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.10
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/aws/aws-sdk-go v1.42.25
	github.com/fatih/structtag v1.2.0
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-logr/zapr v0.4.0
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/mongodb-forks/digest v1.0.3
	github.com/mxschmitt/playwright-go v0.1400.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pborman/uuid v1.2.1
	github.com/sethvargo/go-password v0.2.0
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/atlas v0.7.3-0.20210315115044-4b1d3f428c24
	go.mongodb.org/mongo-driver v1.8.1
	go.uber.org/zap v1.19.1
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
)
