package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// FindByID is mock of UserRepository FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	ret := m.Called(ctx, uid)

	var returnVal0 *model.User
	if ret.Get(0) != nil {
		returnVal0 = ret.Get(0).(*model.User)
	}

	var returnVal1 error

	if ret.Get(1) != nil {
		returnVal1 = ret.Get(1).(error)
	}

	return returnVal0, returnVal1
}

// Create is a mock for UserRepository Create
func (m *MockUserRepository) Create(ctx context.Context, u *model.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
