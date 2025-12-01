package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	conv "github.com/kust1q/Zapp/backend/internal/controllers/http/conv"
	"github.com/sirupsen/logrus"
)

func (h *Handler) deleteTweetMedia(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to get tweet media - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
		return
	}

	if err := h.mediaService.DeleteTweetMedia(c.Request.Context(), tweetID, userID.(int)); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to delete tweet media - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID.(int),
		"tweet_id": tweetID,
	}).Info("successfully delete tweet media")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete tweet media",
	})
}

func (h *Handler) getTweetMedia(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to get tweet media - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
		return
	}

	tweetMedia, err := h.mediaService.GetMediaDataByTweetID(c.Request.Context(), tweetID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get tweet media - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
	logrus.WithFields(logrus.Fields{
		"tweet_id": tweetID,
	}).Info("successfully get tweet media")
	c.JSON(http.StatusOK, conv.FromDomainToMediaResponse(tweetMedia))
}

func (h *Handler) getAvatar(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID == 0 {
		logrus.WithError(err).Error("failed to get avatar - invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	avatar, err := h.mediaService.GetAvatarDataByUserID(c.Request.Context(), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("failed to get avatar - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
	logrus.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("successfully get avatar")
	c.JSON(http.StatusOK, conv.FromDomainToAvatarResponse(avatar))
}
