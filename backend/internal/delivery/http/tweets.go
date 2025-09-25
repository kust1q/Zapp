package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/sirupsen/logrus"
)

const (
	maxMemoryForm = 1024 * 1024 * 1024 // 1 GB
)

func (h *Handler) createTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var request dto.CreateTweetRequest
	var tweet *dto.TweetResponse

	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(maxMemoryForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
			return
		}
		request.Content = c.PostForm("content")
	} else {
		if err := c.BindJSON(&request); err != nil {
			logrus.WithError(err).Error("failed to create tweet - invalid request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("create tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if fileHeader == nil && strings.TrimSpace(request.Content) == "" {
		logrus.WithError(err).Error("failed to create tweet - impossible create empty tweet")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible create empty tweet"})
		return
	} else if fileHeader != nil {
		openedFile, err := fileHeader.Open()
		if err != nil {
			logrus.WithError(err).Error("failed to create tweet - open file error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "open file error"})
			return
		}
		defer openedFile.Close()
		file := dto.FileData{
			File:   openedFile,
			Header: fileHeader,
		}
		tweet, err = h.tweetService.CreateTweetWithMedia(c.Request.Context(), userID.(int), &request, &file)
	} else {
		tweet, err = h.tweetService.CreateTweet(c.Request.Context(), userID.(int), &request)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("create tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweet.ID,
	}).Info("tweet created")
	c.JSON(http.StatusCreated, *tweet)
}

func (h *Handler) updateTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to update tweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	var request dto.UpdateTweetRequest
	if err := c.BindJSON(&request); err != nil {
		logrus.WithError(err).Error("failed to update tweet - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.tweetService.UpdateTweet(c.Request.Context(), userID.(int), tweetID, &request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("update tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, *response)
}

func (h *Handler) likeTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to like tweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err = h.tweetService.LikeTweet(c.Request.Context(), userID.(int), tweetID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("like tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweetID,
	}).Info("tweet liked")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully like tweet",
	})
}

func (h *Handler) unlikeTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to unlike tweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err = h.tweetService.UnlikeTweet(c.Request.Context(), userID.(int), tweetID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("unlike tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweetID,
	}).Info("tweet unliked")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully unlike tweet",
	})
}

func (h *Handler) retweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to retweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err = h.tweetService.CreateRetweet(c.Request.Context(), userID.(int), tweetID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("retweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweetID,
	}).Info("successfully retweet")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully retweet",
	})
}

func (h *Handler) deleteRetweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	retweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || retweetID == 0 {
		logrus.WithError(err).Error("failed to delete retweet - invalid retweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retweet id"})
		return
	}

	if err = h.tweetService.DeleteRetweet(c.Request.Context(), userID.(int), retweetID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":    userID,
			"retweet_id": retweetID,
			"error":      err,
		}).Error("retweet delete failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":    userID,
		"retweet_id": retweetID,
	}).Info("successfully delete retweet")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete retweet",
	})
}

func (h *Handler) replyToTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	var request dto.CreateTweetRequest
	var tweet *dto.TweetResponse

	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(maxMemoryForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
			return
		}
		request.Content = c.PostForm("content")
	} else {
		if err := c.BindJSON(&request); err != nil {
			logrus.WithError(err).Error("failed to create tweet - invalid request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("create tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if fileHeader == nil && strings.TrimSpace(request.Content) == "" {
		logrus.WithError(err).Error("failed to create tweet - impossible create empty tweet")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible create empty tweet"})
		return
	} else if fileHeader != nil {
		openedFile, err := fileHeader.Open()
		if err != nil {
			logrus.WithError(err).Error("failed to create tweet - open file error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "open file error"})
			return
		}
		defer openedFile.Close()
		file := dto.FileData{
			File:   openedFile,
			Header: fileHeader,
		}
		tweet, err = h.tweetService.ReplyToTweetWithMedia(c.Request.Context(), userID.(int), tweetID, &request, &file)
	} else {
		tweet, err = h.tweetService.ReplyToTweet(c.Request.Context(), userID.(int), tweetID, &request)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("create tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweetID,
		"error":    err,
	}).Info("tweet created")
	c.JSON(http.StatusCreated, *tweet)
}

func (h *Handler) getReplies(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	replies, err := h.tweetService.GetRepliesToTweet(c.Request.Context(), tweetID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get replies - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	logrus.WithField("tweet_id", tweetID).Info("replies got")
	c.JSON(http.StatusOK, replies)
}

func (h *Handler) getTweetById(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	tweet, err := h.tweetService.GetTweetById(c.Request.Context(), tweetID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get tweet - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	logrus.WithField("tweet_id", tweetID).Info("tweet got")
	c.JSON(http.StatusOK, *tweet)
}

func (h *Handler) getTweetsAndRetweetsByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		logrus.Error("failed to get tweet by username - invalid username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	tweets, err := h.tweetService.GetTweetsAndRetweetsByUsername(c.Request.Context(), username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get tweets - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	logrus.WithField("username", username).Info("tweets got")
	c.JSON(http.StatusOK, tweets)

}

func (h *Handler) getLikes(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	likes, err := h.tweetService.GetLikes(c.Request.Context(), tweetID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get likes - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	logrus.WithField("tweet_id", tweetID).Info("likes got")
	c.JSON(http.StatusOK, likes)
}

func (h *Handler) deleteTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to delete - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err := h.tweetService.DeleteTweet(c.Request.Context(), userID.(int), tweetID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("tweet delete failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweetID,
	}).Info("successfully delete tweet")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete tweet",
	})
}
