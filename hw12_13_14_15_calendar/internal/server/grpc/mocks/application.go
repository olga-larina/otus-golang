// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	app "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

type Application_Expecter struct {
	mock *mock.Mock
}

func (_m *Application) EXPECT() *Application_Expecter {
	return &Application_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, eventDto
func (_m *Application) Create(ctx context.Context, eventDto app.EventDto) (uint64, error) {
	ret := _m.Called(ctx, eventDto)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, app.EventDto) (uint64, error)); ok {
		return rf(ctx, eventDto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, app.EventDto) uint64); ok {
		r0 = rf(ctx, eventDto)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, app.EventDto) error); ok {
		r1 = rf(ctx, eventDto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Application_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Application_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - eventDto app.EventDto
func (_e *Application_Expecter) Create(ctx interface{}, eventDto interface{}) *Application_Create_Call {
	return &Application_Create_Call{Call: _e.mock.On("Create", ctx, eventDto)}
}

func (_c *Application_Create_Call) Run(run func(ctx context.Context, eventDto app.EventDto)) *Application_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.EventDto))
	})
	return _c
}

func (_c *Application_Create_Call) Return(_a0 uint64, _a1 error) *Application_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Application_Create_Call) RunAndReturn(run func(context.Context, app.EventDto) (uint64, error)) *Application_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, userID, eventID
func (_m *Application) Delete(ctx context.Context, userID uint64, eventID uint64) error {
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

// Application_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type Application_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - eventID uint64
func (_e *Application_Expecter) Delete(ctx interface{}, userID interface{}, eventID interface{}) *Application_Delete_Call {
	return &Application_Delete_Call{Call: _e.mock.On("Delete", ctx, userID, eventID)}
}

func (_c *Application_Delete_Call) Run(run func(ctx context.Context, userID uint64, eventID uint64)) *Application_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64))
	})
	return _c
}

func (_c *Application_Delete_Call) Return(_a0 error) *Application_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_Delete_Call) RunAndReturn(run func(context.Context, uint64, uint64) error) *Application_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, userID, eventID
func (_m *Application) GetByID(ctx context.Context, userID uint64, eventID uint64) (*app.EventDto, error) {
	ret := _m.Called(ctx, userID, eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *app.EventDto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) (*app.EventDto, error)); ok {
		return rf(ctx, userID, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) *app.EventDto); ok {
		r0 = rf(ctx, userID, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*app.EventDto)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64) error); ok {
		r1 = rf(ctx, userID, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Application_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type Application_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - eventID uint64
func (_e *Application_Expecter) GetByID(ctx interface{}, userID interface{}, eventID interface{}) *Application_GetByID_Call {
	return &Application_GetByID_Call{Call: _e.mock.On("GetByID", ctx, userID, eventID)}
}

func (_c *Application_GetByID_Call) Run(run func(ctx context.Context, userID uint64, eventID uint64)) *Application_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64))
	})
	return _c
}

func (_c *Application_GetByID_Call) Return(_a0 *app.EventDto, _a1 error) *Application_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Application_GetByID_Call) RunAndReturn(run func(context.Context, uint64, uint64) (*app.EventDto, error)) *Application_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// ListForDay provides a mock function with given fields: ctx, userID, date
func (_m *Application) ListForDay(ctx context.Context, userID uint64, date time.Time) ([]*app.EventDto, error) {
	ret := _m.Called(ctx, userID, date)

	if len(ret) == 0 {
		panic("no return value specified for ListForDay")
	}

	var r0 []*app.EventDto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time) ([]*app.EventDto, error)); ok {
		return rf(ctx, userID, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time) []*app.EventDto); ok {
		r0 = rf(ctx, userID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*app.EventDto)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, time.Time) error); ok {
		r1 = rf(ctx, userID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Application_ListForDay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListForDay'
type Application_ListForDay_Call struct {
	*mock.Call
}

// ListForDay is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - date time.Time
func (_e *Application_Expecter) ListForDay(ctx interface{}, userID interface{}, date interface{}) *Application_ListForDay_Call {
	return &Application_ListForDay_Call{Call: _e.mock.On("ListForDay", ctx, userID, date)}
}

func (_c *Application_ListForDay_Call) Run(run func(ctx context.Context, userID uint64, date time.Time)) *Application_ListForDay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(time.Time))
	})
	return _c
}

func (_c *Application_ListForDay_Call) Return(_a0 []*app.EventDto, _a1 error) *Application_ListForDay_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Application_ListForDay_Call) RunAndReturn(run func(context.Context, uint64, time.Time) ([]*app.EventDto, error)) *Application_ListForDay_Call {
	_c.Call.Return(run)
	return _c
}

// ListForMonth provides a mock function with given fields: ctx, userID, startDate
func (_m *Application) ListForMonth(ctx context.Context, userID uint64, startDate time.Time) ([]*app.EventDto, error) {
	ret := _m.Called(ctx, userID, startDate)

	if len(ret) == 0 {
		panic("no return value specified for ListForMonth")
	}

	var r0 []*app.EventDto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time) ([]*app.EventDto, error)); ok {
		return rf(ctx, userID, startDate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time) []*app.EventDto); ok {
		r0 = rf(ctx, userID, startDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*app.EventDto)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, time.Time) error); ok {
		r1 = rf(ctx, userID, startDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Application_ListForMonth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListForMonth'
type Application_ListForMonth_Call struct {
	*mock.Call
}

// ListForMonth is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - startDate time.Time
func (_e *Application_Expecter) ListForMonth(ctx interface{}, userID interface{}, startDate interface{}) *Application_ListForMonth_Call {
	return &Application_ListForMonth_Call{Call: _e.mock.On("ListForMonth", ctx, userID, startDate)}
}

func (_c *Application_ListForMonth_Call) Run(run func(ctx context.Context, userID uint64, startDate time.Time)) *Application_ListForMonth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(time.Time))
	})
	return _c
}

func (_c *Application_ListForMonth_Call) Return(_a0 []*app.EventDto, _a1 error) *Application_ListForMonth_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Application_ListForMonth_Call) RunAndReturn(run func(context.Context, uint64, time.Time) ([]*app.EventDto, error)) *Application_ListForMonth_Call {
	_c.Call.Return(run)
	return _c
}

// ListForWeek provides a mock function with given fields: ctx, userID, startDate
func (_m *Application) ListForWeek(ctx context.Context, userID uint64, startDate time.Time) ([]*app.EventDto, error) {
	ret := _m.Called(ctx, userID, startDate)

	if len(ret) == 0 {
		panic("no return value specified for ListForWeek")
	}

	var r0 []*app.EventDto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time) ([]*app.EventDto, error)); ok {
		return rf(ctx, userID, startDate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, time.Time) []*app.EventDto); ok {
		r0 = rf(ctx, userID, startDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*app.EventDto)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, time.Time) error); ok {
		r1 = rf(ctx, userID, startDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Application_ListForWeek_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListForWeek'
type Application_ListForWeek_Call struct {
	*mock.Call
}

// ListForWeek is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint64
//   - startDate time.Time
func (_e *Application_Expecter) ListForWeek(ctx interface{}, userID interface{}, startDate interface{}) *Application_ListForWeek_Call {
	return &Application_ListForWeek_Call{Call: _e.mock.On("ListForWeek", ctx, userID, startDate)}
}

func (_c *Application_ListForWeek_Call) Run(run func(ctx context.Context, userID uint64, startDate time.Time)) *Application_ListForWeek_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(time.Time))
	})
	return _c
}

func (_c *Application_ListForWeek_Call) Return(_a0 []*app.EventDto, _a1 error) *Application_ListForWeek_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Application_ListForWeek_Call) RunAndReturn(run func(context.Context, uint64, time.Time) ([]*app.EventDto, error)) *Application_ListForWeek_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, eventDto
func (_m *Application) Update(ctx context.Context, eventDto app.EventDto) error {
	ret := _m.Called(ctx, eventDto)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, app.EventDto) error); ok {
		r0 = rf(ctx, eventDto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type Application_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - eventDto app.EventDto
func (_e *Application_Expecter) Update(ctx interface{}, eventDto interface{}) *Application_Update_Call {
	return &Application_Update_Call{Call: _e.mock.On("Update", ctx, eventDto)}
}

func (_c *Application_Update_Call) Run(run func(ctx context.Context, eventDto app.EventDto)) *Application_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.EventDto))
	})
	return _c
}

func (_c *Application_Update_Call) Return(_a0 error) *Application_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_Update_Call) RunAndReturn(run func(context.Context, app.EventDto) error) *Application_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
