package model

// Full Data set for the current test case
type TestDataProvider struct {
	ConfPaths       []string                  // init clusters configuration
	ConfUpdatePaths []string                  // update configuration
	Resources       UserInputs                // struct of all user resoucers project,clusters,databaseusers
	Actions         []func(*TestDataProvider) // additional actions for the current data set
	PortGroup       int                       // ports for the test application starts from _
}

func NewTestDataProvider(keyTestPrefix string, r *AtlasKeyType, initClusterConfigs []string, updateClusterConfig []string, users []DBUser, portGroup int, actions []func(*TestDataProvider)) TestDataProvider {
	var data TestDataProvider
	data.ConfPaths = initClusterConfigs
	data.ConfUpdatePaths = updateClusterConfig
	data.Resources = NewUserInputs(keyTestPrefix, users, r)
	data.Actions = actions
	data.PortGroup = portGroup
	for i := range initClusterConfigs {
		data.Resources.Clusters = append(data.Resources.Clusters, LoadUserClusterConfig(data.ConfPaths[i]))
	}
	return data
}
