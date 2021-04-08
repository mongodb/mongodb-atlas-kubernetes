package kube

import (
	"fmt"
	"os"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/sethvargo/go-password/password"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// GenKubeVersion
func GetVersionOutput() {
	session := cli.Execute("helm", "version")
	ExpectWithOffset(1, session).Should(Say("version.BuildInfo"), "Please, install HELM")
}

func Uninstall(name string, ns string) {
	session := cli.Execute("helm", "uninstall", name, "--namespace", ns)
	EventuallyWithOffset(1, session.Wait()).Should(Say("uninstalled"), "HELM. Can't unninstall "+name)
}

func Install(args ...string) {
	args = append([]string{"install"}, args...)
	session := cli.Execute("helm", args...)
	EventuallyWithOffset(1, session.Wait()).Should(Say("STATUS: deployed"), "HELM. Can't install release")
}

func InstallTestApplication(input model.UserInputs, user model.DBUser, port string) {
	Install(
		"test-app-"+user.Spec.Username,
		config.HELMTestAppPath,
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Clusters[0].Spec.Name, user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
	)
}

func InstallCRDToNamespace(input model.UserInputs) {
	Install(
		"mongodb-atlas-operator-crds",
		"mongodb/mongodb-atlas-operator-crds",
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

func InstallCRD(ns string) {
	Install(
		"mongodb-atlas-operator-crds",
		"mongodb/mongodb-atlas-operator-crds",
	)
}

func UninstallCRDInNamespace(input model.UserInputs) {
	Uninstall("mongodb-atlas-operator-crds", input.Namespace)
}

func InstallKubernetesOperatorWide(input model.UserInputs) {
	Install(
		"atlas-operator",
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasURL),
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

func InstallKubernetesOperatorNS(input model.UserInputs) {
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasURL),
		"--set-string", fmt.Sprintf("watchNamespaces=%s", input.Namespace),
		"--namespace="+input.Namespace,
		"--create-namespace",
	)
}

func UninstallKubernetesOperatorNS(input model.UserInputs) {
	Uninstall("atlas-operator-"+input.Project.GetProjectName(), input.Namespace)
}

func AddMongoDBRepo() {
	session := cli.Execute("helm", "repo", "add", "mongodb", "https://mongodb.github.io/helm-charts")
	cli.SessionShouldExit(session)
}

// chart values https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml
func InstallCluster(input model.UserInputs) {
	// TODO input can have more than one ipadresses/users (need generate args here)
	var args []string
	args = append(args,
		input.Clusters[0].Spec.Name,
		"mongodb/atlas-cluster",
		"--set-string", fmt.Sprintf("atlas.orgId=%s", os.Getenv("MCLI_ORG_ID")),
		"--set-string", fmt.Sprintf("atlas.publicApiKey=%s", os.Getenv("MCLI_PUBLIC_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.privateApiKey=%s", os.Getenv("MCLI_PRIVATE_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.connectionSecretName=%s", input.KeyName),

		"--set-string", fmt.Sprintf("project.fullnameOverride=%s", input.Project.GetK8sMetaName()),
		"--set-string", fmt.Sprintf("project.atlasProjectName=%s", input.Project.GetProjectName()),
		"--set-string", fmt.Sprintf("project.projectIpAccessList[0].ipAddress=%s,project.projectIpAccessList[0].comment=%s",
			input.Project.Spec.ProjectIPAccessList[0].IPAddress, input.Project.Spec.ProjectIPAccessList[0].Comment),

		"--set-string", fmt.Sprintf("fullnameOverride=%s", input.Clusters[0].ObjectMeta.Name),
		"--set-string", fmt.Sprintf("mongodb.providerSettings.providerName=%s", input.Clusters[0].Spec.ProviderSettings.ProviderName),
		"--set-string", fmt.Sprintf("mongodb.providerSettings.regionName=%s", input.Clusters[0].Spec.ProviderSettings.RegionName),
		"--set-string", fmt.Sprintf("mongodb.providerSettings.backingProviderName=%s", input.Clusters[0].Spec.ProviderSettings.BackingProviderName),
		"--namespace="+input.Namespace,
		"--create-namespace",
	)
	args = append(args, genSetStringForUsers(input)...)
	Install(args...)
}

func genSetStringForUsers(input model.UserInputs) []string {
	// var args []string
	args := make([]string, 0)
	for i, user := range input.Users {
		var roles []string
		secret, _ := password.Generate(10, 3, 0, false, false)
		for k, role := range user.Spec.Roles {
			roles = append(roles,
				"--set", fmt.Sprintf("users[%d].roles[%d].databaseName=%s", i, k, returnNullIfEmpty(role.DatabaseName)),
				"--set", fmt.Sprintf("users[%d].roles[%d].roleName=%s", i, k, returnNullIfEmpty(role.RoleName)),
				"--set", fmt.Sprintf("users[%d].roles[%d].collectionName=%s", i, k, returnNullIfEmpty(role.CollectionName)),
			)
		}
		args = append(args,
			"--set", fmt.Sprintf("users[%d].username=%s", i, user.Spec.Username),
			"--set", fmt.Sprintf("users[%d].password=%s", i, secret),
			"--set", fmt.Sprintf("users[%d].databaseName=%s", i, returnNullIfEmpty(user.Spec.DatabaseName)),
		)
		args = append(args, roles...)
	}
	return args
}

// rerunNull if empty. req for the HELM chart
func returnNullIfEmpty(line string) string {
	if line == "" {
		return "null"
	}
	return line
}
