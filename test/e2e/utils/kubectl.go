package utils

import (
	"fmt"
	"strings"
	// "github.com/onsi/gomega/gexec"
)

// GenKubeVersion
func GenKubeVersion(fullVersion string) string {
	version := strings.Split(fullVersion, ".")
	return fmt.Sprintf("Major:\"%s\", Minor:\"%s\"", version[0], version[1])
}

// GetPodStatus status.phase
func GetPodStatus(ns string) func() string {
	return func() string {
		session := Execute("kubectl", "get", "pods", "-l", "control-plane=controller-manager", "-o", "jsonpath={.items[0].status.phase}", "-n", ns)
		return string(session.Wait("1m").Out.Contents())
	}
}

// GetGeneration .status.observedGeneration
func GetGeneration(nc string) func() string {
	return func() string {
		session := Execute("kubectl", "get", "atlascluster.atlas.mongodb.com/atlascluster-sample", "-n", nc, "-o", "jsonpath={.status.observedGeneration}")
		return string(session.Wait("1m").Out.Contents())
	}
}

// GetStatus .status.conditions.type=Ready.status
func GetStatus(nc string, atlasname string) func() string {
	return func() string {
		session := Execute("kubectl", "get", atlasname, "-n", nc, "-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
		return string(session.Wait("1m").Out.Contents())
	}
}
