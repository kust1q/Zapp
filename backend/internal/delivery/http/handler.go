package http

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")

	{
		auth := api.Group("/auth")
		{
			auth.POST("/sigh-up", h.signUp)
			auth.POST("/sigh-in", h.signIn)
			auth.POST("/sigh-out", h.signOut)
			auth.POST("/request-verify-token", h.reqVerify)
			auth.POST("/verify", h.verify)
			auth.POST("/forgot-password", h.forgotPasswod)
			auth.POST("/reset-password", h.resetPassword)
		}

		tweet := api.Group("/tweet")
		{
			tweet.POST("/", h.postTweet)
			tweet.POST("/:id/like", h.postLike)
			tweet.POST("/:id/reply", h.postReply)

			tweet.GET("/:id", h.getTweetById)
			tweet.GET("/:id/replies", h.getReplies)
			tweet.GET("/user/:user_id", h.getTweetByUserId)
			tweet.GET("/:id/likes", h.getLikes)

			tweet.DELETE("/:id", h.deleteTweet)
			tweet.DELETE("/:id/like", h.deleteLike)

			retweet := tweet.Group("/:id/retweet")
			{
				retweet.POST("/", h.postRetweet)
				retweet.DELETE("/", h.deleteRetweet)
			}
		}

		user := api.Group("/user")
		{
			user.GET("/me", h.getMe)
			user.PUT("/me", h.putMe)
			user.GET("/:username", h.getByUsername)
			user.POST("/avatar", h.postAvatar)
			user.GET("/:username/followers", h.followers)
			user.GET("/:username/following", h.following)
			user.DELETE("/me", h.delMe)
			user.DELETE("/avatar", h.delAvatar)

			follow := user.Group("/:username/follow")
			{
				follow.POST("/", h.postFollow)
				follow.DELETE("/", h.delFollow)
			}
		}

		feed := api.Group("/feed")
		{
			feed.GET("/", h.getFeed)
			feed.GET("/explore", h.expFeed)
			feed.GET("/trends", h.trendsFeed)
		}

		media := api.Group("/media")
		{
			media.POST("/upload", h.upMedia)
			media.DELETE("/:id", h.delMedia)
		}

		search := api.Group("/search")
		{
			search.GET("/", h.search)
			search.GET("/users", h.searchUser)
			search.GET("/tags", h.searchTags)
		}
	}

	return router
}
