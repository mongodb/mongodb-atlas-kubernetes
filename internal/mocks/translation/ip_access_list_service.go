// Code generated by mockery. DO NOT EDIT.

package translation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	ipaccesslist "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

// IPAccessListServiceMock is an autogenerated mock type for the IPAccessListService type
type IPAccessListServiceMock struct {
	mock.Mock
}

type IPAccessListServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *IPAccessListServiceMock) EXPECT() *IPAccessListServiceMock_Expecter {
	return &IPAccessListServiceMock_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, projectID, entries
func (_m *IPAccessListServiceMock) Add(ctx context.Context, projectID string, entries ipaccesslist.IPAccessEntries) error {
	ret := _m.Called(ctx, projectID, entries)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ipaccesslist.IPAccessEntries) error); ok {
		r0 = rf(ctx, projectID, entries)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IPAccessListServiceMock_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type IPAccessListServiceMock_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - entries ipaccesslist.IPAccessEntries
func (_e *IPAccessListServiceMock_Expecter) Add(ctx interface{}, projectID interface{}, entries interface{}) *IPAccessListServiceMock_Add_Call {
	return &IPAccessListServiceMock_Add_Call{Call: _e.mock.On("Add", ctx, projectID, entries)}
}

func (_c *IPAccessListServiceMock_Add_Call) Run(run func(ctx context.Context, projectID string, entries ipaccesslist.IPAccessEntries)) *IPAccessListServiceMock_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(ipaccesslist.IPAccessEntries))
	})
	return _c
}

func (_c *IPAccessListServiceMock_Add_Call) Return(_a0 error) *IPAccessListServiceMock_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IPAccessListServiceMock_Add_Call) RunAndReturn(run func(context.Context, string, ipaccesslist.IPAccessEntries) error) *IPAccessListServiceMock_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, projectID, entry
func (_m *IPAccessListServiceMock) Delete(ctx context.Context, projectID string, entry *ipaccesslist.IPAccessEntry) error {
	ret := _m.Called(ctx, projectID, entry)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *ipaccesslist.IPAccessEntry) error); ok {
		r0 = rf(ctx, projectID, entry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IPAccessListServiceMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type IPAccessListServiceMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - entry *ipaccesslist.IPAccessEntry
func (_e *IPAccessListServiceMock_Expecter) Delete(ctx interface{}, projectID interface{}, entry interface{}) *IPAccessListServiceMock_Delete_Call {
	return &IPAccessListServiceMock_Delete_Call{Call: _e.mock.On("Delete", ctx, projectID, entry)}
}

func (_c *IPAccessListServiceMock_Delete_Call) Run(run func(ctx context.Context, projectID string, entry *ipaccesslist.IPAccessEntry)) *IPAccessListServiceMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*ipaccesslist.IPAccessEntry))
	})
	return _c
}

func (_c *IPAccessListServiceMock_Delete_Call) Return(_a0 error) *IPAccessListServiceMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IPAccessListServiceMock_Delete_Call) RunAndReturn(run func(context.Context, string, *ipaccesslist.IPAccessEntry) error) *IPAccessListServiceMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, projectID
func (_m *IPAccessListServiceMock) List(ctx context.Context, projectID string) (ipaccesslist.IPAccessEntries, error) {
	ret := _m.Called(ctx, projectID)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 ipaccesslist.IPAccessEntries
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (ipaccesslist.IPAccessEntries, error)); ok {
		return rf(ctx, projectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) ipaccesslist.IPAccessEntries); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ipaccesslist.IPAccessEntries)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IPAccessListServiceMock_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type IPAccessListServiceMock_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
func (_e *IPAccessListServiceMock_Expecter) List(ctx interface{}, projectID interface{}) *IPAccessListServiceMock_List_Call {
	return &IPAccessListServiceMock_List_Call{Call: _e.mock.On("List", ctx, projectID)}
}

func (_c *IPAccessListServiceMock_List_Call) Run(run func(ctx context.Context, projectID string)) *IPAccessListServiceMock_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *IPAccessListServiceMock_List_Call) Return(_a0 ipaccesslist.IPAccessEntries, _a1 error) *IPAccessListServiceMock_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IPAccessListServiceMock_List_Call) RunAndReturn(run func(context.Context, string) (ipaccesslist.IPAccessEntries, error)) *IPAccessListServiceMock_List_Call {
	_c.Call.Return(run)
	return _c
}

// Status provides a mock function with given fields: ctx, projectID, entry
func (_m *IPAccessListServiceMock) Status(ctx context.Context, projectID string, entry *ipaccesslist.IPAccessEntry) (string, error) {
	ret := _m.Called(ctx, projectID, entry)

	if len(ret) == 0 {
		panic("no return value specified for Status")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *ipaccesslist.IPAccessEntry) (string, error)); ok {
		return rf(ctx, projectID, entry)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *ipaccesslist.IPAccessEntry) string); ok {
		r0 = rf(ctx, projectID, entry)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *ipaccesslist.IPAccessEntry) error); ok {
		r1 = rf(ctx, projectID, entry)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IPAccessListServiceMock_Status_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Status'
type IPAccessListServiceMock_Status_Call struct {
	*mock.Call
}

// Status is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID string
//   - entry *ipaccesslist.IPAccessEntry
func (_e *IPAccessListServiceMock_Expecter) Status(ctx interface{}, projectID interface{}, entry interface{}) *IPAccessListServiceMock_Status_Call {
	return &IPAccessListServiceMock_Status_Call{Call: _e.mock.On("Status", ctx, projectID, entry)}
}

func (_c *IPAccessListServiceMock_Status_Call) Run(run func(ctx context.Context, projectID string, entry *ipaccesslist.IPAccessEntry)) *IPAccessListServiceMock_Status_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*ipaccesslist.IPAccessEntry))
	})
	return _c
}

func (_c *IPAccessListServiceMock_Status_Call) Return(_a0 string, _a1 error) *IPAccessListServiceMock_Status_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IPAccessListServiceMock_Status_Call) RunAndReturn(run func(context.Context, string, *ipaccesslist.IPAccessEntry) (string, error)) *IPAccessListServiceMock_Status_Call {
	_c.Call.Return(run)
	return _c
}

// NewIPAccessListServiceMock creates a new instance of IPAccessListServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPAccessListServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPAccessListServiceMock {
	mock := &IPAccessListServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
