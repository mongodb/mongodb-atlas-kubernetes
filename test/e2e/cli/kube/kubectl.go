package kube

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/sethvargo/go-password/password"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

// GenKubeVersion
func GenKubeVersion(fullVersion string) string {
	version := strings.Split(fullVersion, ".")
	return fmt.Sprintf("Major:\"%s\", Minor:\"%s\"", version[0], version[1])
}

// GetPodStatus status.phase
func GetPodStatus(ns string) func() string {
	return func() string {
		session := cli.Execute("kubectl", "get", "pods", "-l", "app.kubernetes.io/instance=mongodb-atlas-kubernetes-operator", "-o", "jsonpath={.items[0].status.phase}", "-n", ns)
		return string(session.Wait("1m").Out.Contents())
	}
}

// DescribeOperatorPod performs "kubectl describe" to get Operator pod information
func DescribeOperatorPod(ns string) string {
	session := cli.Execute("kubectl", "describe", "pods", "-l", "app.kubernetes.io/instance=mongodb-atlas-kubernetes-operator", "-n", ns)
	return string(session.Wait("1m").Out.Contents())
}

// GetGeneration .status.observedGeneration
func GetGeneration(ns, resourceName string) string {
	session := cli.Execute("kubectl", "get", resourceName, "-n", ns, "-o", "jsonpath={.status.observedGeneration}")
	return string(session.Wait("1m").Out.Contents())
}

// GetStatusCondition .status.conditions.type=Ready.status
func GetStatusCondition(ns string, atlasname string) func() string {
	return func() string {
		session := cli.Execute("kubectl", "get", atlasname, "-n", ns, "-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
		return string(session.Wait("1m").Out.Contents())
	}
}

func GetStatusPhase(ns string, args ...string) string {
	args = append([]string{"get"}, args...)
	args = append(args, "-o", "jsonpath={..status.phase}", "-n", ns)
	session := cli.Execute("kubectl", args...)
	return string(session.Wait("1m").Out.Contents())
}

// GetProjectResource
func GetProjectResource(namespace, rName string) v1.AtlasProject {
	session := cli.Execute("kubectl", "get", rName, "-n", namespace, "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var project v1.AtlasProject
	ExpectWithOffset(1, json.Unmarshal(output, &project)).ShouldNot(HaveOccurred())
	return project
}

// GetClusterResource
func GetClusterResource(namespace, rName string) v1.AtlasCluster {
	session := cli.Execute("kubectl", "get", rName, "-n", namespace, "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var cluster v1.AtlasCluster
	ExpectWithOffset(1, json.Unmarshal(output, &cluster)).ShouldNot(HaveOccurred())
	return cluster
}

func GetK8sClusterStateName(ns, rName string) string {
	return GetClusterResource(ns, rName).Status.StateName
}

func DeleteNamespace(ns string) *Buffer {
	session := cli.Execute("kubectl", "delete", "namespace", ns)
	return session.Wait("2m").Out
}

func SwitchContext(name string) {
	session := cli.Execute("kubectl", "config", "use-context", name)
	EventuallyWithOffset(1, session.Wait()).Should(Say("created"))
}

func GetVersionOutput() *Buffer {
	session := cli.Execute("kubectl", "version")
	return session.Wait().Out
}

func Apply(args ...string) *Buffer {
	if args[0] == "-k" {
		args = append([]string{"apply"}, args...)
	} else {
		args = append([]string{"apply", "-f"}, args...)
	}
	session := cli.Execute("kubectl", args...)
	EventuallyWithOffset(1, session).ShouldNot(Say("error"))
	return session.Wait().Out
}

func Delete(args ...string) *Buffer {
	args = append([]string{"delete", "-f"}, args...)
	session := cli.Execute("kubectl", args...)
	return session.Wait("10m").Out
}

func DeleteResource(rType, name, ns string) {
	session := cli.Execute("kubectl", "delete", rType, name, "-n", ns)
	cli.SessionShouldExit(session)
}

func CreateNamespace(name string) *Buffer {
	session := cli.Execute("kubectl", "create", "namespace", name)
	result := cli.GetSessionExitMsg(session)
	ExpectWithOffset(1, result).Should(SatisfyAny(Say("created"), Say("already exists")), "Can't create namespace")
	return session.Out
}

func CreateUserSecret(name, ns string) {
	secret, _ := password.Generate(10, 3, 0, false, false)
	session := cli.ExecuteWithoutWriter("kubectl", "create", "secret", "generic", name,
		"--from-literal=password="+secret,
		"-n", ns,
	)
	EventuallyWithOffset(1, session.Wait()).Should(Say(name + " created"))
}

func CreateApiKeySecret(keyName, ns string) { // TODO add ns
	session := cli.ExecuteWithoutWriter("kubectl", "create", "secret", "generic", keyName,
		"--from-literal=orgId="+os.Getenv("MCLI_ORG_ID"),
		"--from-literal=publicApiKey="+os.Getenv("MCLI_PUBLIC_API_KEY"),
		"--from-literal=privateApiKey="+os.Getenv("MCLI_PRIVATE_API_KEY"),
		"-n", ns,
	)
	EventuallyWithOffset(1, session.Wait()).Should(Say(keyName + " created"))
}

func CreateApiKeySecretFrom(keyName, ns, orgId, public, private string) { // TODO
	session := cli.Execute("kubectl", "create", "secret", "generic", keyName,
		"--from-literal=orgId="+os.Getenv("MCLI_ORG_ID"),
		"--from-literal=publicApiKey="+public,
		"--from-literal=privateApiKey="+private,
		"-n", ns,
	)
	EventuallyWithOffset(1, session.Wait()).Should(Say(keyName + " created"))
}

func DeleteApiKeySecret(keyName, ns string) {
	session := cli.Execute("kubectl", "delete", "secret", keyName, "-n", ns)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}

func GetManagerLogs(ns string) []byte {
	session := cli.ExecuteWithoutWriter("kubectl", "logs", "deploy/mongodb-atlas-operator", "manager", "-n", ns)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	return session.Out.Contents()
}

func GetTestAppLogs(label, ns string) []byte {
	session := cli.ExecuteWithoutWriter("kubectl", "logs", "-l", label, "-n", ns)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	return session.Out.Contents()
}

func DescribeTestApp(label, ns string) []byte {
	session := cli.Execute("kubectl", "describe", "pods", "-l", label, "-n", ns)
	return session.Wait("1m").Out.Contents()
}

func GetYamlResource(resource string, ns string) []byte {
	session := cli.ExecuteWithoutWriter("kubectl", "get", resource, "-o", "yaml", "-n", ns)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	return session.Out.Contents()
}

func CreateConfigMapWithLiterals(configName string, ns string, keys ...string) {
	args := append([]string{"create", "configmap", configName, "-n", ns}, keys...)
	session := cli.Execute("kubectl", args...)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}

func HasConfigMap(configName, ns string) bool {
	session := cli.Execute("kubectl", "get", "configmap", configName, "-n", ns)
	cli.SessionShouldExit(session)
	return session.ExitCode() == 0
}

func GetResourceCreationTimestamp(resource, name, ns string) []byte {
	session := cli.Execute("kubectl", "get", resource, name, "-n", ns, "-o", "jsonpath={.metadata.creationTimestamp}")
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	return session.Out.Contents()
}
