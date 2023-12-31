// Code generated by mockery v2.35.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entities "sportsync/entities"

	mock "github.com/stretchr/testify/mock"

	models "sportsync/models"
)

// TeamRepository is an autogenerated mock type for the TeamRepository type
type TeamRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: c, team
func (_m *TeamRepository) Create(c context.Context, team *entities.Team) error {
	ret := _m.Called(c, team)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Team) error); ok {
		r0 = rf(c, team)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: c, id
func (_m *TeamRepository) GetByID(c context.Context, id string) (entities.Team, error) {
	ret := _m.Called(c, id)

	var r0 entities.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (entities.Team, error)); ok {
		return rf(c, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) entities.Team); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Get(0).(entities.Team)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNameAndSport provides a mock function with given fields: c, name, sportName
func (_m *TeamRepository) GetByNameAndSport(c context.Context, name string, sportName string) (entities.Team, error) {
	ret := _m.Called(c, name, sportName)

	var r0 entities.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (entities.Team, error)); ok {
		return rf(c, name, sportName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) entities.Team); ok {
		r0 = rf(c, name, sportName)
	} else {
		r0 = ret.Get(0).(entities.Team)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(c, name, sportName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMyTeam provides a mock function with given fields: c, filter, userId
func (_m *TeamRepository) GetMyTeam(c context.Context, filter models.GetMyTeamBody, userId string) ([]entities.Team, models.Page, error) {
	ret := _m.Called(c, filter, userId)

	var r0 []entities.Team
	var r1 models.Page
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, models.GetMyTeamBody, string) ([]entities.Team, models.Page, error)); ok {
		return rf(c, filter, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.GetMyTeamBody, string) []entities.Team); ok {
		r0 = rf(c, filter, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.GetMyTeamBody, string) models.Page); ok {
		r1 = rf(c, filter, userId)
	} else {
		r1 = ret.Get(1).(models.Page)
	}

	if rf, ok := ret.Get(2).(func(context.Context, models.GetMyTeamBody, string) error); ok {
		r2 = rf(c, filter, userId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewTeamRepository creates a new instance of TeamRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTeamRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TeamRepository {
	mock := &TeamRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
