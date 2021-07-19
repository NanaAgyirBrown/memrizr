package handlerTest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanaagyirbrown/memrizr/handler"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"github.com/nanaagyirbrown/memrizr/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMe(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.AnythingOfType("*context.emptyCtx"), uid).Return(mockUserResp, nil)

		// a response recorder for getting written http response
		responRecorder := httptest.NewRecorder()

		// use a middleware to set context for test
		// the only claims we care about in this test
		// is the UID
		router := gin.Default()

		// Some sort of dummy middleware
		router.Use(func(c *gin.Context) {
			c.Set("user", &model.User{
				UID: uid,
				},
			)
		})

		handler.NewHandler(&handler.Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responRecorder, request)

		respBody, err := json.Marshal(gin.H{
			"user": mockUserResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, 200, responRecorder.Code)
		assert.Equal(t, respBody, responRecorder.Body.Bytes())
		// assert that UserService.Get was called
		mockUserService.AssertExpectations(t)
	})

	t.Run("NoContextUser", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

		// a response recorder for getting written http response
		responRecorder := httptest.NewRecorder()

		// do not append user to context
		router := gin.Default()
		handler.NewHandler(&handler.Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responRecorder, request)

		assert.Equal(t, 500, responRecorder.Code)
		mockUserService.AssertNotCalled(t, "Get", mock.Anything)
	})

	t.Run("NotFound", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down call chain"))

		// a response recorder for getting written http response
		responRecorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &model.User{
				UID: uid,
			},
			)
		})

		handler.NewHandler(&handler.Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responRecorder, request)

		respErr := apperrors.NewNotFound("user", uid.String())

		respBody, err := json.Marshal(gin.H{
			"error": respErr,
		})
		assert.NoError(t, err)

		assert.Equal(t, respErr.Status(), responRecorder.Code)
		assert.Equal(t, respBody, responRecorder.Body.Bytes())
		mockUserService.AssertExpectations(t) // assert that UserService.Get was called
	})
}
