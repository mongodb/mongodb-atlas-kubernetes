package actions

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"
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
	By(fmt.Sprintf("Create namespace %s", userData.Resources.Namespace))
	Expect(kubecli.CreateNamespace(userData.Context, userData.K8SClient, userData.Resources.Namespace)).Should(Succeed())
	kubecli.CreateDefaultSecret(userData.Context, userData.K8SClient, config.DefaultOperatorGlobalKey, userData.Resources.Namespace)
	if !userData.Resources.AtlasKeyAccessType.GlobalLevelKey {
		CreateConnectionAtlasKey(userData)
	}
	mgr, err := kubecli.RunOperator(&kubecli.Config{
		Namespace: userData.Resources.Namespace,
		WatchedNamespaces: map[string]bool{
			userData.Resources.Namespace: true,
		},
		GlobalAPISecret: client.ObjectKey{
			Namespace: userData.Resources.Namespace,
			Name:      config.DefaultOperatorGlobalKey,
		},
		LogFileName: fmt.Sprintf("namespace-%s.log", userData.Resources.Namespace),
	})
	Expect(err).NotTo(HaveOccurred())
	return mgr
}
