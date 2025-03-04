// Code generated by mockery. DO NOT EDIT.

package translation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	teams "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
)

// TeamsServiceMock is an autogenerated mock type for the TeamsService type
type TeamsServiceMock struct {
	mock.Mock
}

type TeamsServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *TeamsServiceMock) EXPECT() *TeamsServiceMock_Expecter {
	return &TeamsServiceMock_Expecter{mock: &_m.Mock}
}

// AddUsers provides a mock function with given fields: ctx, usersToAdd, orgID, teamID
func (_m *TeamsServiceMock) AddUsers(ctx context.Context, usersToAdd *[]teams.TeamUser, orgID string, teamID string) error {
	ret := _m.Called(ctx, usersToAdd, orgID, teamID)

	if len(ret) == 0 {
		panic("no return value specified for AddUsers")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *[]teams.TeamUser, string, string) error); ok {
		r0 = rf(ctx, usersToAdd, orgID, teamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TeamsServiceMock_AddUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUsers'
type TeamsServiceMock_AddUsers_Call struct {
	*mock.Call
}

// AddUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - usersToAdd *[]teams.TeamUser
//   - orgID string
//   - teamID string
func (_e *TeamsServiceMock_Expecter) AddUsers(ctx interface{}, usersToAdd interface{}, orgID interface{}, teamID interface{}) *TeamsServiceMock_AddUsers_Call {
	return &TeamsServiceMock_AddUsers_Call{Call: _e.mock.On("AddUsers", ctx, usersToAdd, orgID, teamID)}
}

func (_c *TeamsServiceMock_AddUsers_Call) Run(run func(ctx context.Context, usersToAdd *[]teams.TeamUser, orgID string, teamID string)) *TeamsServiceMock_AddUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*[]teams.TeamUser), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_AddUsers_Call) Return(_a0 error) *TeamsServiceMock_AddUsers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TeamsServiceMock_AddUsers_Call) RunAndReturn(run func(context.Context, *[]teams.TeamUser, string, string) error) *TeamsServiceMock_AddUsers_Call {
	_c.Call.Return(run)
	return _c
}

// Assign provides a mock function with given fields: ctx, at, projectID
func (_m *TeamsServiceMock) Assign(ctx context.Context, at *[]teams.AssignedTeam, projectID string) error {
	ret := _m.Called(ctx, at, projectID)

	if len(ret) == 0 {
		panic("no return value specified for Assign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *[]teams.AssignedTeam, string) error); ok {
		r0 = rf(ctx, at, projectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TeamsServiceMock_Assign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Assign'
type TeamsServiceMock_Assign_Call struct {
	*mock.Call
}

// Assign is a helper method to define mock.On call
//   - ctx context.Context
//   - at *[]teams.AssignedTeam
//   - projectID string
func (_e *TeamsServiceMock_Expecter) Assign(ctx interface{}, at interface{}, projectID interface{}) *TeamsServiceMock_Assign_Call {
	return &TeamsServiceMock_Assign_Call{Call: _e.mock.On("Assign", ctx, at, projectID)}
}

func (_c *TeamsServiceMock_Assign_Call) Run(run func(ctx context.Context, at *[]teams.AssignedTeam, projectID string)) *TeamsServiceMock_Assign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*[]teams.AssignedTeam), args[2].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_Assign_Call) Return(_a0 error) *TeamsServiceMock_Assign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TeamsServiceMock_Assign_Call) RunAndReturn(run func(context.Context, *[]teams.AssignedTeam, string) error) *TeamsServiceMock_Assign_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, at, orgID
func (_m *TeamsServiceMock) Create(ctx context.Context, at *teams.Team, orgID string) (*teams.Team, error) {
	ret := _m.Called(ctx, at, orgID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *teams.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *teams.Team, string) (*teams.Team, error)); ok {
		return rf(ctx, at, orgID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *teams.Team, string) *teams.Team); ok {
		r0 = rf(ctx, at, orgID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *teams.Team, string) error); ok {
		r1 = rf(ctx, at, orgID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamsServiceMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type TeamsServiceMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - at *teams.Team
//   - orgID string
func (_e *TeamsServiceMock_Expecter) Create(ctx interface{}, at interface{}, orgID interface{}) *TeamsServiceMock_Create_Call {
	return &TeamsServiceMock_Create_Call{Call: _e.mock.On("Create", ctx, at, orgID)}
}

func (_c *TeamsServiceMock_Create_Call) Run(run func(ctx context.Context, at *teams.Team, orgID string)) *TeamsServiceMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*teams.Team), args[2].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_Create_Call) Return(_a0 *teams.Team, _a1 error) *TeamsServiceMock_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamsServiceMock_Create_Call) RunAndReturn(run func(context.Context, *teams.Team, string) (*teams.Team, error)) *TeamsServiceMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamByID provides a mock function with given fields: ctx, orgID, teamID
func (_m *TeamsServiceMock) GetTeamByID(ctx context.Context, orgID string, teamID string) (*teams.AssignedTeam, error) {
	ret := _m.Called(ctx, orgID, teamID)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamByID")
	}

	var r0 *teams.AssignedTeam
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*teams.AssignedTeam, error)); ok {
		return rf(ctx, orgID, teamID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *teams.AssignedTeam); ok {
		r0 = rf(ctx, orgID, teamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.AssignedTeam)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, orgID, teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamsServiceMock_GetTeamByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamByID'
type TeamsServiceMock_GetTeamByID_Call struct {
	*mock.Call
}

// GetTeamByID is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamID string
func (_e *TeamsServiceMock_Expecter) GetTeamByID(ctx interface{}, orgID interface{}, teamID interface{}) *TeamsServiceMock_GetTeamByID_Call {
	return &TeamsServiceMock_GetTeamByID_Call{Call: _e.mock.On("GetTeamByID", ctx, orgID, teamID)}
}

func (_c *TeamsServiceMock_GetTeamByID_Call) Run(run func(ctx context.Context, orgID string, teamID string)) *TeamsServiceMock_GetTeamByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_GetTeamByID_Call) Return(_a0 *teams.AssignedTeam, _a1 error) *TeamsServiceMock_GetTeamByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamsServiceMock_GetTeamByID_Call) RunAndReturn(run func(context.Context, string, string) (*teams.AssignedTeam, error)) *TeamsServiceMock_GetTeamByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamByName provides a mock function with given fields: ctx, orgID, teamName
func (_m *TeamsServiceMock) GetTeamByName(ctx context.Context, orgID string, teamName string) (*teams.AssignedTeam, error) {
	ret := _m.Called(ctx, orgID, teamName)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamByName")
	}

	var r0 *teams.AssignedTeam
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*teams.AssignedTeam, error)); ok {
		return rf(ctx, orgID, teamName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *teams.AssignedTeam); ok {
		r0 = rf(ctx, orgID, teamName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.AssignedTeam)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, orgID, teamName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamsServiceMock_GetTeamByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamByName'
type TeamsServiceMock_GetTeamByName_Call struct {
	*mock.Call
}

// GetTeamByName is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamName string
func (_e *TeamsServiceMock_Expecter) GetTeamByName(ctx interface{}, orgID interface{}, teamName interface{}) *TeamsServiceMock_GetTeamByName_Call {
	return &TeamsServiceMock_GetTeamByName_Call{Call: _e.mock.On("GetTeamByName", ctx, orgID, teamName)}
}

func (_c *TeamsServiceMock_GetTeamByName_Call) Run(run func(ctx context.Context, orgID string, teamName string)) *TeamsServiceMock_GetTeamByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_GetTeamByName_Call) Return(_a0 *teams.AssignedTeam, _a1 error) *TeamsServiceMock_GetTeamByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamsServiceMock_GetTeamByName_Call) RunAndReturn(run func(context.Context, string, string) (*teams.AssignedTeam, error)) *TeamsServiceMock_GetTeamByName_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamUsers provides a mock function with given fields: ctx, orgID, teamID
func (_m *TeamsServiceMock) GetTeamUsers(ctx context.Context, orgID string, teamID string) ([]teams.TeamUser, error) {
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

// TeamsServiceMock_GetTeamUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamUsers'
type TeamsServiceMock_GetTeamUsers_Call struct {
	*mock.Call
}

// GetTeamUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamID string
func (_e *TeamsServiceMock_Expecter) GetTeamUsers(ctx interface{}, orgID interface{}, teamID interface{}) *TeamsServiceMock_GetTeamUsers_Call {
	return &TeamsServiceMock_GetTeamUsers_Call{Call: _e.mock.On("GetTeamUsers", ctx, orgID, teamID)}
}

func (_c *TeamsServiceMock_GetTeamUsers_Call) Run(run func(ctx context.Context, orgID string, teamID string)) *TeamsServiceMock_GetTeamUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_GetTeamUsers_Call) Return(_a0 []teams.TeamUser, _a1 error) *TeamsServiceMock_GetTeamUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamsServiceMock_GetTeamUsers_Call) RunAndReturn(run func(context.Context, string, string) ([]teams.TeamUser, error)) *TeamsServiceMock_GetTeamUsers_Call {
	_c.Call.Return(run)
	return _c
}

// ListProjectTeams provides a mock function with given fields: ctx, projectID
func (_m *TeamsServiceMock) ListProjectTeams(ctx context.Context, projectID string) ([]teams.AssignedTeam, error) {
	ret := _m.Called(ctx, projectID)

	if len(ret) == 0 {
		panic("no return value specified for ListProjectTeams")
	}

	var r0 []teams.AssignedTeam
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]teams.AssignedTeam, error)); ok {
		return rf(ctx, projectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []teams.AssignedTeam); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]teams.AssignedTeam)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamsServiceMock_ListProjectTeams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListProjectTeams'
type TeamsServiceMock_ListProjectTeams_Call struct {
	*mock.Call
}

// ListProjectTeams is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
func (_e *TeamsServiceMock_Expecter) ListProjectTeams(ctx interface{}, projectID interface{}) *TeamsServiceMock_ListProjectTeams_Call {
	return &TeamsServiceMock_ListProjectTeams_Call{Call: _e.mock.On("ListProjectTeams", ctx, projectID)}
}

func (_c *TeamsServiceMock_ListProjectTeams_Call) Run(run func(ctx context.Context, projectID string)) *TeamsServiceMock_ListProjectTeams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_ListProjectTeams_Call) Return(_a0 []teams.AssignedTeam, _a1 error) *TeamsServiceMock_ListProjectTeams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamsServiceMock_ListProjectTeams_Call) RunAndReturn(run func(context.Context, string) ([]teams.AssignedTeam, error)) *TeamsServiceMock_ListProjectTeams_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveUser provides a mock function with given fields: ctx, orgID, teamID, userID
func (_m *TeamsServiceMock) RemoveUser(ctx context.Context, orgID string, teamID string, userID string) error {
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

// TeamsServiceMock_RemoveUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveUser'
type TeamsServiceMock_RemoveUser_Call struct {
	*mock.Call
}

// RemoveUser is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - teamID string
//   - userID string
func (_e *TeamsServiceMock_Expecter) RemoveUser(ctx interface{}, orgID interface{}, teamID interface{}, userID interface{}) *TeamsServiceMock_RemoveUser_Call {
	return &TeamsServiceMock_RemoveUser_Call{Call: _e.mock.On("RemoveUser", ctx, orgID, teamID, userID)}
}

func (_c *TeamsServiceMock_RemoveUser_Call) Run(run func(ctx context.Context, orgID string, teamID string, userID string)) *TeamsServiceMock_RemoveUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_RemoveUser_Call) Return(_a0 error) *TeamsServiceMock_RemoveUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TeamsServiceMock_RemoveUser_Call) RunAndReturn(run func(context.Context, string, string, string) error) *TeamsServiceMock_RemoveUser_Call {
	_c.Call.Return(run)
	return _c
}

// RenameTeam provides a mock function with given fields: ctx, at, orgID, newName
func (_m *TeamsServiceMock) RenameTeam(ctx context.Context, at *teams.AssignedTeam, orgID string, newName string) (*teams.AssignedTeam, error) {
	ret := _m.Called(ctx, at, orgID, newName)

	if len(ret) == 0 {
		panic("no return value specified for RenameTeam")
	}

	var r0 *teams.AssignedTeam
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *teams.AssignedTeam, string, string) (*teams.AssignedTeam, error)); ok {
		return rf(ctx, at, orgID, newName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *teams.AssignedTeam, string, string) *teams.AssignedTeam); ok {
		r0 = rf(ctx, at, orgID, newName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*teams.AssignedTeam)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *teams.AssignedTeam, string, string) error); ok {
		r1 = rf(ctx, at, orgID, newName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamsServiceMock_RenameTeam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenameTeam'
type TeamsServiceMock_RenameTeam_Call struct {
	*mock.Call
}

// RenameTeam is a helper method to define mock.On call
//   - ctx context.Context
//   - at *teams.AssignedTeam
//   - orgID string
//   - newName string
func (_e *TeamsServiceMock_Expecter) RenameTeam(ctx interface{}, at interface{}, orgID interface{}, newName interface{}) *TeamsServiceMock_RenameTeam_Call {
	return &TeamsServiceMock_RenameTeam_Call{Call: _e.mock.On("RenameTeam", ctx, at, orgID, newName)}
}

func (_c *TeamsServiceMock_RenameTeam_Call) Run(run func(ctx context.Context, at *teams.AssignedTeam, orgID string, newName string)) *TeamsServiceMock_RenameTeam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*teams.AssignedTeam), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_RenameTeam_Call) Return(_a0 *teams.AssignedTeam, _a1 error) *TeamsServiceMock_RenameTeam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamsServiceMock_RenameTeam_Call) RunAndReturn(run func(context.Context, *teams.AssignedTeam, string, string) (*teams.AssignedTeam, error)) *TeamsServiceMock_RenameTeam_Call {
	_c.Call.Return(run)
	return _c
}

// Unassign provides a mock function with given fields: ctx, projectID, teamID
func (_m *TeamsServiceMock) Unassign(ctx context.Context, projectID string, teamID string) error {
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

// TeamsServiceMock_Unassign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unassign'
type TeamsServiceMock_Unassign_Call struct {
	*mock.Call
}

// Unassign is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - teamID string
func (_e *TeamsServiceMock_Expecter) Unassign(ctx interface{}, projectID interface{}, teamID interface{}) *TeamsServiceMock_Unassign_Call {
	return &TeamsServiceMock_Unassign_Call{Call: _e.mock.On("Unassign", ctx, projectID, teamID)}
}

func (_c *TeamsServiceMock_Unassign_Call) Run(run func(ctx context.Context, projectID string, teamID string)) *TeamsServiceMock_Unassign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *TeamsServiceMock_Unassign_Call) Return(_a0 error) *TeamsServiceMock_Unassign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TeamsServiceMock_Unassign_Call) RunAndReturn(run func(context.Context, string, string) error) *TeamsServiceMock_Unassign_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRoles provides a mock function with given fields: ctx, at, projectID, newRoles
func (_m *TeamsServiceMock) UpdateRoles(ctx context.Context, at *teams.AssignedTeam, projectID string, newRoles []v1.TeamRole) error {
	ret := _m.Called(ctx, at, projectID, newRoles)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRoles")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *teams.AssignedTeam, string, []v1.TeamRole) error); ok {
		r0 = rf(ctx, at, projectID, newRoles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TeamsServiceMock_UpdateRoles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRoles'
type TeamsServiceMock_UpdateRoles_Call struct {
	*mock.Call
}

// UpdateRoles is a helper method to define mock.On call
//   - ctx context.Context
//   - at *teams.AssignedTeam
//   - projectID string
//   - newRoles []v1.TeamRole
func (_e *TeamsServiceMock_Expecter) UpdateRoles(ctx interface{}, at interface{}, projectID interface{}, newRoles interface{}) *TeamsServiceMock_UpdateRoles_Call {
	return &TeamsServiceMock_UpdateRoles_Call{Call: _e.mock.On("UpdateRoles", ctx, at, projectID, newRoles)}
}

func (_c *TeamsServiceMock_UpdateRoles_Call) Run(run func(ctx context.Context, at *teams.AssignedTeam, projectID string, newRoles []v1.TeamRole)) *TeamsServiceMock_UpdateRoles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*teams.AssignedTeam), args[2].(string), args[3].([]v1.TeamRole))
	})
	return _c
}

func (_c *TeamsServiceMock_UpdateRoles_Call) Return(_a0 error) *TeamsServiceMock_UpdateRoles_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TeamsServiceMock_UpdateRoles_Call) RunAndReturn(run func(context.Context, *teams.AssignedTeam, string, []v1.TeamRole) error) *TeamsServiceMock_UpdateRoles_Call {
	_c.Call.Return(run)
	return _c
}

// NewTeamsServiceMock creates a new instance of TeamsServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTeamsServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TeamsServiceMock {
	mock := &TeamsServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
