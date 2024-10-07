// Code generated by mockery. DO NOT EDIT.

package translation

import (
	context "context"

	admin "go.mongodb.org/atlas-sdk/v20231115008/admin"

	mock "github.com/stretchr/testify/mock"

	teams "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

// AtlasTeamsServiceMock is an autogenerated mock type for the AtlasTeamsService type
type AtlasTeamsServiceMock struct {
	mock.Mock
}

type AtlasTeamsServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *AtlasTeamsServiceMock) EXPECT() *AtlasTeamsServiceMock_Expecter {
	return &AtlasTeamsServiceMock_Expecter{mock: &_m.Mock}
}

// AddUsers provides a mock function with given fields: ctx, usersToAdd, orgID, teamID
func (_m *AtlasTeamsServiceMock) AddUsers(ctx context.Context, usersToAdd *[]admin.AddUserToTeam, orgID string, teamID string) error {
	ret := _m.Called(ctx, usersToAdd, orgID, teamID)

	if len(ret) == 0 {
		panic("no return value specified for AddUsers")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *[]admin.AddUserToTeam, string, string) error); ok {
		r0 = rf(ctx, usersToAdd, orgID, teamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AtlasTeamsServiceMock_AddUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUsers'
type AtlasTeamsServiceMock_AddUsers_Call struct {
	*mock.Call
}

// AddUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - usersToAdd *[]admin.AddUserToTeam
//   - orgID string
//   - teamID string
func (_e *AtlasTeamsServiceMock_Expecter) AddUsers(ctx interface{}, usersToAdd interface{}, orgID interface{}, teamID interface{}) *AtlasTeamsServiceMock_AddUsers_Call {
	return &AtlasTeamsServiceMock_AddUsers_Call{Call: _e.mock.On("AddUsers", ctx, usersToAdd, orgID, teamID)}
}

func (_c *AtlasTeamsServiceMock_AddUsers_Call) Run(run func(ctx context.Context, usersToAdd *[]admin.AddUserToTeam, orgID string, teamID string)) *AtlasTeamsServiceMock_AddUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*[]admin.AddUserToTeam), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_AddUsers_Call) Return(_a0 error) *AtlasTeamsServiceMock_AddUsers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AtlasTeamsServiceMock_AddUsers_Call) RunAndReturn(run func(context.Context, *[]admin.AddUserToTeam, string, string) error) *AtlasTeamsServiceMock_AddUsers_Call {
	_c.Call.Return(run)
	return _c
}

// Assign provides a mock function with given fields: ctx, at, projectID
func (_m *AtlasTeamsServiceMock) Assign(ctx context.Context, at *[]teams.Team, projectID string) error {
	ret := _m.Called(ctx, at, projectID)

	if len(ret) == 0 {
		panic("no return value specified for Assign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *[]teams.Team, string) error); ok {
		r0 = rf(ctx, at, projectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AtlasTeamsServiceMock_Assign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Assign'
type AtlasTeamsServiceMock_Assign_Call struct {
	*mock.Call
}

// Assign is a helper method to define mock.On call
//   - ctx context.Context
//   - at *[]teams.Team
//   - projectID string
func (_e *AtlasTeamsServiceMock_Expecter) Assign(ctx interface{}, at interface{}, projectID interface{}) *AtlasTeamsServiceMock_Assign_Call {
	return &AtlasTeamsServiceMock_Assign_Call{Call: _e.mock.On("Assign", ctx, at, projectID)}
}

func (_c *AtlasTeamsServiceMock_Assign_Call) Run(run func(ctx context.Context, at *[]teams.Team, projectID string)) *AtlasTeamsServiceMock_Assign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*[]teams.Team), args[2].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_Assign_Call) Return(_a0 error) *AtlasTeamsServiceMock_Assign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AtlasTeamsServiceMock_Assign_Call) RunAndReturn(run func(context.Context, *[]teams.Team, string) error) *AtlasTeamsServiceMock_Assign_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, at, orgID
func (_m *AtlasTeamsServiceMock) Create(ctx context.Context, at *teams.AssignedTeam, orgID string) (*teams.Team, error) {
	ret := _m.Called(ctx, at, orgID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *teams.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *teams.AssignedTeam, string) (*teams.Team, error)); ok {
		return rf(ctx, at, orgID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *teams.AssignedTeam, string) *teams.Team); ok {
		r0 = rf(ctx, at, orgID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *teams.AssignedTeam, string) error); ok {
		r1 = rf(ctx, at, orgID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AtlasTeamsServiceMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type AtlasTeamsServiceMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - at *teams.AssignedTeam
//   - orgID string
func (_e *AtlasTeamsServiceMock_Expecter) Create(ctx interface{}, at interface{}, orgID interface{}) *AtlasTeamsServiceMock_Create_Call {
	return &AtlasTeamsServiceMock_Create_Call{Call: _e.mock.On("Create", ctx, at, orgID)}
}

func (_c *AtlasTeamsServiceMock_Create_Call) Run(run func(ctx context.Context, at *teams.AssignedTeam, orgID string)) *AtlasTeamsServiceMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*teams.AssignedTeam), args[2].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_Create_Call) Return(_a0 *teams.Team, _a1 error) *AtlasTeamsServiceMock_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AtlasTeamsServiceMock_Create_Call) RunAndReturn(run func(context.Context, *teams.AssignedTeam, string) (*teams.Team, error)) *AtlasTeamsServiceMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamByID provides a mock function with given fields: ctx, orgID, teamID
func (_m *AtlasTeamsServiceMock) GetTeamByID(ctx context.Context, orgID string, teamID string) (*teams.Team, error) {
	ret := _m.Called(ctx, orgID, teamID)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamByID")
	}

	var r0 *teams.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*teams.Team, error)); ok {
		return rf(ctx, orgID, teamID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *teams.Team); ok {
		r0 = rf(ctx, orgID, teamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, orgID, teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AtlasTeamsServiceMock_GetTeamByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamByID'
type AtlasTeamsServiceMock_GetTeamByID_Call struct {
	*mock.Call
}

// GetTeamByID is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamID string
func (_e *AtlasTeamsServiceMock_Expecter) GetTeamByID(ctx interface{}, orgID interface{}, teamID interface{}) *AtlasTeamsServiceMock_GetTeamByID_Call {
	return &AtlasTeamsServiceMock_GetTeamByID_Call{Call: _e.mock.On("GetTeamByID", ctx, orgID, teamID)}
}

func (_c *AtlasTeamsServiceMock_GetTeamByID_Call) Run(run func(ctx context.Context, orgID string, teamID string)) *AtlasTeamsServiceMock_GetTeamByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_GetTeamByID_Call) Return(_a0 *teams.Team, _a1 error) *AtlasTeamsServiceMock_GetTeamByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AtlasTeamsServiceMock_GetTeamByID_Call) RunAndReturn(run func(context.Context, string, string) (*teams.Team, error)) *AtlasTeamsServiceMock_GetTeamByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamByName provides a mock function with given fields: ctx, orgID, teamName
func (_m *AtlasTeamsServiceMock) GetTeamByName(ctx context.Context, orgID string, teamName string) (*teams.Team, error) {
	ret := _m.Called(ctx, orgID, teamName)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamByName")
	}

	var r0 *teams.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*teams.Team, error)); ok {
		return rf(ctx, orgID, teamName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *teams.Team); ok {
		r0 = rf(ctx, orgID, teamName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, orgID, teamName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AtlasTeamsServiceMock_GetTeamByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamByName'
type AtlasTeamsServiceMock_GetTeamByName_Call struct {
	*mock.Call
}

// GetTeamByName is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamName string
func (_e *AtlasTeamsServiceMock_Expecter) GetTeamByName(ctx interface{}, orgID interface{}, teamName interface{}) *AtlasTeamsServiceMock_GetTeamByName_Call {
	return &AtlasTeamsServiceMock_GetTeamByName_Call{Call: _e.mock.On("GetTeamByName", ctx, orgID, teamName)}
}

func (_c *AtlasTeamsServiceMock_GetTeamByName_Call) Run(run func(ctx context.Context, orgID string, teamName string)) *AtlasTeamsServiceMock_GetTeamByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_GetTeamByName_Call) Return(_a0 *teams.Team, _a1 error) *AtlasTeamsServiceMock_GetTeamByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AtlasTeamsServiceMock_GetTeamByName_Call) RunAndReturn(run func(context.Context, string, string) (*teams.Team, error)) *AtlasTeamsServiceMock_GetTeamByName_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamUsers provides a mock function with given fields: ctx, orgID, teamID
func (_m *AtlasTeamsServiceMock) GetTeamUsers(ctx context.Context, orgID string, teamID string) ([]teams.TeamUser, error) {
	ret := _m.Called(ctx, orgID, teamID)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamUsers")
	}

	var r0 []teams.TeamUser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]teams.TeamUser, error)); ok {
		return rf(ctx, orgID, teamID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []teams.TeamUser); ok {
		r0 = rf(ctx, orgID, teamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]teams.TeamUser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, orgID, teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AtlasTeamsServiceMock_GetTeamUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamUsers'
type AtlasTeamsServiceMock_GetTeamUsers_Call struct {
	*mock.Call
}

// GetTeamUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamID string
func (_e *AtlasTeamsServiceMock_Expecter) GetTeamUsers(ctx interface{}, orgID interface{}, teamID interface{}) *AtlasTeamsServiceMock_GetTeamUsers_Call {
	return &AtlasTeamsServiceMock_GetTeamUsers_Call{Call: _e.mock.On("GetTeamUsers", ctx, orgID, teamID)}
}

func (_c *AtlasTeamsServiceMock_GetTeamUsers_Call) Run(run func(ctx context.Context, orgID string, teamID string)) *AtlasTeamsServiceMock_GetTeamUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_GetTeamUsers_Call) Return(_a0 []teams.TeamUser, _a1 error) *AtlasTeamsServiceMock_GetTeamUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AtlasTeamsServiceMock_GetTeamUsers_Call) RunAndReturn(run func(context.Context, string, string) ([]teams.TeamUser, error)) *AtlasTeamsServiceMock_GetTeamUsers_Call {
	_c.Call.Return(run)
	return _c
}

// ListProjectTeams provides a mock function with given fields: ctx, projectID
func (_m *AtlasTeamsServiceMock) ListProjectTeams(ctx context.Context, projectID string) ([]teams.Team, error) {
	ret := _m.Called(ctx, projectID)

	if len(ret) == 0 {
		panic("no return value specified for ListProjectTeams")
	}

	var r0 []teams.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]teams.Team, error)); ok {
		return rf(ctx, projectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []teams.Team); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]teams.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AtlasTeamsServiceMock_ListProjectTeams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListProjectTeams'
type AtlasTeamsServiceMock_ListProjectTeams_Call struct {
	*mock.Call
}

// ListProjectTeams is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
func (_e *AtlasTeamsServiceMock_Expecter) ListProjectTeams(ctx interface{}, projectID interface{}) *AtlasTeamsServiceMock_ListProjectTeams_Call {
	return &AtlasTeamsServiceMock_ListProjectTeams_Call{Call: _e.mock.On("ListProjectTeams", ctx, projectID)}
}

func (_c *AtlasTeamsServiceMock_ListProjectTeams_Call) Run(run func(ctx context.Context, projectID string)) *AtlasTeamsServiceMock_ListProjectTeams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_ListProjectTeams_Call) Return(_a0 []teams.Team, _a1 error) *AtlasTeamsServiceMock_ListProjectTeams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AtlasTeamsServiceMock_ListProjectTeams_Call) RunAndReturn(run func(context.Context, string) ([]teams.Team, error)) *AtlasTeamsServiceMock_ListProjectTeams_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveUser provides a mock function with given fields: ctx, orgID, teamID, userID
func (_m *AtlasTeamsServiceMock) RemoveUser(ctx context.Context, orgID string, teamID string, userID string) error {
	ret := _m.Called(ctx, orgID, teamID, userID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, orgID, teamID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AtlasTeamsServiceMock_RemoveUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveUser'
type AtlasTeamsServiceMock_RemoveUser_Call struct {
	*mock.Call
}

// RemoveUser is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamID string
//   - userID string
func (_e *AtlasTeamsServiceMock_Expecter) RemoveUser(ctx interface{}, orgID interface{}, teamID interface{}, userID interface{}) *AtlasTeamsServiceMock_RemoveUser_Call {
	return &AtlasTeamsServiceMock_RemoveUser_Call{Call: _e.mock.On("RemoveUser", ctx, orgID, teamID, userID)}
}

func (_c *AtlasTeamsServiceMock_RemoveUser_Call) Run(run func(ctx context.Context, orgID string, teamID string, userID string)) *AtlasTeamsServiceMock_RemoveUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_RemoveUser_Call) Return(_a0 error) *AtlasTeamsServiceMock_RemoveUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AtlasTeamsServiceMock_RemoveUser_Call) RunAndReturn(run func(context.Context, string, string, string) error) *AtlasTeamsServiceMock_RemoveUser_Call {
	_c.Call.Return(run)
	return _c
}

// RenameTeam provides a mock function with given fields: ctx, at, orgID, newName
func (_m *AtlasTeamsServiceMock) RenameTeam(ctx context.Context, at *teams.Team, orgID string, newName string) (*teams.Team, error) {
	ret := _m.Called(ctx, at, orgID, newName)

	if len(ret) == 0 {
		panic("no return value specified for RenameTeam")
	}

	var r0 *teams.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *teams.Team, string, string) (*teams.Team, error)); ok {
		return rf(ctx, at, orgID, newName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *teams.Team, string, string) *teams.Team); ok {
		r0 = rf(ctx, at, orgID, newName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *teams.Team, string, string) error); ok {
		r1 = rf(ctx, at, orgID, newName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AtlasTeamsServiceMock_RenameTeam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenameTeam'
type AtlasTeamsServiceMock_RenameTeam_Call struct {
	*mock.Call
}

// RenameTeam is a helper method to define mock.On call
//   - ctx context.Context
//   - at *teams.Team
//   - orgID string
//   - newName string
func (_e *AtlasTeamsServiceMock_Expecter) RenameTeam(ctx interface{}, at interface{}, orgID interface{}, newName interface{}) *AtlasTeamsServiceMock_RenameTeam_Call {
	return &AtlasTeamsServiceMock_RenameTeam_Call{Call: _e.mock.On("RenameTeam", ctx, at, orgID, newName)}
}

func (_c *AtlasTeamsServiceMock_RenameTeam_Call) Run(run func(ctx context.Context, at *teams.Team, orgID string, newName string)) *AtlasTeamsServiceMock_RenameTeam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*teams.Team), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_RenameTeam_Call) Return(_a0 *teams.Team, _a1 error) *AtlasTeamsServiceMock_RenameTeam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AtlasTeamsServiceMock_RenameTeam_Call) RunAndReturn(run func(context.Context, *teams.Team, string, string) (*teams.Team, error)) *AtlasTeamsServiceMock_RenameTeam_Call {
	_c.Call.Return(run)
	return _c
}

// Unassign provides a mock function with given fields: ctx, projectID, teamID
func (_m *AtlasTeamsServiceMock) Unassign(ctx context.Context, projectID string, teamID string) error {
	ret := _m.Called(ctx, projectID, teamID)

	if len(ret) == 0 {
		panic("no return value specified for Unassign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, projectID, teamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AtlasTeamsServiceMock_Unassign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unassign'
type AtlasTeamsServiceMock_Unassign_Call struct {
	*mock.Call
}

// Unassign is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - teamID string
func (_e *AtlasTeamsServiceMock_Expecter) Unassign(ctx interface{}, projectID interface{}, teamID interface{}) *AtlasTeamsServiceMock_Unassign_Call {
	return &AtlasTeamsServiceMock_Unassign_Call{Call: _e.mock.On("Unassign", ctx, projectID, teamID)}
}

func (_c *AtlasTeamsServiceMock_Unassign_Call) Run(run func(ctx context.Context, projectID string, teamID string)) *AtlasTeamsServiceMock_Unassign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_Unassign_Call) Return(_a0 error) *AtlasTeamsServiceMock_Unassign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AtlasTeamsServiceMock_Unassign_Call) RunAndReturn(run func(context.Context, string, string) error) *AtlasTeamsServiceMock_Unassign_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRoles provides a mock function with given fields: ctx, at, projectID, newRoles
func (_m *AtlasTeamsServiceMock) UpdateRoles(ctx context.Context, at *teams.Team, projectID string, newRoles []v1.TeamRole) error {
	ret := _m.Called(ctx, at, projectID, newRoles)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRoles")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *teams.Team, string, []v1.TeamRole) error); ok {
		r0 = rf(ctx, at, projectID, newRoles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AtlasTeamsServiceMock_UpdateRoles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRoles'
type AtlasTeamsServiceMock_UpdateRoles_Call struct {
	*mock.Call
}

// UpdateRoles is a helper method to define mock.On call
//   - ctx context.Context
//   - at *teams.Team
//   - projectID string
//   - newRoles []v1.TeamRole
func (_e *AtlasTeamsServiceMock_Expecter) UpdateRoles(ctx interface{}, at interface{}, projectID interface{}, newRoles interface{}) *AtlasTeamsServiceMock_UpdateRoles_Call {
	return &AtlasTeamsServiceMock_UpdateRoles_Call{Call: _e.mock.On("UpdateRoles", ctx, at, projectID, newRoles)}
}

func (_c *AtlasTeamsServiceMock_UpdateRoles_Call) Run(run func(ctx context.Context, at *teams.Team, projectID string, newRoles []v1.TeamRole)) *AtlasTeamsServiceMock_UpdateRoles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*teams.Team), args[2].(string), args[3].([]v1.TeamRole))
	})
	return _c
}

func (_c *AtlasTeamsServiceMock_UpdateRoles_Call) Return(_a0 error) *AtlasTeamsServiceMock_UpdateRoles_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AtlasTeamsServiceMock_UpdateRoles_Call) RunAndReturn(run func(context.Context, *teams.Team, string, []v1.TeamRole) error) *AtlasTeamsServiceMock_UpdateRoles_Call {
	_c.Call.Return(run)
	return _c
}

// NewAtlasTeamsServiceMock creates a new instance of AtlasTeamsServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAtlasTeamsServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *AtlasTeamsServiceMock {
	mock := &AtlasTeamsServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}