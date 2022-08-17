package e2e_importer_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	k8sClient   client.Client
	atlasClient *atlas.Atlas
)

var _ = BeforeSuite(func() {
	// Add CRDs definitions to client scheme
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))

	// Instantiate the client to interact with k8s cluster
	kubeConfig, err := config.GetConfig()
	if err != nil {
		Fail("Failed to retrieve kube config")
	}
	k8sClient, _ = client.New(kubeConfig, client.Options{
		Scheme: scheme,
	})

	atlasClient = atlas.GetClientOrFail()
})

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}
