package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getFeed(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	feed, err := h.feedService.GetUserFeedByUserId(c.Request.Context(), userID.(int))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to get feed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("user_id", userID.(int)).Error("successfuly get feed")
	c.JSON(http.StatusOK, feed)
}
