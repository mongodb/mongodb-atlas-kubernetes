package kube

import (
	"fmt"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

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

func GetDeploymentDump(output string) {
	outputFolder := fmt.Sprintf("--output-directory=%s", output)
	session := cli.Execute("kubectl", "cluster-info", "dump", "--all-namespaces", outputFolder)
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
}
