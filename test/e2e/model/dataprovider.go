package model

// Full Data set for the current test case
type TestDataProvider struct {
	ConfPaths       []string                  // init clusters configuration
	ConfUpdatePaths []string                  // update configuration
	Resources       UserInputs                // struct of all user resoucers project,clusters,databaseusers
	Actions         []func(*TestDataProvider) // additional actions for the current data set
	PortGroup       int                       // ports for the test application starts from _
}

func NewTestDataProvider(initClusterConfigs []string, updateClusterConfig []string, users []DBUser, portGroup int, actions []func(*TestDataProvider)) TestDataProvider {
	var data TestDataProvider
	data.ConfPaths = initClusterConfigs
	data.ConfUpdatePaths = updateClusterConfig
	data.Resources = NewUserInputs("my-atlas-key", users)
	data.Actions = actions
	data.PortGroup = portGroup
	return data
}

// func GetNew
