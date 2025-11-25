package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/controllers/http/dto/request"
	conv "github.com/kust1q/Zapp/backend/internal/pkg/conv/dto"
	"github.com/kust1q/Zapp/backend/internal/service/auth"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c *gin.Context) {
	var req request.SignUp
	if err := c.BindJSON(&req); err != nil {
		logrus.WithError(err).Error("failed to sign up")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.authService.SignUp(c.Request.Context(), conv.FromSignUpRequestToDomain(&req))
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyUsed) {
			logrus.WithField("email", req.Email).Warn("sign up failed - email already used")
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already used",
			})
			return
		}
		if errors.Is(err, auth.ErrUsernameAlreadyUsed) {
			logrus.WithField("username", req.Username).Warn("sign up failed - username already used")
			c.JSON(http.StatusConflict, gin.H{
				"error": "username already used",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("sign up failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithFields(logrus.Fields{
		"email":   req.Email,
		"user_id": user.ID,
	}).Info("user registered")
	c.JSON(http.StatusCreated, conv.FromDomainToSignUpResponse(user))
}

func (h *Handler) signIn(c *gin.Context) {
	var req request.SignIn
	if err := c.BindJSON(&req); err != nil {
		logrus.WithError(err).Error("failed to sign in - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	tokens, err := h.authService.SignIn(c.Request.Context(), conv.FromSignInRequestToDomain(&req))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			logrus.WithField("email", req.Email).Warn("sign in failed - invalid credentials")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("sign in failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)

	c.SetCookie(
		RefreshTokenCookieName,
		tokens.Refresh.Refresh,
		int(h.authService.GetRefreshTTL().Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, conv.FromDomainToAccessResponse(tokens))
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		logrus.WithError(err).Error("failed to get refresh token from cookie")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tokens, err := h.authService.Refresh(c.Request.Context(), conv.FromRefreshRequestToDomain(refreshToken))
	if err != nil {
		if errors.Is(err, auth.ErrTokenNotFound) || errors.Is(err, auth.ErrInvalidRefreshToken) {
			c.SetCookie(RefreshTokenCookieName, "", -1, "/", "", false, true)
			logrus.Warn("token refresh failed - invalid token, cookie removed")
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

	c.SetSameSite(http.SameSiteStrictMode)

	c.SetCookie(
		RefreshTokenCookieName,
		tokens.Refresh.Refresh,
		int(h.authService.GetRefreshTTL().Seconds()),
		"/",
		"",
		false,
		true,
	)

	logrus.WithFields(logrus.Fields{
		"old_refresh": refreshToken[:10] + "...",
		"new_refresh": tokens.Refresh.Refresh[:10] + "...",
		"access":      tokens.Access.Access[:10] + "...",
	}).Info("tokens refreshed")
	c.JSON(http.StatusOK, conv.FromDomainToAccessResponse(tokens))
}

func (h *Handler) signOut(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		c.SetCookie(RefreshTokenCookieName, "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "successfully signed out"})
		return
	}

	if err := h.authService.SignOut(c.Request.Context(), conv.FromRefreshRequestToDomain(refreshToken)); err != nil {
		logrus.WithError(err).Error("failed to sign out - internal server error")
		c.SetCookie(RefreshTokenCookieName, "", -1, "/", "", false, true) // Удалить cookie
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	logrus.Info("user signed out")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed out",
	})
}

func (h *Handler) updatePassword(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req request.UpdatePassword
	if err := c.BindJSON(&req); err != nil {
		logrus.WithError(err).Error("failed to reset password - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.authService.UpdatePassword(c.Request.Context(), conv.FromResetPasswordRequestToDomain(userID.(int), &req)); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("reset password failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("user_id", userID.(int)).Info("password reset")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully update password",
	})
}

func (h *Handler) forgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	if err := c.BindJSON(&req); err != nil {
		logrus.WithError(err).Error("failed with forgotten password - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	recovery, err := h.authService.ForgotPassword(c.Request.Context(), conv.FromForgotPasswordRequestToDomain(&req))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("forgot password failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("email", req.Email).Info("password reset requested	")
	c.JSON(http.StatusOK, conv.FromDomainToRecoveryResponse(recovery))
}

func (h *Handler) recoveryPassword(c *gin.Context) {
	var req request.RecoveryPassword
	if err := c.BindJSON(&req); err != nil {
		logrus.WithError(err).Error("failed to reset password - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.authService.RecoveryPassword(c.Request.Context(), conv.FromRecoveryPasswordRequestToDomain(&req)); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":          err,
			"recovery_token": req.RecoveryToken[:10] + "...",
		}).Error("recovery password failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("recovery_token", req.RecoveryToken[:10]+"...").Info("password recovery")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully recovery password",
	})
}
