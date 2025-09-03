package http

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

type AuthService interface {
	SignUp(ctx context.Context, user *dto.SignUpRequest) (dto.SignUpResponse, error)
	SignIn(ctx context.Context, credential *dto.SignInRequest) (dto.SignInResponse, error)
	Refresh(ctx context.Context, token *dto.RefreshRequest) (dto.SignInResponse, error)
	SignOut(ctx context.Context, token *dto.RefreshRequest) error
	VerifyAccessToken(tokenString string) (int, error)
	UpdateSecuritySettings(ctx context.Context, userID int, req *dto.UpdateSecuritySettingsRequest) error
	ResetPassword(ctx context.Context, userID int, req *dto.ResetPasswordRequest) error
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
}

type TweetService interface {
	CreateTweet(ctx context.Context, userID int, tweet *dto.CreateTweetRequest) (dto.TweetResponse, error)
	UpdateTweet(ctx context.Context, userID, tweetID int, req *dto.UpdateTweetRequest) (dto.UpdateTweetResponse, error)
	LikeTweet(ctx context.Context, tweetID int) error
	UnLikeTweet(ctx context.Context, tweetID int) error
	ReplyToTweet(ctx context.Context, userID, tweetID int, tweet *dto.CreateTweetRequest) (dto.TweetResponse, error)
	CreateRetweet(ctx context.Context, userID, tweetID int) error
	DeleteRetweet(ctx context.Context, userID, retweetID int) error
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

	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/sign-out", h.signOut)
		auth.POST("/forgot-password", h.forgotPassword)
	}

	public := api.Group("/public")
	{
		public.GET("/tweets/:id", h.getTweetById)
		public.GET("/tweets/:id/replies", h.getReplies)
		public.GET("/tweets/:id/likes", h.getLikes)
		public.GET("/users/:username", h.getByUsername)
		public.GET("/users/:username/tweets", h.getTweetsByUsername)
		public.GET("/users/:username/followers", h.followers)
		public.GET("/users/:username/following", h.following)
		public.GET("/search", h.search)
	}

	protected := api.Group("/protected", h.authMiddleware)
	{
		protected.PUT("/reset-password", h.resetPassword)
		protected.PUT("/security-settings", h.updateSecuritySettings)

		tweets := protected.Group("/tweets")
		{
			tweets.POST("", h.createTweet)
			tweets.PUT("/:id", h.updateTweet)
			tweets.DELETE("/:id", h.deleteTweet)

			tweets.POST("/:id/like", h.likeTweet)
			tweets.DELETE("/:id/like", h.unlikeTweet)

			tweets.POST("/:id/reply", h.replyToTweet)
			tweets.GET("/:id/replies", h.getRepliesToTweet)

			tweets.POST("/:id/retweet", h.retweet)
			tweets.DELETE("/retweets/:id", h.deleteRetweet)
		}

		media := protected.Group("/media")
		{
			media.POST("", h.uploadMedia)
			media.GET("/:id", h.getMedia)
			media.DELETE("/:id", h.deleteMedia)
		}

		users := protected.Group("/users")
		{
			users.GET("/me", h.getMe)
			users.PUT("/me", h.updateMe)
			users.DELETE("/me", h.deleteMe)
			users.POST("/:username/follow", h.followUser)
			users.DELETE("/:username/follow", h.unfollowUser)
			users.GET("/replies", h.mineReplies)
		}

		feed := protected.Group("/feed")
		{
			feed.GET("/", h.getFeed)
			feed.GET("/explore", h.exploreFeed)
			feed.GET("/trends", h.trendsFeed)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
	})

	return router
}
