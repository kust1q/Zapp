package http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	RefreshTokenCookieName = "refresh_token"
)

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

	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.PATCH("/refresh", h.refresh)
		auth.DELETE("/sign-out", h.signOut)
		auth.POST("/forgot-password", h.forgotPassword)
		auth.PATCH("/recovery-password", h.recoveryPassword)
	}

	public := api.Group("/public")
	{
		public.GET("/tweets/:tweet_id", h.getTweetById)
		public.GET("/tweets/:tweet_id/replies", h.getReplies)
		public.GET("/tweets/:tweet_id/likes", h.getLikes)
		public.GET("/tweets/media/:tweet_id", h.getTweetMedia)
		public.GET("/tweets/search", h.searchTweets)
		public.GET("/users/:username/profile", h.getUserProfile)
		public.GET("/users/:username/tweets", h.getTweetsAndRetweetsByUsername)
		public.GET("/users/:username/followers", h.followers)
		public.GET("/users/:username/following", h.following)
		public.GET("/users/avatar/:user_id", h.getAvatar)
		public.GET("/users/search", h.searchUsers)
	}

	protected := api.Group("/protected", h.authMiddleware)
	{
		protected.PUT("/reset-password", h.updatePassword)

		tweets := protected.Group("/tweets")
		{
			tweets.POST("", h.createTweet)
			tweets.PATCH("/:tweet_id", h.updateTweet)
			tweets.DELETE("/:tweet_id", h.deleteTweet)

			tweets.POST("/:tweet_id/like", h.likeTweet)
			tweets.DELETE("/:tweet_id/like", h.unlikeTweet)

			tweets.POST("/:tweet_id/reply", h.replyToTweet)

			tweets.POST("/:tweet_id/retweet", h.retweet)
			tweets.DELETE("/:tweet_id/retweet", h.deleteRetweet)

			tweets.DELETE("/:tweet_id", h.deleteTweetMedia)
		}

		users := protected.Group("/users")
		{
			users.GET("/me", h.getMe)
			users.PATCH("/me", h.updateMe)
			users.DELETE("/me", h.deleteMe)
			users.POST("/:user_id/follow", h.followUser)
			users.DELETE("/:user_id/follow", h.unfollowUser)
		}

		protected.GET("/feed", h.getFeed)
	}

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
	})

	return router
}
