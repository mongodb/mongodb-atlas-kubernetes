package actions

import (
	"context"
	"fmt"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func ProjectCreationFlow(userData *model.TestDataProvider) {
	By("Prepare operator configurations", func() {
		mgr := PrepareOperatorConfigurations(userData)
		ctx := context.Background()
		go func(ctx context.Context) context.Context {
			err := mgr.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			return ctx
		}(ctx)
		deploy.CreateProject(userData)
		userData.ManagerContext = ctx
	})
}

func PrepareOperatorConfigurations(userData *model.TestDataProvider) manager.Manager {
	CreateNamespaceAndSecrets(userData)
	logPath := path.Join("output", userData.Resources.Namespace)
	mgr, err := k8s.RunOperator(&k8s.Config{
		Namespace: userData.Resources.Namespace,
		WatchedNamespaces: map[string]bool{
			userData.Resources.Namespace: true,
		},
		GlobalAPISecret: client.ObjectKey{
			Namespace: userData.Resources.Namespace,
			Name:      config.DefaultOperatorGlobalKey,
		},
		LogDir: logPath,
	})
	Expect(err).NotTo(HaveOccurred())
	return mgr
}

func CreateNamespaceAndSecrets(userData *model.TestDataProvider) {
	By(fmt.Sprintf("Create namespace %s", userData.Resources.Namespace))
	Expect(k8s.CreateNamespace(userData.Context, userData.K8SClient, userData.Resources.Namespace)).Should(Succeed())
	k8s.CreateDefaultSecret(userData.Context, userData.K8SClient, config.DefaultOperatorGlobalKey, userData.Resources.Namespace)
	if !userData.Resources.AtlasKeyAccessType.GlobalLevelKey {
		CreateConnectionAtlasKey(userData)
	}
}
