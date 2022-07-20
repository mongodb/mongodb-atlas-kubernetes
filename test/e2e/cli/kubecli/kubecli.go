package kube

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/sethvargo/go-password/password"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
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
	cli.SessionShouldExit(session)
	return string(session.Out.Contents())
}

// GetGeneration .status.observedGeneration
func GetGeneration(ns, resourceName string) string {
	session := cli.Execute("kubectl", "get", resourceName, "-n", ns, "-o", "jsonpath={.status.observedGeneration}")
	return string(session.Wait("1m").Out.Contents())
}

// GetStatusCondition .status.conditions.type=Ready.status
func GetStatusCondition(statusType, ns string, atlasname string) string {
	jsonpath := fmt.Sprintf("jsonpath={.status.conditions[?(@.type=='%s')].status}", statusType)
	session := cli.Execute("kubectl", "get", atlasname, "-n", ns, "-o", jsonpath)
	cli.SessionShouldExit(session)
	return string(session.Out.Contents())
}

func GetStatusPhase(ns string, args ...string) string {
	args = append([]string{"get"}, args...)
	args = append(args, "-o", "jsonpath={..status.phase}", "-n", ns)
	session := cli.Execute("kubectl", args...)
	cli.SessionShouldExit(session)
	return string(session.Out.Contents())
}

// GetProjectResource
func GetProjectResource(namespace, rName string) []byte {
	session := cli.Execute("kubectl", "get", rName, "-n", namespace, "-o", "json")
	return session.Wait("1m").Out.Contents()
}

// GetDeploymentResource
func GetDeploymentResource(namespace, rName string) v1.AtlasDeployment {
	session := cli.Execute("kubectl", "get", rName, "-n", namespace, "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var deployment v1.AtlasDeployment
	ExpectWithOffset(1, json.Unmarshal(output, &deployment)).ShouldNot(HaveOccurred())
	return deployment
}

func GetK8sDeploymentStateName(ns, rName string) string {
	return GetDeploymentResource(ns, rName).Status.StateName
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
	EventuallyWithOffset(1, session).ShouldNot(Say("invalid"))
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

func DeleteClusterResource(rType, name string) {
	session := cli.Execute("kubectl", "delete", rType, name)
	cli.SessionShouldExit(session)
}

func CreateNamespace(name string) *Buffer {
	session := cli.Execute("kubectl", "create", "namespace", name)
	result := cli.GetSessionExitMsg(session)
	ExpectWithOffset(1, result).Should(SatisfyAny(Say("created"), Say("already exists")), "Can't create namespace")
	return session.Out
}

func CreateRandomUserSecret(name, ns string) {
	secret, _ := password.Generate(10, 3, 0, false, false)
	CreateUserSecret(secret, name, ns)
}

func CreateUserSecret(secret, name, ns string) {
	session := cli.ExecuteWithoutWriter("kubectl", "create", "secret", "generic", name,
		"--from-literal=password="+secret,
		"-n", ns,
	)
	result := cli.GetSessionExitMsg(session)
	EventuallyWithOffset(1, result).Should(SatisfyAny(Say(name+" created"), Say("already exists")), "Can't create user secret"+name)

	labels := map[string]string{
		connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
	}

	// apply all labels to the secret
	for k, v := range labels {
		session = cli.ExecuteWithoutWriter("kubectl", "label", "secret", name, fmt.Sprintf("%s=%s", k, v), "-n", ns, "--overwrite")
		result = cli.GetSessionExitMsg(session)
		Eventually(result).Should(SatisfyAny(Say("secret/"+name+" labeled"), Say("secret/"+name+" not labeled")))
	}
}

func CreateApiKeySecret(keyName, ns string) {
	session := cli.ExecuteWithoutWriter("kubectl", "create", "secret", "generic", keyName,
		"--from-literal=orgId="+os.Getenv("MCLI_ORG_ID"),
		"--from-literal=publicApiKey="+os.Getenv("MCLI_PUBLIC_API_KEY"),
		"--from-literal=privateApiKey="+os.Getenv("MCLI_PRIVATE_API_KEY"),
		"-n", ns,
	)
	result := cli.GetSessionExitMsg(session)
	EventuallyWithOffset(1, result).Should(SatisfyAny(Say(keyName+" created"), Say("already exists")), "Can't create secret"+keyName)

	session = cli.Execute("kubectl", "label", "secret", keyName, fmt.Sprintf("%s=%s", connectionsecret.TypeLabelKey, connectionsecret.CredLabelVal), "-n", ns, "--overwrite")
	result = cli.GetSessionExitMsg(session)

	// the output is "not labeled" if a label attempt is made and the label already exists with the same value.
	Eventually(result).Should(SatisfyAny(Say("secret/"+keyName+" labeled"), Say("secret/"+keyName+" not labeled")))
}

func CreateApiKeySecretFrom(keyName, ns, orgId, public, private string) {
	session := cli.ExecuteWithoutWriter("kubectl", "create", "secret", "generic", keyName,
		"--from-literal=orgId="+os.Getenv("MCLI_ORG_ID"),
		"--from-literal=publicApiKey="+public,
		"--from-literal=privateApiKey="+private,
		"-n", ns,
	)
	result := cli.GetSessionExitMsg(session)
	EventuallyWithOffset(1, result).Should(SatisfyAny(Say(keyName+" created"), Say("already exists")), "Can't create secret"+keyName)

	session = cli.Execute("kubectl", "label", "secret", keyName, fmt.Sprintf("%s=%s", connectionsecret.TypeLabelKey, connectionsecret.CredLabelVal), "-n", ns, "--overwrite")
	result = cli.GetSessionExitMsg(session)

	// the output is "not labeled" if a label attempt is made and the label already exists with the same value.
	Eventually(result).Should(SatisfyAny(Say("secret/"+keyName+" labeled"), Say("secret/"+keyName+" not labeled")))
}

func DeleteApiKeySecret(keyName, ns string) {
	session := cli.Execute("kubectl", "delete", "secret", keyName, "-n", ns)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}

func CreateX509Secret(keyName, ns string) {
	cert, _, _, err := utils.GenerateX509Cert()
	Expect(err).To(BeNil())

	certFileName := "x509cert.pem"
	certFile, err := os.Create(certFileName)
	Expect(err).To(BeNil())

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	Expect(err).To(BeNil())
	err = certFile.Close()
	Expect(err).To(BeNil())

	session := cli.ExecuteWithoutWriter("kubectl", "create", "secret", "generic", keyName,
		"--from-file="+certFileName,
		"-n", ns,
	)
	result := cli.GetSessionExitMsg(session)
	EventuallyWithOffset(1, result).Should(SatisfyAny(Say(keyName+" created"), Say("already exists")), "Can't create secret"+keyName)

	session = cli.Execute("kubectl", "label", "secret", keyName, fmt.Sprintf("%s=%s", connectionsecret.TypeLabelKey, connectionsecret.CredLabelVal), "-n", ns, "--overwrite")
	result = cli.GetSessionExitMsg(session)

	// the output is "not labeled" if a label attempt is made and the label already exists with the same value.
	Eventually(result).Should(SatisfyAny(Say("secret/"+keyName+" labeled"), Say("secret/"+keyName+" not labeled")))
}

func GetManagerLogs(ns string) []byte {
	session := cli.ExecuteWithoutWriter("kubectl", "logs", "deploy/mongodb-atlas-operator", "manager", "-n", ns)
	cli.SessionShouldExit(session)
	return session.Out.Contents()
}

func GetLogs(label, ns string) []byte {
	session := cli.ExecuteWithoutWriter("kubectl", "logs", "-l", label, "-n", ns)
	cli.SessionShouldExit(session)
	return session.Out.Contents()
}

func DescribeTestApp(label, ns string) []byte {
	session := cli.Execute("kubectl", "describe", "pods", "-l", label, "-n", ns)
	cli.SessionShouldExit(session)
	return session.Out.Contents()
}

func GetYamlResource(resource string, ns string) []byte {
	session := cli.ExecuteWithoutWriter("kubectl", "get", resource, "-o", "yaml", "-n", ns)
	cli.SessionShouldExit(session)
	return session.Out.Contents()
}

func GetJsonResource(resource string, ns string) []byte {
	session := cli.Execute("kubectl", "get", resource, "-n", ns, "-o", "json")
	cli.SessionShouldExit(session)
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
	cli.SessionShouldExit(session)
	return session.Out.Contents()
}

func Annotate(resource, annotation, ns string) {
	session := cli.Execute("kubectl", "annotate", resource, annotation, "-n", ns, "--overwrite=true")
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}

func LabelResourceByLabel(resource, newLabel, ns, labeled string) {
	session := cli.Execute("kubectl", "label", resource, newLabel, "-l", labeled, "-n", ns)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}

func GetPrivateEndpoint(resource, ns string) []byte { // TODO do we need []byte?
	session := cli.Execute("kubectl", "get", resource, "-n", ns, "-o", "jsonpath={.status.privateEndpoints}")
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	return session.Out.Contents()
}

func GetDeploymentDump(output string) {
	outputFolder := fmt.Sprintf("--output-directory=%s", output)
	session := cli.Execute("kubectl", "cluster-info", "dump", "--all-namespaces", outputFolder)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}
