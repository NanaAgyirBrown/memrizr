package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"log"
	"net/http"
)

// signupReq won't be exported
type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req signupReq

	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	err := h.UserService.Signup(ctx, u)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})

		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to clear/delete the created user in the database

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}