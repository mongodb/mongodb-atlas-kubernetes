package kube

import (
	"fmt"
	"os"
	"path"
	"strings"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/sethvargo/go-password/password"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
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

func Upgrade(args ...string) {
	args = append([]string{"upgrade"}, args...)
	session := cli.Execute("helm", args...)
	EventuallyWithOffset(1, session.Wait()).Should(Say("STATUS: deployed"), "HELM. Can't upgrade release")
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

func InstallCRD(input model.UserInputs) {
	Install(
		"mongodb-atlas-operator-crds",
		"mongodb/mongodb-atlas-operator-crds",
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

func UninstallCRD(input model.UserInputs) {
	Uninstall("mongodb-atlas-operator-crds", input.Namespace)
}

func InstallKubernetesOperatorWide(input model.UserInputs) {
	repo, tag := splitDockerImage()
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasURL),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

func InstallKubernetesOperatorNS(input model.UserInputs) {
	repo, tag := splitDockerImage()
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasURL),
		"--set-string", fmt.Sprintf("watchNamespaces=%s", input.Namespace),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--namespace="+input.Namespace,
		"--create-namespace",
	)
}

func splitDockerImage() (string, string) {
	url := strings.Split(os.Getenv("IMAGE_URL"), ":")
	Expect(len(url)).Should(Equal(2), "Can't split DOCKER IMAGE")
	return url[0], url[1]
}

func UninstallKubernetesOperator(input model.UserInputs) {
	Uninstall("atlas-operator-"+input.Project.GetProjectName(), input.Namespace)
}

func AddMongoDBRepo() {
	session := cli.Execute("helm", "repo", "add", "mongodb", "https://mongodb.github.io/helm-charts")
	cli.SessionShouldExit(session)
}

func PrepareHelmChartValuesFile(input model.UserInputs) {
	type usersType struct {
		model.UserSpec
		Password string `json:"password,omitempty"`
	}
	type values struct {
		Project model.ProjectSpec `json:"project,omitempty"`
		Mongodb model.ClusterSpec `json:"mongodb,omitempty"`
		Users   []usersType       `json:"users,omitempty"`
	}
	convertType := func(user model.DBUser) usersType {
		var newUser usersType
		newUser.DatabaseName = user.Spec.DatabaseName
		newUser.Labels = user.Spec.Labels
		newUser.Roles = user.Spec.Roles
		newUser.Scopes = user.Spec.Scopes
		newUser.PasswordSecret = user.Spec.PasswordSecret
		newUser.Username = user.Spec.Username
		newUser.DeleteAfterDate = user.Spec.DeleteAfterDate
		return newUser
	}
	newValues := values{input.Project.Spec, input.Clusters[0].Spec, []usersType{}}
	for i := range input.Users {
		secret, _ := password.Generate(10, 3, 0, false, false)
		currentUser := convertType(input.Users[i])
		currentUser.Password = secret
		newValues.Users = append(newValues.Users, currentUser)
	}
	utils.SaveToFile(
		pathToAtlasClusterValuesFile(input),
		utils.JSONToYAMLConvert(newValues),
	)
}

// chart values https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml
func InstallCluster(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	args := prepareHelmChartArgs(input)
	Install(args...)
}

func UpgradeAtlasClusterChart(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	Upgrade(prepareHelmChartArgs(input)...)
}

func prepareHelmChartArgs(input model.UserInputs) []string {
	args := []string{
		input.Clusters[0].Spec.Name,
		"mongodb/atlas-cluster",
		"--set-string", fmt.Sprintf("atlas.orgId=%s", os.Getenv("MCLI_ORG_ID")),
		"--set-string", fmt.Sprintf("atlas.publicApiKey=%s", os.Getenv("MCLI_PUBLIC_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.privateApiKey=%s", os.Getenv("MCLI_PRIVATE_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.connectionSecretName=%s", input.KeyName),

		"--set-string", fmt.Sprintf("project.fullnameOverride=%s", input.Project.GetK8sMetaName()),
		"--set-string", fmt.Sprintf("project.atlasProjectName=%s", input.Project.GetProjectName()),
		"--set-string", fmt.Sprintf("fullnameOverride=%s", input.Clusters[0].ObjectMeta.Name),

		"-f", pathToAtlasClusterValuesFile(input),
		"--namespace=" + input.Namespace,
		"--create-namespace",
	}
	if input.Clusters[0].Spec.ProviderSettings.BackingProviderName == "" {
		args = append(args, "--set", "mongodb.providerSettings.backingProviderName=null") // TODO check
	}
	return args
}

func genSetStringForUsers(input model.UserInputs) []string { // nolint
	// var args []string
	args := make([]string, 0)
	for i, user := range input.Users {
		var roles []string
		secret, _ := password.Generate(10, 3, 0, false, false)
		for k, role := range user.Spec.Roles {
			roles = append(roles,
				"--set", fmt.Sprintf("users[%d].roles[%d].databaseName=%s", i, k, returnNullIfEmpty(role.DatabaseName)),
				"--set", fmt.Sprintf("users[%d].roles[%d].roleName=%s", i, k, returnNullIfEmpty(role.RoleName)),
			)
			if role.CollectionName != "" {
				roles = append(roles,
					"--set", fmt.Sprintf("users[%d].roles[%d].collectionName=%s", i, k, returnNullIfEmpty(role.CollectionName)),
				)
			}
		}
		args = append(args,
			"--set", fmt.Sprintf("users[%d].username=%s", i, user.Spec.Username),
			"--set", fmt.Sprintf("users[%d].password=%s", i, secret),
		)
		if user.Spec.DatabaseName != "" {
			args = append(args,
				"--set", fmt.Sprintf("users[%d].databaseName=%s", i, returnNullIfEmpty(user.Spec.DatabaseName)),
			)
		}
		args = append(args, roles...)
	}
	return args
}

// returnNullIfEmpty if empty. req for the HELM chart
func returnNullIfEmpty(line string) string { // nolint
	if line == "" {
		return "null"
	}
	return line
}

// pathToAtlasClusterValuesFile values for the  atlas-cluster helm chart https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml
func pathToAtlasClusterValuesFile(input model.UserInputs) string {
	return path.Join(input.ProjectPath, "values.yaml")
}
