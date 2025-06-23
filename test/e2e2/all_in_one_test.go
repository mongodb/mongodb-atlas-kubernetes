package e2e2_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
)

var _ = Describe("all-in-one.yaml", Ordered, Label("all-in-one"), func() {
	var kubeClient client.Client

	_ = BeforeAll(func() {
		ctx := context.Background()
		c, err := kube.NewK8sTest(ctx)
		kubeClient = c
		Expect(err).To(Succeed())
	})

	It("applies all-in-one.yaml", func() {
		Expect(cli.Execute("kubectl", "apply", "-f", "../../deploy/all-in-one.yaml").Wait().ExitCode()).Should(Equal(0))
	})

	It("waits for mongodb-atlas-operator deployment to be Ready", func() {
		Eventually(func(g Gomega, ctx context.Context) {
			var deployment appsv1.Deployment
			err := kubeClient.Get(ctx, client.ObjectKey{
				Namespace: "mongodb-atlas-system",
				Name:      "mongodb-atlas-operator",
			}, &deployment)
			g.Expect(err).ToNot(HaveOccurred())

			var ready bool
			for _, cond := range deployment.Status.Conditions {
				if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
					ready = true
					break
				}
			}
			g.Expect(ready).To(BeTrue(), "deployment is not Ready")
		}).WithContext(context.Background()).WithPolling(time.Second).WithTimeout(time.Minute).Should(Succeed())
	})
})
