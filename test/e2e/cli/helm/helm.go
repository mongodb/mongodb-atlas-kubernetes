package kube

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
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

func matchHelmSearch(match string) string {
	session := cli.Execute("helm", "search", "repo", "mongodb")
	EventuallyWithOffset(1, session, "5m", "10s").Should(gexec.Exit(0))
	content := session.Out.Contents()

	Expect(regexp.MustCompile(match).Match(content)).Should(BeTrue())
	version := regexp.MustCompile(match).FindStringSubmatch(string(content))
	Expect(version).Should(HaveLen(2))
	GinkgoWriter.Write([]byte(fmt.Sprintf("Found version %s for match %s", version[1], match)))
	return version[1]
}

func GetChartVersion(name string) string {
	match := fmt.Sprintf("%s[\\s ]+([\\d.]+)", name)
	return matchHelmSearch(match)
}

func GetAppVersion(name string) string {
	match := fmt.Sprintf("%s[\\s ]+[\\d.]+[\\s ]+([\\d.]+)", name)
	return matchHelmSearch(match)
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
		config.HelmTestAppPath,
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Clusters[0].Spec.Name, user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
	)
}

func RestartTestApplication(input model.UserInputs, user model.DBUser, port string) {
	Upgrade(
		"test-app-"+user.Spec.Username,
		config.HelmTestAppPath,
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Clusters[0].Spec.Name, user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
		"--recreate-pods",
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

func InstallK8sOperatorWide(input model.UserInputs) {
	repo, tag := splitDockerImage()
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

func InstallLatestReleaseOperatorNS(input model.UserInputs) {
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("watchNamespaces=%s", input.Namespace),
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--namespace="+input.Namespace,
		"--create-namespace",
	)
}

func InstallK8sOperatorNS(input model.UserInputs) {
	repo, tag := splitDockerImage()
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
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

func PrepareHelmChartValuesFile(input model.UserInputs, installFromPackage bool) {
	var version string
	if installFromPackage {
		GinkgoWriter.Write([]byte("install from package get version from there"))
		data, _ := ioutil.ReadFile(filepath.Join(config.HelmAtlasResourcesChartPath, "Chart.yaml"))
		match := "version:[\\s ]+([\\d.]+)"
		r, err := regexp.Compile(match)
		Expect(err).ShouldNot(HaveOccurred())
		version = r.FindStringSubmatch(string(data))[1]
	} else {
		version = GetChartVersion("atlas-cluster")
	}

	GinkgoWriter.Write([]byte(version))
	if version == "0.1.6" { // TODO need smth better
		GinkgoWriter.Write([]byte("old version of atlas-cluster chart"))
		PrepareHelmChartValuesFileVersion05(input)
	} else {
		GinkgoWriter.Write([]byte("new version of atlas-cluster chart"))
		PrepareHelmChartValuesFileVersion06(input)
	}
}

// chart values https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml
func PrepareHelmChartValuesFileVersion05(input model.UserInputs) {
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
func PrepareHelmChartValuesFileVersion06(input model.UserInputs) {
	type usersType struct {
		model.UserSpec
		Password string `json:"password,omitempty"`
	}
	type values struct {
		Project  model.ProjectSpec   `json:"project,omitempty"`
		Clusters []model.ClusterSpec `json:"clusters,omitempty"`
		Users    []usersType         `json:"users,omitempty"`
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
	newValues := values{input.Project.Spec, []model.ClusterSpec{}, []usersType{}}
	for i := range input.Clusters {
		newValues.Clusters = append(newValues.Clusters, input.Clusters[i].Spec)
	}
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

func InstallCluster(input model.UserInputs) {
	PrepareHelmChartValuesFile(input, false)
	args := prepareHelmChartArgs(input, "mongodb/atlas-cluster")
	Install(args...)
}

func UpgradeOperatorChart(input model.UserInputs) {
	repo, tag := splitDockerImage()
	packageChart(config.HelmCRDChartPath, filepath.Join(config.HelmOperatorChartPath, "charts"))
	Upgrade(
		"atlas-operator-"+input.Project.GetProjectName(),
		config.HelmOperatorChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"-n", input.Namespace,
		"--wait", "--timeout", "5m",
	)
}

func UpgradeAtlasClusterChart(input model.UserInputs) {
	PrepareHelmChartValuesFile(input, false)
	Upgrade(prepareHelmChartArgs(input, "mongodb/atlas-cluster")...)
}

func UpgradeAtlasClusterChartDev(input model.UserInputs) {
	PrepareHelmChartValuesFile(input, true)
	Upgrade(prepareHelmChartArgs(input, config.HelmAtlasResourcesChartPath)...)
}

func packageChart(sPath, dPath string) {
	session := cli.Execute("helm", "package", sPath, "-d", dPath)
	cli.SessionShouldExit(session)
}

func prepareHelmChartArgs(input model.UserInputs, chartName string) []string {
	args := []string{
		input.Clusters[0].Spec.Name,
		chartName,
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
		"--set", fmt.Sprintf("mongodb.providerSettings.backingProviderName=%s", returnNullIfEmpty(input.Clusters[0].Spec.ProviderSettings.BackingProviderName)),
	}
	return args
}

// returnNullIfEmpty if empty. HELM chart uses --set key
func returnNullIfEmpty(line string) string {
	if line == "" {
		return "null"
	}
	return line
}

// pathToAtlasClusterValuesFile generate path to values file (HELM chart)
// values for the  atlas-cluster helm chart https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml
func pathToAtlasClusterValuesFile(input model.UserInputs) string {
	return path.Join(input.ProjectPath, "values.yaml")
}
