package model

// Full Data set for the current test case
type TestDataProvider struct {
	ConfPaths                []string                  // init Deployments configuration
	ConfUpdatePaths          []string                  // update configuration
	Resources                UserInputs                // struct of all user resoucers project,Deployments,databaseusers
	Actions                  []func(*TestDataProvider) // additional actions for the current data set
	PortGroup                int                       // ports for the test application starts from _
	SkipAppConnectivityCheck bool
}

func NewTestDataProvider(keyTestPrefix string, project AProject, r *AtlasKeyType, initDeploymentConfigs []string, updateDeploymentConfig []string, users []DBUser, portGroup int, actions []func(*TestDataProvider)) TestDataProvider {
	var data TestDataProvider
	data.ConfPaths = initDeploymentConfigs
	data.ConfUpdatePaths = updateDeploymentConfig
	data.Resources = NewUserInputs(keyTestPrefix, project, users, r)
	data.Actions = actions
	data.PortGroup = portGroup
	for i := range initDeploymentConfigs {
		data.Resources.Deployments = append(data.Resources.Deployments, LoadUserDeploymentConfig(data.ConfPaths[i]))
	}
	return data
}
