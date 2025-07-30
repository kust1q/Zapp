package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (h *Handler) signUp(c *gin.Context) {
	var in dto.CreateUserRequest
	if err := c.BindJSON(&in); err != nil {
		newErrorReponse(c, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *Handler) signIn(c *gin.Context) {

}

func (h *Handler) signOut(c *gin.Context) {

}

func (h *Handler) reqVerify(c *gin.Context) {

}

func (h *Handler) verify(c *gin.Context) {

}

func (h *Handler) resetPassword(c *gin.Context) {

}

func (h *Handler) forgotPasswod(c *gin.Context) {

}
