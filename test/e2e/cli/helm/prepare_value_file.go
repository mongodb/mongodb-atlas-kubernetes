package helm

import (
	"github.com/sethvargo/go-password/password"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/utils"
)

// Prepare chart values file for project, deployments, users https://github.com/mongodb/helm-charts/blob/main/charts/atlas-deployment/values.yaml
func PrepareHelmChartValuesFile(input model.UserInputs) {
	type usersType struct {
		model.UserSpec
		Password string `json:"password,omitempty"`
	}
	type values struct {
		Project     model.ProjectSpec      `json:"project,omitempty"`
		Deployments []model.DeploymentSpec `json:"deployments,omitempty"`
		Users       []usersType            `json:"users,omitempty"`
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
	newValues := values{input.Project.Spec, []model.DeploymentSpec{}, []usersType{}}
	for i := range input.Deployments {
		newValues.Deployments = append(newValues.Deployments, input.Deployments[i].Spec)
	}
	for i := range input.Users {
		secret, _ := password.Generate(10, 3, 0, false, false)
		currentUser := convertType(input.Users[i])
		currentUser.Password = secret
		newValues.Users = append(newValues.Users, currentUser)
	}
	utils.SaveToFile(
		pathToAtlasDeploymentValuesFile(input),
		utils.JSONToYAMLConvert(newValues),
	)
}
