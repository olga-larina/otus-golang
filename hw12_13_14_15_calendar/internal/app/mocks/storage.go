// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	storage "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

type Storage_Expecter struct {
	mock *mock.Mock
}

func (_m *Storage) EXPECT() *Storage_Expecter {
	return &Storage_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, event
func (_m *Storage) Create(ctx context.Context, event *storage.Event) (uint64, error) {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.Event) (uint64, error)); ok {
		return rf(ctx, event)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *storage.Event) uint64); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *storage.Event) error); ok {
		r1 = rf(ctx, event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storage_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Storage_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - event *storage.Event
func (_e *Storage_Expecter) Create(ctx interface{}, event interface{}) *Storage_Create_Call {
	return &Storage_Create_Call{Call: _e.mock.On("Create", ctx, event)}
}

func (_c *Storage_Create_Call) Run(run func(ctx context.Context, event *storage.Event)) *Storage_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*storage.Event))
	})
	return _c
}

func (_c *Storage_Create_Call) Return(_a0 uint64, _a1 error) *Storage_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Storage_Create_Call) RunAndReturn(run func(context.Context, *storage.Event) (uint64, error)) *Storage_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, userID, eventID
func (_m *Storage) Delete(ctx context.Context, userID uint64, eventID uint64) error {
	ret := _m.Called(ctx, userID, eventID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) error); ok {
		r0 = rf(ctx, userID, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storage_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type Storage_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - eventID uint64
func (_e *Storage_Expecter) Delete(ctx interface{}, userID interface{}, eventID interface{}) *Storage_Delete_Call {
	return &Storage_Delete_Call{Call: _e.mock.On("Delete", ctx, userID, eventID)}
}

func (_c *Storage_Delete_Call) Run(run func(ctx context.Context, userID uint64, eventID uint64)) *Storage_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64))
	})
	return _c
}

func (_c *Storage_Delete_Call) Return(_a0 error) *Storage_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Storage_Delete_Call) RunAndReturn(run func(context.Context, uint64, uint64) error) *Storage_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, userID, eventID
func (_m *Storage) GetByID(ctx context.Context, userID uint64, eventID uint64) (*storage.Event, error) {
	ret := _m.Called(ctx, userID, eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *storage.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) (*storage.Event, error)); ok {
		return rf(ctx, userID, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) *storage.Event); ok {
		r0 = rf(ctx, userID, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64) error); ok {
		r1 = rf(ctx, userID, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storage_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type Storage_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - eventID uint64
func (_e *Storage_Expecter) GetByID(ctx interface{}, userID interface{}, eventID interface{}) *Storage_GetByID_Call {
	return &Storage_GetByID_Call{Call: _e.mock.On("GetByID", ctx, userID, eventID)}
}

func (_c *Storage_GetByID_Call) Run(run func(ctx context.Context, userID uint64, eventID uint64)) *Storage_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64))
	})
	return _c
}

func (_c *Storage_GetByID_Call) Return(_a0 *storage.Event, _a1 error) *Storage_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Storage_GetByID_Call) RunAndReturn(run func(context.Context, uint64, uint64) (*storage.Event, error)) *Storage_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// ListForPeriod provides a mock function with given fields: ctx, userID, startDate, endDateExclusive
func (_m *Storage) ListForPeriod(ctx context.Context, userID uint64, startDate time.Time, endDateExclusive time.Time) ([]*storage.Event, error) {
	ret := _m.Called(ctx, userID, startDate, endDateExclusive)

	if len(ret) == 0 {
		panic("no return value specified for ListForPeriod")
	}

	var r0 []*storage.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time, time.Time) ([]*storage.Event, error)); ok {
		return rf(ctx, userID, startDate, endDateExclusive)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time, time.Time) []*storage.Event); ok {
		r0 = rf(ctx, userID, startDate, endDateExclusive)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*storage.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, time.Time, time.Time) error); ok {
		r1 = rf(ctx, userID, startDate, endDateExclusive)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storage_ListForPeriod_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListForPeriod'
type Storage_ListForPeriod_Call struct {
	*mock.Call
}

// ListForPeriod is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - startDate time.Time
//   - endDateExclusive time.Time
func (_e *Storage_Expecter) ListForPeriod(ctx interface{}, userID interface{}, startDate interface{}, endDateExclusive interface{}) *Storage_ListForPeriod_Call {
	return &Storage_ListForPeriod_Call{Call: _e.mock.On("ListForPeriod", ctx, userID, startDate, endDateExclusive)}
}

func (_c *Storage_ListForPeriod_Call) Run(run func(ctx context.Context, userID uint64, startDate time.Time, endDateExclusive time.Time)) *Storage_ListForPeriod_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(time.Time), args[3].(time.Time))
	})
	return _c
}

func (_c *Storage_ListForPeriod_Call) Return(_a0 []*storage.Event, _a1 error) *Storage_ListForPeriod_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Storage_ListForPeriod_Call) RunAndReturn(run func(context.Context, uint64, time.Time, time.Time) ([]*storage.Event, error)) *Storage_ListForPeriod_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, event
func (_m *Storage) Update(ctx context.Context, event *storage.Event) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.Event) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storage_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type Storage_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - event *storage.Event
func (_e *Storage_Expecter) Update(ctx interface{}, event interface{}) *Storage_Update_Call {
	return &Storage_Update_Call{Call: _e.mock.On("Update", ctx, event)}
}

func (_c *Storage_Update_Call) Run(run func(ctx context.Context, event *storage.Event)) *Storage_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*storage.Event))
	})
	return _c
}

func (_c *Storage_Update_Call) Return(_a0 error) *Storage_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Storage_Update_Call) RunAndReturn(run func(context.Context, *storage.Event) error) *Storage_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}