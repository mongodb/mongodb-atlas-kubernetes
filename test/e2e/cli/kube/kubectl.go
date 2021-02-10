package kube

import (
	"fmt"
	"strings"

	"encoding/json"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

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
		session := cli.Execute("kubectl", "get", "pods", "-l", "control-plane=controller-manager", "-o", "jsonpath={.items[0].status.phase}", "-n", ns)
		return string(session.Wait("1m").Out.Contents())
	}
}

// GetGeneration .status.observedGeneration
func GetGeneration(ns string) func() string {
	return func() string {
		session := cli.Execute("kubectl", "get", "atlascluster.atlas.mongodb.com/atlascluster-sample", "-n", ns, "-o", "jsonpath={.status.observedGeneration}")
		return string(session.Wait("1m").Out.Contents())
	}
}

// GetStatusCondition .status.conditions.type=Ready.status
func GetStatusCondition(ns string, atlasname string) func() string {
	return func() string {
		session := cli.Execute("kubectl", "get", atlasname, "-n", ns, "-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
		return string(session.Wait("1m").Out.Contents())
	}
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

func GetK8sClusterStateName(ns, rName string) func() string {
	return func() string {
		return GetClusterResource(ns, rName).Status.StateName
	}
}

func DeleteNamespace(ns string) *Buffer {
	session := cli.Execute("kubectl", "delete", "namespace", ns)
	return session.Wait().Out
}
