package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

func (h *Handler) search(c *gin.Context) {
	query := c.Query("q")
	searchType := c.Query("type")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	if searchType == "user" {
		users, err := h.searchService.SearchUsers(c.Request.Context(), query)
		if err != nil {
			logrus.WithError(err).WithField("query", query).Error("failed to search users")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}
		if users == nil {
			users = []entity.User{}
		}
		c.JSON(http.StatusOK, gin.H{"users": users})

	}

	tweets, err := h.searchService.SearchTweets(c.Request.Context(), query)
	if err != nil {
		logrus.WithError(err).WithField("query", query).Error("failed to search tweets")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	if tweets == nil {
		tweets = []entity.Tweet{}
	}
	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}
