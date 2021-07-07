package ServiceTest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/nanaagyirbrown/memrizr/mocks"
	"github.com/nanaagyirbrown/memrizr/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test(t *testing.T){
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uid,_ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID: uid,
			Email: "bob@bob.com",
			Name: "Bobby Bobson",
		}

		//Instantiate
		mockUserRepository := new(mocks.MockUserRepository)
		us := service.NewUserService(&service.USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx,uid)

		assert.NoError(t,err)
		assert.Equal(t,u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		uid,_ := uuid.NewRandom()

		//Instantiate
		mockUserRepository := new(mocks.MockUserRepository)
		us := service.NewUserService(&service.USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("some error down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx,uid)

		assert.Error(t,err)
		assert.Nil(t,u)
		mockUserRepository.AssertExpectations(t)
	})
}
