package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type AuthService interface {
	CreateUser(entity.User) (int, error)
}

type TweetService interface {
}

type UserService interface {
}

type SearchService interface {
}

type FeedService interface {
}

type MediaService interface {
}

type Handler struct {
	authService   AuthService
	tweetService  TweetService
	userService   UserService
	searchService SearchService
	feedService   FeedService
	mediaService  MediaService
}

func NewHandler(
	authService AuthService,
	tweetService TweetService,
	userService UserService,
	searchService SearchService,
	feedService FeedService,
	mediaService MediaService,
) *Handler {
	return &Handler{
		authService:   authService,
		tweetService:  tweetService,
		userService:   userService,
		searchService: searchService,
		feedService:   feedService,
		mediaService:  mediaService,
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sigh-up", h.signUp)
		auth.POST("/sigh-in", h.signIn)
		auth.POST("/sigh-out", h.signOut)
		auth.POST("/request-verify-token", h.reqVerify)
		auth.POST("/verify", h.verify)
		auth.POST("/forgot-password", h.forgotPasswod)
		auth.POST("/reset-password", h.resetPassword)
	}

	api := router.Group("/api")

	{
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

			media := tweet.Group("/media")
			{
				media.POST("/upload", h.upMedia)
				media.DELETE("/:id", h.delMedia)
			}
		}

		user := api.Group("/user")
		{
			user.GET("/me", h.getMe)
			user.PUT("/me", h.putMe)
			user.GET("/:username", h.getByUsername)
			user.GET("/:username/followers", h.followers)
			user.GET("/:username/following", h.following)
			user.DELETE("/me", h.delMe)

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

		search := api.Group("/search")
		{
			search.GET("/", h.search)
			search.GET("/users", h.searchUser)
			search.GET("/tags", h.searchTags)
		}
	}

	return router
}
