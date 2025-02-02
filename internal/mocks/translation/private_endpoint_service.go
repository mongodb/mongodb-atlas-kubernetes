// Code generated by mockery. DO NOT EDIT.

package translation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	privateendpoint "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
)

// PrivateEndpointServiceMock is an autogenerated mock type for the PrivateEndpointService type
type PrivateEndpointServiceMock struct {
	mock.Mock
}

type PrivateEndpointServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *PrivateEndpointServiceMock) EXPECT() *PrivateEndpointServiceMock_Expecter {
	return &PrivateEndpointServiceMock_Expecter{mock: &_m.Mock}
}

// CreatePrivateEndpointInterface provides a mock function with given fields: ctx, projectID, provider, serviceID, gcpProjectID, peInterface
func (_m *PrivateEndpointServiceMock) CreatePrivateEndpointInterface(ctx context.Context, projectID string, provider string, serviceID string, gcpProjectID string, peInterface privateendpoint.EndpointInterface) (privateendpoint.EndpointInterface, error) {
	ret := _m.Called(ctx, projectID, provider, serviceID, gcpProjectID, peInterface)

	if len(ret) == 0 {
		panic("no return value specified for CreatePrivateEndpointInterface")
	}

	var r0 privateendpoint.EndpointInterface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, privateendpoint.EndpointInterface) (privateendpoint.EndpointInterface, error)); ok {
		return rf(ctx, projectID, provider, serviceID, gcpProjectID, peInterface)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, privateendpoint.EndpointInterface) privateendpoint.EndpointInterface); ok {
		r0 = rf(ctx, projectID, provider, serviceID, gcpProjectID, peInterface)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privateendpoint.EndpointInterface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, privateendpoint.EndpointInterface) error); ok {
		r1 = rf(ctx, projectID, provider, serviceID, gcpProjectID, peInterface)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePrivateEndpointInterface'
type PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call struct {
	*mock.Call
}

// CreatePrivateEndpointInterface is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - provider string
//   - serviceID string
//   - gcpProjectID string
//   - peInterface privateendpoint.EndpointInterface
func (_e *PrivateEndpointServiceMock_Expecter) CreatePrivateEndpointInterface(ctx interface{}, projectID interface{}, provider interface{}, serviceID interface{}, gcpProjectID interface{}, peInterface interface{}) *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call {
	return &PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call{Call: _e.mock.On("CreatePrivateEndpointInterface", ctx, projectID, provider, serviceID, gcpProjectID, peInterface)}
}

func (_c *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call) Run(run func(ctx context.Context, projectID string, provider string, serviceID string, gcpProjectID string, peInterface privateendpoint.EndpointInterface)) *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(string), args[5].(privateendpoint.EndpointInterface))
	})
	return _c
}

func (_c *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call) Return(_a0 privateendpoint.EndpointInterface, _a1 error) *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call) RunAndReturn(run func(context.Context, string, string, string, string, privateendpoint.EndpointInterface) (privateendpoint.EndpointInterface, error)) *PrivateEndpointServiceMock_CreatePrivateEndpointInterface_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePrivateEndpointService provides a mock function with given fields: ctx, projectID, peService
func (_m *PrivateEndpointServiceMock) CreatePrivateEndpointService(ctx context.Context, projectID string, peService privateendpoint.EndpointService) (privateendpoint.EndpointService, error) {
	ret := _m.Called(ctx, projectID, peService)

	if len(ret) == 0 {
		panic("no return value specified for CreatePrivateEndpointService")
	}

	var r0 privateendpoint.EndpointService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, privateendpoint.EndpointService) (privateendpoint.EndpointService, error)); ok {
		return rf(ctx, projectID, peService)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, privateendpoint.EndpointService) privateendpoint.EndpointService); ok {
		r0 = rf(ctx, projectID, peService)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privateendpoint.EndpointService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, privateendpoint.EndpointService) error); ok {
		r1 = rf(ctx, projectID, peService)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrivateEndpointServiceMock_CreatePrivateEndpointService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePrivateEndpointService'
type PrivateEndpointServiceMock_CreatePrivateEndpointService_Call struct {
	*mock.Call
}

// CreatePrivateEndpointService is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - peService privateendpoint.EndpointService
func (_e *PrivateEndpointServiceMock_Expecter) CreatePrivateEndpointService(ctx interface{}, projectID interface{}, peService interface{}) *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call {
	return &PrivateEndpointServiceMock_CreatePrivateEndpointService_Call{Call: _e.mock.On("CreatePrivateEndpointService", ctx, projectID, peService)}
}

func (_c *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call) Run(run func(ctx context.Context, projectID string, peService privateendpoint.EndpointService)) *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(privateendpoint.EndpointService))
	})
	return _c
}

func (_c *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call) Return(_a0 privateendpoint.EndpointService, _a1 error) *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call) RunAndReturn(run func(context.Context, string, privateendpoint.EndpointService) (privateendpoint.EndpointService, error)) *PrivateEndpointServiceMock_CreatePrivateEndpointService_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteEndpointInterface provides a mock function with given fields: ctx, projectID, provider, serviceID, ID
func (_m *PrivateEndpointServiceMock) DeleteEndpointInterface(ctx context.Context, projectID string, provider string, serviceID string, ID string) error {
	ret := _m.Called(ctx, projectID, provider, serviceID, ID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEndpointInterface")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) error); ok {
		r0 = rf(ctx, projectID, provider, serviceID, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PrivateEndpointServiceMock_DeleteEndpointInterface_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteEndpointInterface'
type PrivateEndpointServiceMock_DeleteEndpointInterface_Call struct {
	*mock.Call
}

// DeleteEndpointInterface is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - provider string
//   - serviceID string
//   - ID string
func (_e *PrivateEndpointServiceMock_Expecter) DeleteEndpointInterface(ctx interface{}, projectID interface{}, provider interface{}, serviceID interface{}, ID interface{}) *PrivateEndpointServiceMock_DeleteEndpointInterface_Call {
	return &PrivateEndpointServiceMock_DeleteEndpointInterface_Call{Call: _e.mock.On("DeleteEndpointInterface", ctx, projectID, provider, serviceID, ID)}
}

func (_c *PrivateEndpointServiceMock_DeleteEndpointInterface_Call) Run(run func(ctx context.Context, projectID string, provider string, serviceID string, ID string)) *PrivateEndpointServiceMock_DeleteEndpointInterface_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(string))
	})
	return _c
}

func (_c *PrivateEndpointServiceMock_DeleteEndpointInterface_Call) Return(_a0 error) *PrivateEndpointServiceMock_DeleteEndpointInterface_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PrivateEndpointServiceMock_DeleteEndpointInterface_Call) RunAndReturn(run func(context.Context, string, string, string, string) error) *PrivateEndpointServiceMock_DeleteEndpointInterface_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteEndpointService provides a mock function with given fields: ctx, projectID, provider, ID
func (_m *PrivateEndpointServiceMock) DeleteEndpointService(ctx context.Context, projectID string, provider string, ID string) error {
	ret := _m.Called(ctx, projectID, provider, ID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEndpointService")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, projectID, provider, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PrivateEndpointServiceMock_DeleteEndpointService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteEndpointService'
type PrivateEndpointServiceMock_DeleteEndpointService_Call struct {
	*mock.Call
}

// DeleteEndpointService is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - provider string
//   - ID string
func (_e *PrivateEndpointServiceMock_Expecter) DeleteEndpointService(ctx interface{}, projectID interface{}, provider interface{}, ID interface{}) *PrivateEndpointServiceMock_DeleteEndpointService_Call {
	return &PrivateEndpointServiceMock_DeleteEndpointService_Call{Call: _e.mock.On("DeleteEndpointService", ctx, projectID, provider, ID)}
}

func (_c *PrivateEndpointServiceMock_DeleteEndpointService_Call) Run(run func(ctx context.Context, projectID string, provider string, ID string)) *PrivateEndpointServiceMock_DeleteEndpointService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *PrivateEndpointServiceMock_DeleteEndpointService_Call) Return(_a0 error) *PrivateEndpointServiceMock_DeleteEndpointService_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PrivateEndpointServiceMock_DeleteEndpointService_Call) RunAndReturn(run func(context.Context, string, string, string) error) *PrivateEndpointServiceMock_DeleteEndpointService_Call {
	_c.Call.Return(run)
	return _c
}

// GetPrivateEndpoint provides a mock function with given fields: ctx, projectID, provider, ID
func (_m *PrivateEndpointServiceMock) GetPrivateEndpoint(ctx context.Context, projectID string, provider string, ID string) (privateendpoint.EndpointService, error) {
	ret := _m.Called(ctx, projectID, provider, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetPrivateEndpoint")
	}

	var r0 privateendpoint.EndpointService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (privateendpoint.EndpointService, error)); ok {
		return rf(ctx, projectID, provider, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) privateendpoint.EndpointService); ok {
		r0 = rf(ctx, projectID, provider, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privateendpoint.EndpointService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, projectID, provider, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrivateEndpointServiceMock_GetPrivateEndpoint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPrivateEndpoint'
type PrivateEndpointServiceMock_GetPrivateEndpoint_Call struct {
	*mock.Call
}

// GetPrivateEndpoint is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - provider string
//   - ID string
func (_e *PrivateEndpointServiceMock_Expecter) GetPrivateEndpoint(ctx interface{}, projectID interface{}, provider interface{}, ID interface{}) *PrivateEndpointServiceMock_GetPrivateEndpoint_Call {
	return &PrivateEndpointServiceMock_GetPrivateEndpoint_Call{Call: _e.mock.On("GetPrivateEndpoint", ctx, projectID, provider, ID)}
}

func (_c *PrivateEndpointServiceMock_GetPrivateEndpoint_Call) Run(run func(ctx context.Context, projectID string, provider string, ID string)) *PrivateEndpointServiceMock_GetPrivateEndpoint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *PrivateEndpointServiceMock_GetPrivateEndpoint_Call) Return(_a0 privateendpoint.EndpointService, _a1 error) *PrivateEndpointServiceMock_GetPrivateEndpoint_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PrivateEndpointServiceMock_GetPrivateEndpoint_Call) RunAndReturn(run func(context.Context, string, string, string) (privateendpoint.EndpointService, error)) *PrivateEndpointServiceMock_GetPrivateEndpoint_Call {
	_c.Call.Return(run)
	return _c
}

// ListPrivateEndpoints provides a mock function with given fields: ctx, projectID, provider
func (_m *PrivateEndpointServiceMock) ListPrivateEndpoints(ctx context.Context, projectID string, provider string) ([]privateendpoint.EndpointService, error) {
	ret := _m.Called(ctx, projectID, provider)

	if len(ret) == 0 {
		panic("no return value specified for ListPrivateEndpoints")
	}

	var r0 []privateendpoint.EndpointService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]privateendpoint.EndpointService, error)); ok {
		return rf(ctx, projectID, provider)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []privateendpoint.EndpointService); ok {
		r0 = rf(ctx, projectID, provider)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]privateendpoint.EndpointService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, projectID, provider)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrivateEndpointServiceMock_ListPrivateEndpoints_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListPrivateEndpoints'
type PrivateEndpointServiceMock_ListPrivateEndpoints_Call struct {
	*mock.Call
}

// ListPrivateEndpoints is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - provider string
func (_e *PrivateEndpointServiceMock_Expecter) ListPrivateEndpoints(ctx interface{}, projectID interface{}, provider interface{}) *PrivateEndpointServiceMock_ListPrivateEndpoints_Call {
	return &PrivateEndpointServiceMock_ListPrivateEndpoints_Call{Call: _e.mock.On("ListPrivateEndpoints", ctx, projectID, provider)}
}

func (_c *PrivateEndpointServiceMock_ListPrivateEndpoints_Call) Run(run func(ctx context.Context, projectID string, provider string)) *PrivateEndpointServiceMock_ListPrivateEndpoints_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *PrivateEndpointServiceMock_ListPrivateEndpoints_Call) Return(_a0 []privateendpoint.EndpointService, _a1 error) *PrivateEndpointServiceMock_ListPrivateEndpoints_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PrivateEndpointServiceMock_ListPrivateEndpoints_Call) RunAndReturn(run func(context.Context, string, string) ([]privateendpoint.EndpointService, error)) *PrivateEndpointServiceMock_ListPrivateEndpoints_Call {
	_c.Call.Return(run)
	return _c
}

// NewPrivateEndpointServiceMock creates a new instance of PrivateEndpointServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPrivateEndpointServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *PrivateEndpointServiceMock {
	mock := &PrivateEndpointServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
