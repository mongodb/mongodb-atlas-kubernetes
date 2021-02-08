package utils

import (
	"fmt"
	"os/exec"

	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
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
	json.Unmarshal(output, &projects)
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
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	output := session.Out.Contents()
	var cluster mongodbatlas.Cluster
	ExpectWithOffset(1, json.Unmarshal(output, &cluster)).ShouldNot(HaveOccurred())
	return cluster
}

func DeleteCluster(projectID, clusterName string) *Buffer {
	session := Execute("mongocli", "atlas", "cluster", "delete", clusterName, "--projectId", projectID, "--force")
	return session.Wait().Out
}

func IsProjectExist(name string) bool {
	projects := GetProjects().Results
	for _, p := range projects {
		if p.Name == name {
			return true
		}
	}
	return false
}

func IsClusterExist(projectID string, name string) bool {
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

func GetClusterStateName(projectID string, clusterName string) string {
	result := GetClustersInfo(projectID, clusterName)
	return result.StateName
}
