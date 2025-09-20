package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user, err := h.userService.GetUserByID()
}

func (h *Handler) updateMe(c *gin.Context) {

}

func (h *Handler) deleteMe(c *gin.Context) {

}

func (h *Handler) getByUsername(c *gin.Context) {

}

func (h *Handler) followers(c *gin.Context) {

}

func (h *Handler) following(c *gin.Context) {

}

func (h *Handler) followUser(c *gin.Context) {

}

func (h *Handler) unfollowUser(c *gin.Context) {

}
