package helm

import (
	"io/ioutil"
	"path/filepath"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sethvargo/go-password/password"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func PrepareHelmChartValuesFile(input model.UserInputs, installFromPackage bool) {
	var version string
	if installFromPackage {
		GinkgoWriter.Write([]byte("install from package get version from there"))
		data, _ := ioutil.ReadFile(filepath.Join(config.AtlasClusterHelmChartPath, "Chart.yaml"))
		match := "version:[\\s ]+([\\d.]+)"
		r, err := regexp.Compile(match)
		Expect(err).ShouldNot(HaveOccurred())
		version = r.FindStringSubmatch(string(data))[1]
	} else {
		version = GetChartVersion("atlas-cluster")
	}

	GinkgoWriter.Write([]byte(version))
	if version == "0.1.6" {
		GinkgoWriter.Write([]byte("old version of atlas-cluster chart"))
		PrepareHelmChartValuesFileVersion05(input)
	} else {
		GinkgoWriter.Write([]byte("new version of atlas-cluster chart"))
		PrepareHelmChartValuesFileVersion06(input)
	}
}

// chart values https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml app version 0.5
func PrepareHelmChartValuesFileVersion05(input model.UserInputs) {
	type usersType struct {
		model.UserSpec
		Password string `json:"password,omitempty"`
	}
	type values struct {
		Project model.ProjectSpec `json:"project,omitempty"`
		Mongodb model.ClusterSpec `json:"mongodb,omitempty"`
		Users   []usersType       `json:"users,omitempty"`
	}
	convertType := func(user model.DBUser) usersType {
		var newUser usersType
		newUser.DatabaseName = user.Spec.DatabaseName
		newUser.Labels = user.Spec.Labels
		newUser.Roles = user.Spec.Roles
		newUser.Scopes = user.Spec.Scopes
		newUser.PasswordSecret = user.Spec.PasswordSecret
		newUser.Username = user.Spec.Username
		newUser.DeleteAfterDate = user.Spec.DeleteAfterDate
		return newUser
	}
	newValues := values{input.Project.Spec, input.Clusters[0].Spec, []usersType{}}
	for i := range input.Users {
		secret, _ := password.Generate(10, 3, 0, false, false)
		currentUser := convertType(input.Users[i])
		currentUser.Password = secret
		newValues.Users = append(newValues.Users, currentUser)
	}
	utils.SaveToFile(
		pathToAtlasClusterValuesFile(input),
		utils.JSONToYAMLConvert(newValues),
	)
}

// chart values https://github.com/mongodb/helm-charts/blob/main/charts/atlas-cluster/values.yaml app version 0.6
func PrepareHelmChartValuesFileVersion06(input model.UserInputs) {
	type usersType struct {
		model.UserSpec
		Password string `json:"password,omitempty"`
	}
	type values struct {
		Project  model.ProjectSpec   `json:"project,omitempty"`
		Clusters []model.ClusterSpec `json:"clusters,omitempty"`
		Users    []usersType         `json:"users,omitempty"`
	}
	convertType := func(user model.DBUser) usersType {
		var newUser usersType
		newUser.DatabaseName = user.Spec.DatabaseName
		newUser.Labels = user.Spec.Labels
		newUser.Roles = user.Spec.Roles
		newUser.Scopes = user.Spec.Scopes
		newUser.PasswordSecret = user.Spec.PasswordSecret
		newUser.Username = user.Spec.Username
		newUser.DeleteAfterDate = user.Spec.DeleteAfterDate
		return newUser
	}
	newValues := values{input.Project.Spec, []model.ClusterSpec{}, []usersType{}}
	for i := range input.Clusters {
		newValues.Clusters = append(newValues.Clusters, input.Clusters[i].Spec)
	}
	for i := range input.Users {
		secret, _ := password.Generate(10, 3, 0, false, false)
		currentUser := convertType(input.Users[i])
		currentUser.Password = secret
		newValues.Users = append(newValues.Users, currentUser)
	}
	utils.SaveToFile(
		pathToAtlasClusterValuesFile(input),
		utils.JSONToYAMLConvert(newValues),
	)
}
