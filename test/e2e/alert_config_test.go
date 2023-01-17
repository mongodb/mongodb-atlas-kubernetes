package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("Alert configuration tests", Label("alert-config"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, alertConfigurations []v1.AlertConfiguration) {
			testData = test
			actions.ProjectCreationFlow(test)
			alertConfigFlow(test, alertConfigurations)
		},
		Entry("Test[alert-configs-1]: Project with 2 identical alert configs", Label("alert-configs-1"),
			model.DataProvider(
				"alert-configs-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.AlertConfiguration{
				{
					EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:       true,
					Threshold: &v1.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "1",
						Units:     "HOURS",
					},
					Notifications: []v1.Notification{
						{
							IntervalMin:  5,
							DelayMin:     toptr.MakePtr(5),
							EmailEnabled: toptr.MakePtr(true),
							SMSEnabled:   toptr.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
				{
					EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:       true,
					Threshold: &v1.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "1",
						Units:     "HOURS",
					},
					Notifications: []v1.Notification{
						{
							IntervalMin:  5,
							DelayMin:     toptr.MakePtr(5),
							EmailEnabled: toptr.MakePtr(true),
							SMSEnabled:   toptr.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
			},
		),
		Entry("Test[alert-configs-2]: Project with 2 different alert configs", Label("alert-configs-2"),
			model.DataProvider(
				"alert-configs-2",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.AlertConfiguration{
				{
					EventTypeName: "JOINED_GROUP",
					Enabled:       true,
					Notifications: []v1.Notification{
						{
							IntervalMin:  60,
							DelayMin:     toptr.MakePtr(0),
							EmailEnabled: toptr.MakePtr(true),
							SMSEnabled:   toptr.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
				{
					EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:       true,
					Threshold: &v1.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "1",
						Units:     "HOURS",
					},
					Notifications: []v1.Notification{
						{
							IntervalMin:  5,
							DelayMin:     toptr.MakePtr(5),
							EmailEnabled: toptr.MakePtr(true),
							SMSEnabled:   toptr.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
			},
		),
	)

})

func alertConfigFlow(userData *model.TestDataProvider, alertConfigs []v1.AlertConfiguration) {
	Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
		Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
	userData.Project.Spec.AlertConfigurationSyncEnabled = true
	userData.Project.Spec.AlertConfigurations = append(userData.Project.Spec.AlertConfigurations, alertConfigs...)
	Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())

	actions.WaitForConditionsToBecomeTrue(userData, status.AlertConfigurationReadyType, status.ReadyType)
	Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
	Expect(userData.Project.Status.AlertConfigurations).Should(HaveLen(len(alertConfigs)))

	atlasClient := atlas.GetClientOrFail()
	alertConfigurations, _, err := atlasClient.Client.AlertConfigurations.List(userData.Context, userData.Project.ID(), nil)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(alertConfigurations).Should(HaveLen(len(alertConfigs)), "Atlas alert configurations", alertConfigurations)

	atlasIDList := make([]string, 0, len(alertConfigurations))
	for _, alertConfig := range alertConfigurations {
		atlasIDList = append(atlasIDList, alertConfig.ID)
	}
	statusIDList := make([]string, 0, len(userData.Project.Status.AlertConfigurations))
	for _, alertConfig := range userData.Project.Status.AlertConfigurations {
		statusIDList = append(statusIDList, alertConfig.ID)
	}
	Expect(util.IsEqualWithoutOrder(statusIDList, atlasIDList)).Should(BeTrue())
}
