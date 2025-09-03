package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/sirupsen/logrus"
)

func (h *Handler) createTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input dto.CreateTweetRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to create tweet - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tweet, err := h.tweetService.CreateTweet(c.Request.Context(), userID.(int), &input)
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
	c.JSON(http.StatusCreated, tweet)
}

func (h *Handler) updateTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to update tweet - invalid tweetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
		return
	}

	var input dto.UpdateTweetRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to update tweet - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.tweetService.UpdateTweet(c.Request.Context(), userID.(int), tweetID, &input)
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

	c.JSON(http.StatusOK, response)
}

func (h *Handler) likeTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to like tweet - invalid tweetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
		return
	}

	if err = h.tweetService.LikeTweet(c.Request.Context(), tweetID); err != nil {
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
	tweetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to unlike tweet - invalid tweetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
		return
	}

	if err = h.tweetService.UnLikeTweet(c.Request.Context(), tweetID); err != nil {
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
	tweetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to retweet - invalid tweetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
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
	retweetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || retweetID == 0 {
		logrus.WithError(err).Error("failed to delete retweet - invalid retweetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retweetID"})
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
	tweetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweetID"})
		return
	}

	var input dto.CreateTweetRequest
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Error("failed to reply - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tweet, err := h.tweetService.ReplyToTweet(c.Request.Context(), userID.(int), tweetID, &input)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID,
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to reply - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"tweet_id": tweet.ID,
	}).Info("reply created")
	c.JSON(http.StatusCreated, tweet)

}

func (h *Handler) getRepliesToTweet(c *gin.Context) {

}

func (h *Handler) mineReplies(c *gin.Context) {

}

func (h *Handler) getTweetById(c *gin.Context) {

}

func (h *Handler) getReplies(c *gin.Context) {

}

func (h *Handler) getTweetsByUsername(c *gin.Context) {

}

func (h *Handler) getLikes(c *gin.Context) {

}

func (h *Handler) deleteTweet(c *gin.Context) {

}
