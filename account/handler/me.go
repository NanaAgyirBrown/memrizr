package handler

import (
	"github.com/gin-gonic/gin"
	model2 "github.com/nanaagyirbrown/memrizr/handler/model"
	apperrors2 "github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"log"
	"net/http"
)

// Me handler calls services for getting
// a user's details
func (h *Handler) Me(c *gin.Context){

	user, exists := c.Get("user")

	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors2.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*model2.User).UID

	// gin.Context satisfies go's context.Context interface
	ctx := c.Request.Context()
	u, err := h.UserService.Get(ctx, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors2.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":u,
	})
}
