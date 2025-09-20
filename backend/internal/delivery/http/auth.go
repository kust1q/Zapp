package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/kust1q/Zapp/backend/internal/service/auth"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c *gin.Context) {
	var input dto.SignUpRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to sign up")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.authService.SignUp(c.Request.Context(), &input)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyUsed) {
			logrus.WithField("email", input.Email).Warn("sign up failed - email already used")
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already used",
			})
			return
		}
		if errors.Is(err, auth.ErrUsernameAlreadyUsed) {
			logrus.WithField("username", input.Username).Warn("sign up failed - username already used")
			c.JSON(http.StatusConflict, gin.H{
				"error": "username already used",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err,
		}).Error("sign up failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithFields(logrus.Fields{
		"email":   input.Email,
		"user_id": user.ID,
	}).Info("user registered")
	c.JSON(http.StatusCreated, *user)
}

func (h *Handler) signIn(c *gin.Context) {
	var input dto.SignInRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to sign in - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	response, err := h.authService.SignIn(c.Request.Context(), &input)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			logrus.WithField("email", input.Email).Warn("sign in failed - invalid credentials")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err,
		}).Error("sign in failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, *response)
}

func (h *Handler) refresh(c *gin.Context) {
	var input dto.RefreshRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("token refresh failed - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	response, err := h.authService.Refresh(c.Request.Context(), &input)
	if err != nil {
		if errors.Is(err, auth.ErrTokenNotFound) || errors.Is(err, auth.ErrInvalidRefreshToken) {
			logrus.Warn("token refresh failed - invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		} else if errors.Is(err, auth.ErrUserNotFound) {
			logrus.Warn("token refresh failed - user not found")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}
		logrus.WithError(err).Error("token refresh failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithFields(logrus.Fields{
		"old_refresh": input.Refresh[:10] + "...",
		"access":      response.Access[:10] + "...",
	}).Info("token refreshed")
	c.JSON(http.StatusOK, *response)
}

func (h *Handler) signOut(c *gin.Context) {
	var input dto.RefreshRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed sign out - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authService.SignOut(c.Request.Context(), &input); err != nil {
		logrus.WithError(err).Error("failed to sign out - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.Info("user signed out")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed out",
	})
}

func (h *Handler) updateSecuritySettings(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input dto.UpdateSecuritySettingsRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to update security settings - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.authService.UpdateSecuritySettings(c.Request.Context(), userID.(int), &input); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("update security settings failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("user_id", userID).Info("security settings updated")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully update security settings",
	})
}

func (h *Handler) resetPassword(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input dto.ResetPasswordRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to reset password - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), userID.(int), &input); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("reset password failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("user_id", userID).Info("password reset")
	c.JSON(http.StatusOK, gin.H{"message": "successfully update password"})
}

func (h *Handler) forgotPassword(c *gin.Context) {
	var input dto.ForgotPasswordRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed with forgotten password - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), &input); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err,
		}).Error("forgot password failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("email", input.Email).Info("password reset requested")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully update password",
	})
}
