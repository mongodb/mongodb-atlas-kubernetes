package helm

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
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
	session := cli.Execute("helm", "uninstall", name, "--namespace", ns, "--wait")
	EventuallyWithOffset(1, session.Wait()).Should(Say("uninstalled"), "HELM. Can't unninstall "+name)
}

func Install(args ...string) {
	dependencyAsFileForCRD()
	args = append([]string{"install"}, args...)
	session := cli.Execute("helm", args...)
	EventuallyWithOffset(1, session.Wait()).Should(Say("STATUS: deployed"), "HELM. Can't install release")
}

func Upgrade(args ...string) {
	dependencyAsFileForCRD()
	args = append([]string{"upgrade"}, args...)
	session := cli.Execute("helm", args...)
	EventuallyWithOffset(1, session.Wait()).Should(Say("STATUS: deployed"), "HELM. Can't upgrade release")
}

func InstallTestApplication(input model.UserInputs, user model.DBUser, port string) {
	Install(
		"test-app-"+user.Spec.Username,
		config.TestAppHelmChartPath,
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Clusters[0].Spec.GetClusterName(), user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
	)
}

func RestartTestApplication(input model.UserInputs, user model.DBUser, port string) {
	Upgrade(
		"test-app-"+user.Spec.Username,
		config.TestAppHelmChartPath,
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Clusters[0].Spec.GetClusterName(), user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
		"--recreate-pods",
	)
}

func InstallCRD(input model.UserInputs) {
	Install(
		"mongodb-atlas-operator-crds",
		config.AtlasOperatorCRDHelmChartPath,
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

func UninstallCRD(input model.UserInputs) {
	Uninstall("mongodb-atlas-operator-crds", input.Namespace)
}

func InstallOperatorWideSubmodule(input model.UserInputs) {
	packageChart(config.AtlasOperatorCRDHelmChartPath, filepath.Join(config.AtlasOperatorHelmChartPath, "charts"))
	repo, tag := splitDockerImage()
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		config.AtlasOperatorHelmChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--namespace", input.Namespace,
		"--create-namespace",
	)
}

// InstallOperatorNamespacedFromLatestRelease install latest released version of the
// Atlas Operator from Helm charts repo.
func InstallOperatorNamespacedFromLatestRelease(input model.UserInputs) {
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set-string", fmt.Sprintf("watchNamespaces=%s", input.Namespace),
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--namespace="+input.Namespace,
		"--create-namespace",
	)
}

// InstallOperatorNamespacedSubmodule installs the operator from `helm-charts` directory.
// It is expected that this directory already exists.
func InstallOperatorNamespacedSubmodule(input model.UserInputs) {
	packageChart(config.AtlasOperatorCRDHelmChartPath, filepath.Join(config.AtlasOperatorHelmChartPath, "charts"))
	repo, tag := splitDockerImage()
	Install(
		"atlas-operator-"+input.Project.GetProjectName(),
		config.AtlasOperatorHelmChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("watchNamespaces=%s", input.Namespace),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--namespace="+input.Namespace,
		"--create-namespace",
	)
}

// splitDockerImage returns the image name and tag.
// It splits on the rightmost ":" character to allow for ports
// to be defined in the image name (like `localhost:5000`).
func splitDockerImage() (string, string) {
	imageUrl := os.Getenv("IMAGE_URL")
	// make sure we split on the tag, not on the port ":"
	sepIdx := strings.LastIndex(imageUrl, ":")
	url := []string{imageUrl[:sepIdx], imageUrl[sepIdx+1:]}
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

// InstallClusterSubmodule install the Atlas Cluster Helm Chart from submodule.
func InstallClusterSubmodule(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	args := prepareHelmChartArgs(input, config.AtlasDeploymentHelmChartPath)
	Install(args...)
}

// InstallClusterRelease from repo
func InstallClusterRelease(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	args := prepareHelmChartArgs(input, "mongodb/atlas-deployment")
	Install(args...)
}

func UpgradeOperatorChart(input model.UserInputs) {
	repo, tag := splitDockerImage()
	packageChart(config.AtlasOperatorCRDHelmChartPath, filepath.Join(config.AtlasOperatorHelmChartPath, "charts"))
	Upgrade(
		"atlas-operator-"+input.Project.GetProjectName(),
		config.AtlasOperatorHelmChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"-n", input.Namespace,
		// "--wait", "--timeout", "5m", // TODO helm upgrade do not exit
	)
}

func UpgradeAtlasDeploymentChartDev(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	Upgrade(prepareHelmChartArgs(input, config.AtlasDeploymentHelmChartPath)...)
}

func packageChart(sPath, dPath string) {
	session := cli.Execute("helm", "package", sPath, "-d", dPath)
	cli.SessionShouldExit(session)
}

func prepareHelmChartArgs(input model.UserInputs, chartName string) []string {
	args := []string{
		input.Clusters[0].Spec.GetClusterName(),
		chartName,
		"--set-string", fmt.Sprintf("atlas.secret.orgId=%s", os.Getenv("MCLI_ORG_ID")),
		"--set-string", fmt.Sprintf("atlas.secret.publicApiKey=%s", os.Getenv("MCLI_PUBLIC_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.secret.privateApiKey=%s", os.Getenv("MCLI_PRIVATE_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.secret.setCustomName=%s", input.KeyName),

		"--set-string", fmt.Sprintf("project.fullnameOverride=%s", input.Project.GetK8sMetaName()),
		"--set-string", fmt.Sprintf("project.atlasProjectName=%s", input.Project.GetProjectName()),
		"--set-string", fmt.Sprintf("fullnameOverride=%s", input.Clusters[0].ObjectMeta.Name),

		"-f", pathToAtlasDeploymentValuesFile(input),
		"--namespace=" + input.Namespace,
		"--create-namespace",
	}
	return args
}

// pathToAtlasDeploymentValuesFile generate path to values file (HELM chart)
// values for the  atlas-deployment helm chart https://github.com/mongodb/helm-charts/blob/main/charts/atlas-deployment/values.yaml
func pathToAtlasDeploymentValuesFile(input model.UserInputs) string {
	return path.Join(input.ProjectPath, "values.yaml")
}
