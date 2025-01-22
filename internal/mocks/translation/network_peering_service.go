// Code generated by mockery. DO NOT EDIT.

package translation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	networkpeering "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

// NetworkPeeringServiceMock is an autogenerated mock type for the NetworkPeeringService type
type NetworkPeeringServiceMock struct {
	mock.Mock
}

type NetworkPeeringServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *NetworkPeeringServiceMock) EXPECT() *NetworkPeeringServiceMock_Expecter {
	return &NetworkPeeringServiceMock_Expecter{mock: &_m.Mock}
}

// CreateContainer provides a mock function with given fields: ctx, projectID, container
func (_m *NetworkPeeringServiceMock) CreateContainer(ctx context.Context, projectID string, container *networkpeering.ProviderContainer) (*networkpeering.ProviderContainer, error) {
	ret := _m.Called(ctx, projectID, container)

	if len(ret) == 0 {
		panic("no return value specified for CreateContainer")
	}

	var r0 *networkpeering.ProviderContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *networkpeering.ProviderContainer) (*networkpeering.ProviderContainer, error)); ok {
		return rf(ctx, projectID, container)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *networkpeering.ProviderContainer) *networkpeering.ProviderContainer); ok {
		r0 = rf(ctx, projectID, container)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*networkpeering.ProviderContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *networkpeering.ProviderContainer) error); ok {
		r1 = rf(ctx, projectID, container)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkPeeringServiceMock_CreateContainer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateContainer'
type NetworkPeeringServiceMock_CreateContainer_Call struct {
	*mock.Call
}

// CreateContainer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - container *networkpeering.ProviderContainer
func (_e *NetworkPeeringServiceMock_Expecter) CreateContainer(ctx interface{}, projectID interface{}, container interface{}) *NetworkPeeringServiceMock_CreateContainer_Call {
	return &NetworkPeeringServiceMock_CreateContainer_Call{Call: _e.mock.On("CreateContainer", ctx, projectID, container)}
}

func (_c *NetworkPeeringServiceMock_CreateContainer_Call) Run(run func(ctx context.Context, projectID string, container *networkpeering.ProviderContainer)) *NetworkPeeringServiceMock_CreateContainer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*networkpeering.ProviderContainer))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_CreateContainer_Call) Return(_a0 *networkpeering.ProviderContainer, _a1 error) *NetworkPeeringServiceMock_CreateContainer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NetworkPeeringServiceMock_CreateContainer_Call) RunAndReturn(run func(context.Context, string, *networkpeering.ProviderContainer) (*networkpeering.ProviderContainer, error)) *NetworkPeeringServiceMock_CreateContainer_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePeer provides a mock function with given fields: ctx, projectID, conn
func (_m *NetworkPeeringServiceMock) CreatePeer(ctx context.Context, projectID string, conn *networkpeering.NetworkPeer) (*networkpeering.NetworkPeer, error) {
	ret := _m.Called(ctx, projectID, conn)

	if len(ret) == 0 {
		panic("no return value specified for CreatePeer")
	}

	var r0 *networkpeering.NetworkPeer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *networkpeering.NetworkPeer) (*networkpeering.NetworkPeer, error)); ok {
		return rf(ctx, projectID, conn)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *networkpeering.NetworkPeer) *networkpeering.NetworkPeer); ok {
		r0 = rf(ctx, projectID, conn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*networkpeering.NetworkPeer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *networkpeering.NetworkPeer) error); ok {
		r1 = rf(ctx, projectID, conn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkPeeringServiceMock_CreatePeer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePeer'
type NetworkPeeringServiceMock_CreatePeer_Call struct {
	*mock.Call
}

// CreatePeer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - conn *networkpeering.NetworkPeer
func (_e *NetworkPeeringServiceMock_Expecter) CreatePeer(ctx interface{}, projectID interface{}, conn interface{}) *NetworkPeeringServiceMock_CreatePeer_Call {
	return &NetworkPeeringServiceMock_CreatePeer_Call{Call: _e.mock.On("CreatePeer", ctx, projectID, conn)}
}

func (_c *NetworkPeeringServiceMock_CreatePeer_Call) Run(run func(ctx context.Context, projectID string, conn *networkpeering.NetworkPeer)) *NetworkPeeringServiceMock_CreatePeer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*networkpeering.NetworkPeer))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_CreatePeer_Call) Return(_a0 *networkpeering.NetworkPeer, _a1 error) *NetworkPeeringServiceMock_CreatePeer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NetworkPeeringServiceMock_CreatePeer_Call) RunAndReturn(run func(context.Context, string, *networkpeering.NetworkPeer) (*networkpeering.NetworkPeer, error)) *NetworkPeeringServiceMock_CreatePeer_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteContainer provides a mock function with given fields: ctx, projectID, containerID
func (_m *NetworkPeeringServiceMock) DeleteContainer(ctx context.Context, projectID string, containerID string) error {
	ret := _m.Called(ctx, projectID, containerID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteContainer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, projectID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NetworkPeeringServiceMock_DeleteContainer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteContainer'
type NetworkPeeringServiceMock_DeleteContainer_Call struct {
	*mock.Call
}

// DeleteContainer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - containerID string
func (_e *NetworkPeeringServiceMock_Expecter) DeleteContainer(ctx interface{}, projectID interface{}, containerID interface{}) *NetworkPeeringServiceMock_DeleteContainer_Call {
	return &NetworkPeeringServiceMock_DeleteContainer_Call{Call: _e.mock.On("DeleteContainer", ctx, projectID, containerID)}
}

func (_c *NetworkPeeringServiceMock_DeleteContainer_Call) Run(run func(ctx context.Context, projectID string, containerID string)) *NetworkPeeringServiceMock_DeleteContainer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_DeleteContainer_Call) Return(_a0 error) *NetworkPeeringServiceMock_DeleteContainer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NetworkPeeringServiceMock_DeleteContainer_Call) RunAndReturn(run func(context.Context, string, string) error) *NetworkPeeringServiceMock_DeleteContainer_Call {
	_c.Call.Return(run)
	return _c
}

// DeletePeer provides a mock function with given fields: ctx, projectID, containerID
func (_m *NetworkPeeringServiceMock) DeletePeer(ctx context.Context, projectID string, containerID string) error {
	ret := _m.Called(ctx, projectID, containerID)

	if len(ret) == 0 {
		panic("no return value specified for DeletePeer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, projectID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NetworkPeeringServiceMock_DeletePeer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeletePeer'
type NetworkPeeringServiceMock_DeletePeer_Call struct {
	*mock.Call
}

// DeletePeer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - containerID string
func (_e *NetworkPeeringServiceMock_Expecter) DeletePeer(ctx interface{}, projectID interface{}, containerID interface{}) *NetworkPeeringServiceMock_DeletePeer_Call {
	return &NetworkPeeringServiceMock_DeletePeer_Call{Call: _e.mock.On("DeletePeer", ctx, projectID, containerID)}
}

func (_c *NetworkPeeringServiceMock_DeletePeer_Call) Run(run func(ctx context.Context, projectID string, containerID string)) *NetworkPeeringServiceMock_DeletePeer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_DeletePeer_Call) Return(_a0 error) *NetworkPeeringServiceMock_DeletePeer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NetworkPeeringServiceMock_DeletePeer_Call) RunAndReturn(run func(context.Context, string, string) error) *NetworkPeeringServiceMock_DeletePeer_Call {
	_c.Call.Return(run)
	return _c
}

// FindContainer provides a mock function with given fields: ctx, projectID, provider, cidrBlock
func (_m *NetworkPeeringServiceMock) FindContainer(ctx context.Context, projectID string, provider string, cidrBlock string) (*networkpeering.ProviderContainer, error) {
	ret := _m.Called(ctx, projectID, provider, cidrBlock)

	if len(ret) == 0 {
		panic("no return value specified for FindContainer")
	}

	var r0 *networkpeering.ProviderContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (*networkpeering.ProviderContainer, error)); ok {
		return rf(ctx, projectID, provider, cidrBlock)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *networkpeering.ProviderContainer); ok {
		r0 = rf(ctx, projectID, provider, cidrBlock)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*networkpeering.ProviderContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, projectID, provider, cidrBlock)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkPeeringServiceMock_FindContainer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindContainer'
type NetworkPeeringServiceMock_FindContainer_Call struct {
	*mock.Call
}

// FindContainer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - provider string
//   - cidrBlock string
func (_e *NetworkPeeringServiceMock_Expecter) FindContainer(ctx interface{}, projectID interface{}, provider interface{}, cidrBlock interface{}) *NetworkPeeringServiceMock_FindContainer_Call {
	return &NetworkPeeringServiceMock_FindContainer_Call{Call: _e.mock.On("FindContainer", ctx, projectID, provider, cidrBlock)}
}

func (_c *NetworkPeeringServiceMock_FindContainer_Call) Run(run func(ctx context.Context, projectID string, provider string, cidrBlock string)) *NetworkPeeringServiceMock_FindContainer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_FindContainer_Call) Return(_a0 *networkpeering.ProviderContainer, _a1 error) *NetworkPeeringServiceMock_FindContainer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NetworkPeeringServiceMock_FindContainer_Call) RunAndReturn(run func(context.Context, string, string, string) (*networkpeering.ProviderContainer, error)) *NetworkPeeringServiceMock_FindContainer_Call {
	_c.Call.Return(run)
	return _c
}

// GetContainer provides a mock function with given fields: ctx, projectID, containerID
func (_m *NetworkPeeringServiceMock) GetContainer(ctx context.Context, projectID string, containerID string) (*networkpeering.ProviderContainer, error) {
	ret := _m.Called(ctx, projectID, containerID)

	if len(ret) == 0 {
		panic("no return value specified for GetContainer")
	}

	var r0 *networkpeering.ProviderContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*networkpeering.ProviderContainer, error)); ok {
		return rf(ctx, projectID, containerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *networkpeering.ProviderContainer); ok {
		r0 = rf(ctx, projectID, containerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*networkpeering.ProviderContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, projectID, containerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkPeeringServiceMock_GetContainer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetContainer'
type NetworkPeeringServiceMock_GetContainer_Call struct {
	*mock.Call
}

// GetContainer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - containerID string
func (_e *NetworkPeeringServiceMock_Expecter) GetContainer(ctx interface{}, projectID interface{}, containerID interface{}) *NetworkPeeringServiceMock_GetContainer_Call {
	return &NetworkPeeringServiceMock_GetContainer_Call{Call: _e.mock.On("GetContainer", ctx, projectID, containerID)}
}

func (_c *NetworkPeeringServiceMock_GetContainer_Call) Run(run func(ctx context.Context, projectID string, containerID string)) *NetworkPeeringServiceMock_GetContainer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_GetContainer_Call) Return(_a0 *networkpeering.ProviderContainer, _a1 error) *NetworkPeeringServiceMock_GetContainer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NetworkPeeringServiceMock_GetContainer_Call) RunAndReturn(run func(context.Context, string, string) (*networkpeering.ProviderContainer, error)) *NetworkPeeringServiceMock_GetContainer_Call {
	_c.Call.Return(run)
	return _c
}

// GetPeer provides a mock function with given fields: ctx, projectID, containerID
func (_m *NetworkPeeringServiceMock) GetPeer(ctx context.Context, projectID string, containerID string) (*networkpeering.NetworkPeer, error) {
	ret := _m.Called(ctx, projectID, containerID)

	if len(ret) == 0 {
		panic("no return value specified for GetPeer")
	}

	var r0 *networkpeering.NetworkPeer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*networkpeering.NetworkPeer, error)); ok {
		return rf(ctx, projectID, containerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *networkpeering.NetworkPeer); ok {
		r0 = rf(ctx, projectID, containerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*networkpeering.NetworkPeer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, projectID, containerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkPeeringServiceMock_GetPeer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPeer'
type NetworkPeeringServiceMock_GetPeer_Call struct {
	*mock.Call
}

// GetPeer is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - containerID string
func (_e *NetworkPeeringServiceMock_Expecter) GetPeer(ctx interface{}, projectID interface{}, containerID interface{}) *NetworkPeeringServiceMock_GetPeer_Call {
	return &NetworkPeeringServiceMock_GetPeer_Call{Call: _e.mock.On("GetPeer", ctx, projectID, containerID)}
}

func (_c *NetworkPeeringServiceMock_GetPeer_Call) Run(run func(ctx context.Context, projectID string, containerID string)) *NetworkPeeringServiceMock_GetPeer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *NetworkPeeringServiceMock_GetPeer_Call) Return(_a0 *networkpeering.NetworkPeer, _a1 error) *NetworkPeeringServiceMock_GetPeer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NetworkPeeringServiceMock_GetPeer_Call) RunAndReturn(run func(context.Context, string, string) (*networkpeering.NetworkPeer, error)) *NetworkPeeringServiceMock_GetPeer_Call {
	_c.Call.Return(run)
	return _c
}

// NewNetworkPeeringServiceMock creates a new instance of NetworkPeeringServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNetworkPeeringServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *NetworkPeeringServiceMock {
	mock := &NetworkPeeringServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
