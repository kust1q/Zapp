package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user, err := h.userService.GetMe(c.Request.Context(), userID.(int))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to get user profile - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("user_id", userID.(int)).Info("successfully get profile")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully get profile",
	})
	c.JSON(http.StatusOK, *user)
}

func (h *Handler) updateMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	var request dto.UpdateBioRequest
	if err := c.BindJSON(&request); err != nil {
		logrus.WithError(err).Error("failed to update user bio - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.userService.Update(c.Request.Context(), userID.(int), &request); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to update bio - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("user_id", userID.(int)).Info("successfully update bio")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully upate bio",
	})
}

func (h *Handler) deleteMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID.(int)); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to delete user - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithField("user_id", userID.(int)).Info("successfully delete user")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete user",
	})
}

func (h *Handler) followers(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		logrus.Error("failed to get tweet by username - invalid username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	followers, err := h.userService.GetFollowers(c.Request.Context(), username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get followers - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("username", username).Info("successfully get followers")
	c.JSON(http.StatusOK, followers)
}

func (h *Handler) following(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		logrus.Error("failed to get tweet by username - invalid username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	followings, err := h.userService.GetFollowings(c.Request.Context(), username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get followings - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("username", username).Info("successfully get followers")
	c.JSON(http.StatusOK, followings)
}

func (h *Handler) followUser(c *gin.Context) {
	followerID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	followingID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || followingID == 0 {
		logrus.WithError(err).Error("failed follow - invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	follow, err := h.userService.FollowToUser(c.Request.Context(), followerID.(int), followingID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"follower_id":  followerID,
			"following_id": followingID,
			"error":        err,
		}).Error("failed to follow - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithFields(logrus.Fields{
		"follower_id":  followerID,
		"following_id": followingID,
	}).Info("successfully follow")
	c.JSON(http.StatusOK, follow)
}

func (h *Handler) unfollowUser(c *gin.Context) {
	followerID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	followingID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || followingID == 0 {
		logrus.WithError(err).Error("failed unfollow - invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.UnfollowUser(c.Request.Context(), followerID.(int), followingID); err != nil {
		logrus.WithFields(logrus.Fields{
			"follower_id":  followerID,
			"following_id": followingID,
			"error":        err,
		}).Error("failed to unfollow - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithFields(logrus.Fields{
		"follower_id":  followerID,
		"following_id": followingID,
	}).Info("successfully unfollow")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully unfollow",
	})
}

func (h *Handler) getUserProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		logrus.Error("failed to get tweet by username - invalid username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	profile, err := h.userService.GetUserProfile(c.Request.Context(), username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get user profile - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("username", username).Info("successfully get profile")
	c.JSON(http.StatusOK, profile)
}
