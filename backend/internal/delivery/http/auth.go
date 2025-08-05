package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/kust1q/Zapp/backend/internal/service/auth"
)

func (h *Handler) signUp(c *gin.Context) {
	var input dto.CreateUserRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorReponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.authService.SignUp(c.Request.Context(), input)
	if err != nil {
		if err == auth.ErrEmailAlreadyUsed || err == auth.ErrUsernameAlreadyUsed {
			newErrorReponse(c, http.StatusConflict, err.Error())
			return
		}
		newErrorReponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) signIn(c *gin.Context) {

}

func (h *Handler) refresh(c *gin.Context) {

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
