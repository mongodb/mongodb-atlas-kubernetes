package utils

import (
	"fmt"
	"os/exec"

	"encoding/json"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"go.mongodb.org/atlas/mongodbatlas"
)

func Execute(command string, args ...string) *gexec.Session {
	// GinkgoWriter.Write([]byte("\n " + command + " " + strings.Join(args, " "))) //TODO for the local run only
	cmd := exec.Command(command, args...)
	session, _ := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	return session
}

func GetClusters(projectID string) []mongodbatlas.Cluster {
	session := Execute("mongocli", "atlas", "clusters", "list", "--projectId", projectID, "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var clusters []mongodbatlas.Cluster
	ExpectWithOffset(1, json.Unmarshal(output, &clusters)).ShouldNot(HaveOccurred())
	return clusters
}

func GetClusterByName(projectID string, name string) mongodbatlas.Cluster {
	clusters := GetClusters(projectID)
	for _, c := range clusters {
		if c.Name == name {
			return c
		}
	}
	panic(fmt.Sprintf("no Cluster with name %s in project %s", name, projectID))
}

func GetProjects() mongodbatlas.Projects {
	session := Execute("mongocli", "iam", "projects", "list", "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var projects mongodbatlas.Projects
	ExpectWithOffset(1, json.Unmarshal(output, &projects)).ShouldNot(HaveOccurred())
	return projects
}

func GetProjectID(name string) string {
	projects := GetProjects()
	for _, p := range projects.Results {
		GinkgoWriter.Write([]byte(p.Name + p.ID + name))
		if p.Name == name {
			return p.ID
		}
	}
	return ""
}

func GetClustersInfo(projectID string, name string) mongodbatlas.Cluster {
	session := Execute("mongocli", "atlas", "clusters", "describe", name, "--projectId", projectID, "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var cluster mongodbatlas.Cluster
	Expect(json.Unmarshal(output, &cluster)).ShouldNot(HaveOccurred())
	return cluster
}

func IsProjectExist(name string) func() bool {
	return func() bool {
		projects := GetProjects().Results
		for _, p := range projects {
			if p.Name == name {
				return true
			}
		}
		return false
	}
}

func IsClusterExist(projectID string, name string) func() bool {
	return func() bool {
		clusters := GetClusters(projectID)
		// if clusters
		for _, c := range clusters {
			GinkgoWriter.Write([]byte(c.Name + name + "\n"))
			if c.Name == name {
				return true
			}
		}
		return false
	}
}

func GetClusterStatus(projectID string, clusterName string) func() string {
	return func() string {
		result := GetClustersInfo(projectID, clusterName)
		return result.StateName
	}
}

//TODO move
func GenKubeVersion(fullVersion string) string {
	version := strings.Split(fullVersion, ".")
	return fmt.Sprintf("Major:\"%s\", Minor:\"%s\"", version[0], version[1])
}

//TODO move
func GetPodStatus(ns string) func() string{
	return func() string{
		session := Execute("kubectl", "get", "pods", "-l", "control-plane=controller-manager", "-o", "jsonpath={.items[0].status.phase}", "-n", ns)
		return string(session.Wait("1m").Out.Contents())
	}
}
