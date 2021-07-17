package service

import (
	"context"
	"github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"log"

	"github.com/google/uuid"
	"github.com/nanaagyirbrown/memrizr/handler/model"
)

// UserService acts as a struct for holding dependencies of this service, likewise injecting an implementation of UserRepository
// for use in service methods
type UserService struct {
	UserRepository model.UserRepository
}

// USConfig (User-Service Configuration) will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository model.UserRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &UserService{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *UserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

// Signup reaches our UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *UserService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to signup user email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	u.Password = pw
	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishUserUpdated(u, true)

	// if err != nil {
	//  return nil, apperrors.NewInternal()
	// }

	return nil
}
