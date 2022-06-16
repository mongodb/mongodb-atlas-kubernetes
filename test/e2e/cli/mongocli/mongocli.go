package mongocli

import (
	"encoding/json"
	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"go.mongodb.org/atlas/mongodbatlas"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

func GetDeployments(projectID string) []mongodbatlas.Cluster {
	session := cli.Execute("mongocli", "atlas", "clusters", "list", "--projectId", projectID, "-o", "json")
	output := session.Wait("1m").Out.Contents()
	var deployments []mongodbatlas.Cluster
	ExpectWithOffset(1, json.Unmarshal(output, &deployments)).ShouldNot(HaveOccurred())
	return deployments
}

func GetDeploymentByName(projectID string, name string) mongodbatlas.Cluster {
	deployments := GetDeployments(projectID)
	for _, c := range deployments {
		if c.Name == name {
			return c
		}
	}
	panic(fmt.Sprintf("no deployment with name %s in project %s", name, projectID))
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

func GetDeploymentsInfo(projectID string, name string) mongodbatlas.Cluster {
	session := cli.Execute("mongocli", "atlas", "clusters", "describe", name, "--projectId", projectID, "-o", "json")
	EventuallyWithOffset(1, session).Should(gexec.Exit(0))
	output := session.Out.Contents()
	var deployment mongodbatlas.Cluster
	ExpectWithOffset(1, json.Unmarshal(output, &deployment)).ShouldNot(HaveOccurred())
	return deployment
}

func IsProjectInfoExist(projectID string) bool {
	session := cli.Execute("mongocli", "iam", "project", "describe", projectID, "-o", "json")
	cli.SessionShouldExit(session)
	return session.ExitCode() == 0
}

func DeleteDeployment(projectID, deploymentName string) *Buffer {
	session := cli.Execute("mongocli", "atlas", "Deployment", "delete", deploymentName, "--projectId", projectID, "--force")
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

func IsDeploymentExist(projectID string, name string) bool {
	deployments := GetDeployments(projectID)
	for _, c := range deployments {
		GinkgoWriter.Write([]byte(c.Name + "<->" + name + "\n"))
		if c.Name == name {
			return true
		}
	}
	return false
}

func GetDeploymentStateName(projectID string, deploymentName string) string {
	result := GetDeploymentsInfo(projectID, deploymentName)
	return result.StateName
}

func GetVersionOutput() {
	session := cli.Execute("mongocli", "--version")
	session.Wait()
}

func GetUser(database, userName, projectID string) *mongodbatlas.DatabaseUser {
	EventuallyWithOffset(1, IsUserExist(database, userName, projectID), "7m", "10s").Should(BeTrue(), "User doesn't exist")
	session := cli.Execute("mongocli", "atlas", "dbusers", "get", userName, "--authDB", database, "--projectId", projectID, "-o", "json")
	cli.SessionShouldExit(session)
	output := session.Out.Contents()
	var user mongodbatlas.DatabaseUser
	ExpectWithOffset(1, json.Unmarshal(output, &user)).ShouldNot(HaveOccurred())
	return &user
}

func IsUserExist(database, userName, projectID string) bool {
	session := cli.Execute("mongocli", "atlas", "dbusers", "get", userName, "--authDB", database, "--projectId", projectID, "-o", "json")
	cli.SessionShouldExit(session)
	return session.ExitCode() == 0
}

func IsUserExistForAdmin(userName, projectID string) bool {
	session := cli.Execute("mongocli", "atlas", "dbusers", "get", userName, "--projectId", projectID, "-o", "json")
	cli.SessionShouldExit(session)
	return session.ExitCode() == 0
}

func CreateAtlasProjectAPIKey(role, projectID string) (string, string) {
	session := cli.ExecuteWithoutWriter("mongocli", "iam", "project", "apikey", "create", "--projectId", projectID, "--desc", "\"created from the test\"", "--role", role)
	EventuallyWithOffset(1, session.Wait()).Should(Say("created"))
	public := regexp.MustCompile("Public API Key (.+)").FindStringSubmatch(string(session.Out.Contents()))[1]
	private := regexp.MustCompile("Private API Key (.+)").FindStringSubmatch(string(session.Out.Contents()))[1]
	return public, private
}
