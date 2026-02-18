// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"encoding/json"
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
	"gopkg.in/yaml.v3"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

// GetVersionOutput returns the helm CLI version
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
	GinkgoWriter.Write(fmt.Appendf(nil, "Found version %s for match %s", version[1], match))
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
	session := cli.Execute("helm", "uninstall", name, "--namespace", ns, "--wait", "--timeout", "15m") // remove timeout
	EventuallyWithOffset(1, session.Wait("15m")).Should(Or(Say("uninstalled"), Say("")), "HELM. Can't uninstall "+name)
}

func Install(args ...string) {
	dependencyAsFileForCRD()
	args = append([]string{"install"}, args...)
	session := cli.Execute("helm", args...)
	msg := cli.GetSessionExitMsg(session)
	ExpectWithOffset(1, msg).Should(SatisfyAny(Say("STATUS: deployed"), Say("already exists"), BeEmpty()),
		"HELM. Can't install release. Message: "+string(msg.Contents()),
	)
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
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Deployments[0].Spec.GetDeploymentName(), user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
	)
}

func RestartTestApplication(input model.UserInputs, user model.DBUser, port string) {
	Upgrade(
		"test-app-"+user.Spec.Username,
		config.TestAppHelmChartPath,
		"--set-string", fmt.Sprintf("connectionSecret=%s-%s-%s", input.Project.GetProjectName(), input.Deployments[0].Spec.GetDeploymentName(), user.Spec.Username),
		"--set-string", fmt.Sprintf("nodePort=%s", port),
		"-n", input.Namespace,
		"--recreate-pods",
	)
}

func InstallCRD(input model.UserInputs) {
	Install("mongodb-atlas-operator-crds"+input.TestID, config.AtlasOperatorCRDHelmChartPath)
}

func UninstallCRD(input model.UserInputs) {
	Uninstall("mongodb-atlas-operator-crds"+input.TestID, "default")
}

func InstallOperatorWideSubmodule(input model.UserInputs) {
	packageChart(config.AtlasOperatorCRDHelmChartPath, filepath.Join(config.AtlasOperatorHelmChartPath, "charts"))
	repo, tag := splitDockerImage()
	createNamespace(input.Namespace)
	installArgs := []string{
		"atlas-operator-" + input.Project.GetProjectName(),
		config.AtlasOperatorHelmChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set", "objectDeletionProtection=false",
		"--set", "subobjectDeletionProtection=false",
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--namespace", input.Namespace,
	}
	pullSecretPassword := os.Getenv("IMAGE_PULL_SECRET_PASSWORD")
	if pullSecretPassword != "" {
		installArgs = addPullSecret(installArgs, pullSecretPassword, input.Namespace)
	}
	Install(installArgs...)
}

// InstallOperatorNamespacedFromLatestRelease install latest released version of the
// Atlas Operator from Helm charts repo.
func InstallOperatorNamespacedFromLatestRelease(input model.UserInputs) {
	createNamespace(input.Namespace)
	installArgs := []string{
		"atlas-operator-" + input.Project.GetProjectName(),
		"mongodb/mongodb-atlas-operator",
		"--set", fmt.Sprintf("watchNamespaces={%s}", input.Namespace),
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set", "objectDeletionProtection=false",
		"--set", "subobjectDeletionProtection=false",
		"--namespace=" + input.Namespace,
	}
	pullSecretPassword := os.Getenv("IMAGE_PULL_SECRET_PASSWORD")
	if pullSecretPassword != "" {
		installArgs = addPullSecret(installArgs, pullSecretPassword, input.Namespace)
	}
	Install(installArgs...)
}

// InstallOperatorNamespacedSubmodule installs the operator from `helm-charts` directory.
// It is expected that this directory already exists.
// mongodb-atlas-operator-crds.enabled=false - because used only for DDT-tests, and CRD deploy there separately
func InstallOperatorNamespacedSubmodule(input model.UserInputs) {
	packageChart(config.AtlasOperatorCRDHelmChartPath, filepath.Join(config.AtlasOperatorHelmChartPath, "charts"))
	repo, tag := splitDockerImage()
	createNamespace(input.Namespace)
	installArgs := []string{
		"atlas-operator-" + input.Project.GetProjectName(),
		config.AtlasOperatorHelmChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--set", fmt.Sprintf("watchNamespaces={%s}", input.Namespace),
		"--set", "mongodb-atlas-operator-crds.enabled=false",
		"--set", "objectDeletionProtection=false",
		"--set", "subobjectDeletionProtection=false",
		fmt.Sprintf("--namespace=%s", input.Namespace),
	}
	pullSecretPassword := os.Getenv("IMAGE_PULL_SECRET_PASSWORD")
	if pullSecretPassword != "" {
		installArgs = addPullSecret(installArgs, pullSecretPassword, input.Namespace)
	}
	Install(installArgs...)
}

func addPullSecret(installArgs []string, pullSecretPassword, namespace string) []string {
	registry := os.Getenv("IMAGE_PULL_SECRET_REGISTRY")
	pullSecretUsername := os.Getenv("IMAGE_PULL_SECRET_USERNAME")
	secretName := pullSecretName(registry)
	createPullSecret(secretName, namespace, registry, pullSecretUsername, pullSecretPassword)
	return appendPullSecretArg(installArgs, secretName)
}

func pullSecretName(registry string) string {
	return fmt.Sprintf("ako-pull-secret-%s", registry)
}

func appendPullSecretArg(installArgs []string, pullSecretName string) []string {
	return append(installArgs,
		"--set-string", fmt.Sprintf("imagePullSecrets[0].name=%s", pullSecretName))
}

func createNamespace(namespace string) {
	session := cli.Execute("kubectl", "create", "namespace", namespace)
	msg := cli.GetSessionExitMsg(session)
	Expect(session.ExitCode()).To(Equal(0), "namespace creation failed: %s", msg)
}

func createPullSecret(secretName, namespace, registry, username, password string) {
	session := cli.Execute(
		"kubectl",
		"create",
		"secret",
		"docker-registry",
		secretName,
		fmt.Sprintf("--namespace=%s", namespace),
		fmt.Sprintf("--docker-server=%s", registry),
		fmt.Sprintf("--docker-username=%s", username),
		fmt.Sprintf("--docker-password=%s", password),
	)
	msg := cli.GetSessionExitMsg(session)
	Expect(session.ExitCode()).To(Equal(0), "pull secret creation failed: %s", msg)
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

// InstallDeploymentSubmodule install the Atlas Deployment Helm Chart from submodule.
func InstallDeploymentSubmodule(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	args := prepareHelmChartArgs(input, config.AtlasDeploymentHelmChartPath)
	Install(args...)
}

// InstallDeploymentRelease from repo
func InstallDeploymentRelease(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	args := prepareHelmChartArgs(input, "mongodb/atlas-deployment")
	Install(args...)
}

func UpgradeOperatorChart(input model.UserInputs) {
	repo, tag := splitDockerImage()
	packageChart(config.AtlasOperatorCRDHelmChartPath, filepath.Join(config.AtlasOperatorHelmChartPath, "charts"))
	upgradeArgs := []string{
		"atlas-operator-" + input.Project.GetProjectName(),
		config.AtlasOperatorHelmChartPath,
		"--set-string", fmt.Sprintf("atlasURI=%s", config.AtlasHost),
		"--set-string", fmt.Sprintf("image.repository=%s", repo),
		"--set-string", fmt.Sprintf("image.tag=%s", tag),
		"--set", "objectDeletionProtection=false",
		"--set", "subobjectDeletionProtection=false",
		"--atomic",
		"-n", input.Namespace,
		// "--wait", "--timeout", "5m", // TODO helm upgrade do not exit
	}
	registry := os.Getenv("IMAGE_PULL_SECRET_REGISTRY")
	upgradeArgs = appendPullSecretArg(upgradeArgs, pullSecretName(registry))
	Upgrade(upgradeArgs...)
}

func UpgradeAtlasDeploymentChartDev(input model.UserInputs) {
	PrepareHelmChartValuesFile(input)
	Upgrade(prepareHelmChartArgs(input, config.AtlasDeploymentHelmChartPath)...)
}

func GetReleasedChartVersion() (string, error) {
	session := cli.Execute("helm", "show", "chart", "mongodb/mongodb-atlas-operator")
	Eventually(session).Should(gexec.Exit(0))
	return getVersionFromChartYAML(session.Out.Contents())
}

func getVersionFromChartYAML(chartYAML []byte) (string, error) {
	chartInfo := map[string]any{}
	err := yaml.Unmarshal(chartYAML, chartInfo)
	if err != nil {
		return "", err
	}
	version, ok := (chartInfo["version"]).(string)
	if !ok {
		return "", fmt.Errorf("not a string at version: %v", chartInfo["version"])
	}
	return version, nil
}

func GetDevelopmentMayorVersion() (string, error) {
	fileBytes, err := os.ReadFile(config.VersionFile)
	if err != nil {
		return "", err
	}

	type versionInfo struct {
		Current string `json:"current"`
		Next    string `json:"next"`
	}
	var versions versionInfo
	if err := json.Unmarshal(fileBytes, &versions); err != nil {
		return "", fmt.Errorf("failed to parse JSON from '%s': %w", config.VersionFile, err)
	}

	parts := strings.Split(versions.Current, ".")
	return parts[0], nil
}

func packageChart(sPath, dPath string) {
	session := cli.Execute("helm", "package", sPath, "-d", dPath)
	cli.SessionShouldExit(session)
}

func prepareHelmChartArgs(input model.UserInputs, chartName string) []string {
	args := []string{
		input.Deployments[0].Spec.GetDeploymentName(),
		chartName,
		"--set-string", fmt.Sprintf("atlas.secret.orgId=%s", os.Getenv("MCLI_ORG_ID")),
		"--set-string", fmt.Sprintf("atlas.secret.publicApiKey=%s", os.Getenv("MCLI_PUBLIC_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.secret.privateApiKey=%s", os.Getenv("MCLI_PRIVATE_API_KEY")),
		"--set-string", fmt.Sprintf("atlas.secret.setCustomName=%s", input.KeyName),

		"--set-string", fmt.Sprintf("project.fullnameOverride=%s", input.Project.GetK8sMetaName()),
		"--set-string", fmt.Sprintf("project.atlasProjectName=%s", input.Project.GetProjectName()),
		"--set-string", fmt.Sprintf("fullnameOverride=%s", input.Deployments[0].ObjectMeta.Name),

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
