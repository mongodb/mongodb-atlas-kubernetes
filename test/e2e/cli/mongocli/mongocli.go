package mongocli

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"go.mongodb.org/atlas/mongodbatlas"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

func GetClusters(projectID string) []mongodbatlas.Cluster {
	session := cli.Execute("mongocli", "atlas", "clusters", "list", "--projectId", projectID, "-o", "json")
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
	// TODO // mongocli iam projects list -o json >> Error: GET https://cloud-qa.mongodb.com/api/atlas/v1.0/groups: 404 (request "INVALID_ATLAS_GROUP") Atlas group 6026c0654c99af06ac632572 is in an invalid state and cannot be loaded.
	session := cli.Execute("mongocli", "iam", "projects", "list", "-o", "json")
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
	session := cli.Execute("mongocli", "atlas", "clusters", "describe", name, "--projectId", projectID, "-o", "json")
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	output := session.Out.Contents()
	var cluster mongodbatlas.Cluster
	ExpectWithOffset(1, json.Unmarshal(output, &cluster)).ShouldNot(HaveOccurred())
	return cluster
}

func IsProjectInfoExist(projectID string) bool {
	session := cli.Execute("mongocli", "iam", "project", "describe", projectID, "-o", "json")
	Eventually(
		func() int {
			return session.ExitCode()
		},
	).ShouldNot(Equal(-1))
	return session.ExitCode() == 0
}

func DeleteCluster(projectID, clusterName string) *Buffer {
	session := cli.Execute("mongocli", "atlas", "cluster", "delete", clusterName, "--projectId", projectID, "--force")
	return session.Wait().Out
}

func IsProjectExist(name string) bool {
	projects := GetProjects().Results
	for _, p := range projects {
		GinkgoWriter.Write([]byte(p.Name + "<->" + name + "\n"))
		if p.Name == name {
			return true
		}
	}
	return false
}

func IsClusterExist(projectID string, name string) bool {
	clusters := GetClusters(projectID)
	for _, c := range clusters {
		GinkgoWriter.Write([]byte(c.Name + "<->" + name + "\n"))
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

func GetVersionOutput() {
	session := cli.Execute("mongocli", "--version")
	session.Wait(10)
}
